package redis

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

//var logger = log.New("redis-balancer")

// TCPConn 是对 TCPConn 的包装
// 目的是为了捕获 close 事件
type TCPConn struct {
	once    sync.Once
	onClose func()
	*net.TCPConn
}

func (c *TCPConn) Close() error {
	if c.onClose != nil {
		c.once.Do(c.onClose)
	}
	return c.TCPConn.Close()
}

type resolver interface {
	solve(addr string) (list []net.TCPAddr, err error)
}

type DNSResolver struct {
	*net.Resolver
}

func newDNSResolver() *DNSResolver {
	return &DNSResolver{
		net.DefaultResolver,
	}
}

func (s *DNSResolver) solve(addr string) (list []net.TCPAddr, err error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	addrs, err := s.LookupHost(context.TODO(), host)
	if err != nil {
		return nil, err
	}
	list = make([]net.TCPAddr, 0, len(addrs))
	portNum, _ := strconv.Atoi(port)
	for _, addr := range addrs {
		list = append(list, net.TCPAddr{
			IP:   net.ParseIP(addr),
			Port: portNum,
		})
	}
	return list, nil
}

// LeastConn 最小链接的负载均衡实现
// 目的是保持 redis 链接池中，各个节点长链接的数量的平衡
type LeastConn struct {
	d        net.Dialer
	addr     string
	mu       sync.Mutex
	table    map[string]int32
	resolver resolver
	trigger  chan struct{}
	stop     chan struct{}
	stopOnce sync.Once
}

// NewLeastConn 创建负载均衡器
// addr 表示目的地址
// resolver 用来解析该地址
func NewLeastConn(addr string, resolver resolver) (*LeastConn, error) {
	l := &LeastConn{
		addr:     addr,
		table:    map[string]int32{},
		resolver: resolver,
		trigger:  make(chan struct{}, 1),
		stop:     make(chan struct{}, 1),
	}
	addrs, err := l.resolver.solve(addr)
	if err != nil {
		return nil, err
	}
	l.update(addrs)
	go l.debounceResolve()
	return l, err
}

func (l *LeastConn) update(list []net.TCPAddr) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// 使用新的table，目的是避免同一个 map 内存可能不回收的问题
	table := make(map[string]int32, len(list))
	for _, addr := range list {
		table[addr.String()] = l.table[addr.String()]
	}
	l.table = table
}

func (l *LeastConn) Close() error {
	l.stopOnce.Do(func() {
		l.stop <- struct{}{}
	})
	return nil
}

const minResolveDebounceTime = time.Millisecond * 300
const minErrDebounceTime = time.Millisecond * 30

// debounceResolve 去抖动解析
// 在 maxResolveDebounceTime 时间内，不再进行重新解析
func (l *LeastConn) debounceResolve() {
	// 上次解析成功的时间
	lastTime := time.Time{}
	// 上次解析失败的时间
	var lastErrTime *time.Time

	for {
		select {
		case <-l.stop:
			return
		case <-l.trigger:
			gap := time.Since(lastTime)
			if (lastErrTime != nil && gap > minErrDebounceTime) ||
				(gap > minResolveDebounceTime) {

				list, err := l.resolver.solve(l.addr)
				if err != nil {
					now := time.Now()
					lastErrTime = &now

					//logger.Errorf(context.TODO(), "host [%s] resolve failed: %s ", l.addr, err)
				} else {
					lastErrTime = nil
					lastTime = time.Now()
					l.update(list)
				}
			}
		}
	}
}

// tryUpdate 尝试触发地址解析
// 当链接创建/回收的时候进行尝试更新
func (l *LeastConn) tryUpdate() {
	select {
	case l.trigger <- struct{}{}:
	default:
	}
}

func (l *LeastConn) OnConnClose(addr string) {
	l.mu.Lock()
	v, ok := l.table[addr]
	if ok {
		l.table[addr] = v - 1
	}
	l.mu.Unlock()

	l.tryUpdate()
}

func (l *LeastConn) DialContext(ctx context.Context, _, _ string) (net.Conn, error) {
	addr, err := l.acquire()
	if err != nil {
		return nil, err
	}
	conn, err := l.d.DialContext(ctx, "tcp", addr)
	if err != nil {
		l.OnConnClose(addr)
		return nil, err
	}
	l.tryUpdate()
	return &TCPConn{
		onClose: func() {
			l.OnConnClose(addr)
		},
		TCPConn: conn.(*net.TCPConn),
	}, nil
}

func (l *LeastConn) acquire() (addr string, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.table) == 0 {
		return "", fmt.Errorf("no such host: %s", l.addr)
	}

	var min int32 = -1
	minAddr := ""
	pairs := make([]pair, 0, len(l.table))
	for k, v := range l.table {
		pairs = append(pairs, pair{k, v})
	}

	rand.Shuffle(len(pairs), func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})

	for _, p := range pairs {
		v := p.num
		addr := p.addr
		if min == -1 {
			min = v
			minAddr = addr
		}

		if v < min {
			min = v
			minAddr = addr
		}
	}
	l.table[minAddr]++
	return minAddr, nil
}

type pair struct {
	addr string
	num  int32
}
