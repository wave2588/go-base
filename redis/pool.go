package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/wave2588/go-base/utils"
	"net"
	"net/url"
	"sync/atomic"
	"time"
)

func init() {
	//redis.SetLogger(tredis.Logger{})
}

type Pool interface {
	Master() Client
	Slave() Client
	Slaves() []Client
	Close() (err error)
}

type Options struct {
	Name       string
	MasterAddr *url.URL
	SlaveAddrs []*url.URL

	// Maximum number of retries before giving up.
	// Default is 3 retries.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 128 milliseconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 256 milliseconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	// Maximum number of socket connections.
	// Default is 60 connections per every CPU.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing new connection is slow.
	// Default is CPU counts.
	MinIdleConns int
	// Connection age at which client retires (closes) the connection.
	// Default is 5 minutes.
	MaxConnAge time.Duration
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration
	// Hooks
	Hooks []redis.Hook
	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that fifo has higher cpu overhead compared to lifo.
	PoolFIFO bool
}

func (opt *Options) init() {
	if opt.MaxRetries == -1 {
		opt.MaxRetries = 3
	}
	switch opt.MinRetryBackoff {
	case -1:
		opt.MinRetryBackoff = 0
	case 0:
		opt.MinRetryBackoff = 8 * time.Millisecond
	}
	switch opt.MaxRetryBackoff {
	case -1:
		opt.MaxRetryBackoff = 0
	case 0:
		opt.MaxRetryBackoff = 512 * time.Millisecond
	}

	if opt.DialTimeout == 0 {
		opt.DialTimeout = 128 * time.Millisecond
	}
	switch opt.ReadTimeout {
	case -1:
		opt.ReadTimeout = 0
	case 0:
		opt.ReadTimeout = 256 * time.Millisecond
	}
	switch opt.WriteTimeout {
	case -1:
		opt.WriteTimeout = 0
	case 0:
		opt.WriteTimeout = opt.ReadTimeout
	}

	if opt.PoolSize == 0 {
		opt.PoolSize = 60 * utils.TotalCPU()
	}
	if opt.MinIdleConns == 0 {
		opt.MinIdleConns = utils.TotalCPU()
	}
	if opt.MaxConnAge == 0 {
		opt.MaxConnAge = 5 * time.Minute
	}
	if opt.PoolTimeout == 0 {
		opt.PoolTimeout = opt.ReadTimeout + time.Second
	}
	if opt.IdleTimeout == 0 {
		opt.IdleTimeout = 5 * time.Minute
	}
	if opt.IdleCheckFrequency == 0 {
		opt.IdleCheckFrequency = time.Minute
	}
	if opt.Hooks == nil {
		opt.Hooks = []redis.Hook{}
	}
}

func (opt *Options) build(addr *url.URL) *redis.Options {
	password, _ := addr.User.Password()
	ropt := &redis.Options{
		Addr:               addr.Host,
		Password:           password,
		MaxRetries:         opt.MaxRetries,
		MinRetryBackoff:    opt.MinRetryBackoff,
		MaxRetryBackoff:    opt.MaxRetryBackoff,
		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,
		PoolFIFO:           opt.PoolFIFO,
	}
	host, _, _ := net.SplitHostPort(addr.Host)
	if net.ParseIP(host) == nil {
		lc, err := NewLeastConn(addr.Host, newDNSResolver())
		if err == nil {
			ropt.Dialer = lc.DialContext
		}
	}
	return ropt
}

func GetPool(name string, masterAddr *url.URL, slaveAddrs []*url.URL) (Pool, error) {
	return GetPoolWithOptions(name, masterAddr, slaveAddrs, &Options{})
}

func GetPoolWithOptions(name string, masterAddr *url.URL, slaveAddrs []*url.URL, opt *Options) (Pool, error) {
	opt.Name = name
	if opt.MasterAddr == nil {
		opt.MasterAddr = masterAddr
		opt.SlaveAddrs = slaveAddrs
	}
	return newPool(opt)
}

func newPool(opt *Options) (Pool, error) {
	opt.init()

	master := redis.NewClient(opt.build(opt.MasterAddr))

	slaves := make([]*redis.Client, 0)
	for _, addr := range opt.SlaveAddrs {
		slave := redis.NewClient(opt.build(addr))
		slaves = append(slaves, slave)
	}

	pool := &poolImpl{
		name:    opt.Name,
		master:  master,
		slaves:  slaves,
		closing: make(chan struct{}),
	}

	go pool.reporter()

	return pool, nil
}

type poolImpl struct {
	name    string
	master  *redis.Client
	slaves  []*redis.Client
	next    uint64
	closing chan struct{}
}

func (p *poolImpl) reporter() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.flushStats(fmt.Sprintf("%s-master", p.name), p.master.PoolStats())

			for i, slave := range p.slaves {
				p.flushStats(fmt.Sprintf("%s-slave-%d", p.name, i), slave.PoolStats())
			}
		case <-p.closing:
			return
		}
	}
}

func (p *poolImpl) flushStats(name string, stats *redis.PoolStats) {
	//statsd.TelemetryGauge(statsd.Join("redis_pool_stats", name, "hits"), float64(stats.Hits))
	//statsd.TelemetryGauge(statsd.Join("redis_pool_stats", name, "misses"), float64(stats.Misses))
	//statsd.TelemetryGauge(statsd.Join("redis_pool_stats", name, "timeouts"), float64(stats.Timeouts))
	//statsd.TelemetryGauge(statsd.Join("redis_pool_stats", name, "total_conns"), float64(stats.TotalConns))
	//statsd.TelemetryGauge(statsd.Join("redis_pool_stats", name, "idle_conns"), float64(stats.IdleConns))
	//statsd.TelemetryGauge(statsd.Join("redis_pool_stats", name, "stale_conns"), float64(stats.StaleConns))
}

func (p *poolImpl) Master() Client {
	return p.master
}

func (p *poolImpl) Slave() Client {
	slaveNum := uint64(len(p.slaves))
	if slaveNum == 0 {
		return p.master
	}
	return p.slaves[atomic.AddUint64(&p.next, 1)%slaveNum]
}

func (p *poolImpl) Slaves() []Client {
	slaves := make([]Client, 0, len(p.slaves))
	for _, slave := range p.slaves {
		slaves = append(slaves, slave)
	}
	return slaves
}

func (p *poolImpl) Close() (err error) {
	close(p.closing)

	if err = p.master.Close(); err != nil {
		return err
	}
	for _, slave := range p.slaves {
		if err = slave.Close(); err != nil {
			return err
		}
	}
	return nil
}
