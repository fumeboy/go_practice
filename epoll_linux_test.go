package main_test

import (
	"sync"
	"syscall"
	"testing"
	"unsafe"
)

const (
	maxEvents = 1024

	EventRead  uint32 = 0x1
	EventWrite uint32 = 0x2
	EventErr   uint32 = 0x80

	EPOLLET_ = 0x80000000
	// syscall.EPOLLET 的值是 -0x80000000，应该是错误的， 很奇怪
)

type event struct {
	FD    int
	Event uint32
}

type epoll struct {
	fd                     int
	wakeUpEventFd          int
	wakeUpEventFdSignalOut []byte

	toAdd      []int // 可以考虑改成 chan， 不过 chan 有容量限制
	toAddMutex sync.Mutex

	closeChan chan struct{}
	closeOnce sync.Once
	finishChan chan struct{}
}

func (p *epoll) New() (*epoll, error) {
	fd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}
	r0, _, e0 := syscall.Syscall(syscall.SYS_EVENTFD2, 0, 0, 0) // 调用 eventfd 函数
	if e0 != 0 {
		syscall.Close(fd)
		return nil, err
	}
	if err := syscall.EpollCtl(fd, syscall.EPOLL_CTL_ADD, int(r0),
		&syscall.EpollEvent{Fd: int32(r0),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		syscall.Close(fd)
		syscall.Close(int(r0))
		return nil, err
	}
	p.fd = fd
	p.wakeUpEventFd = int(r0)
	p.wakeUpEventFdSignalOut = make([]byte, 8)
	p.closeChan = make(chan struct{})
	return p, err
}

func (p *epoll) Add(fd int) error {
	p.toAddMutex.Lock()
	p.toAdd = append(p.toAdd, fd)
	p.toAddMutex.Unlock()

	return p.wakeup() // 使用 wakeup 使 epoll wait 结束阻塞
}

func (p *epoll) Close() error {
	p.closeOnce.Do(func() {
		close(p.closeChan)
	})
	return p.wakeup()
}

func (p *epoll) wakeup() error {
	var x uint64 = 1 // 非 0 值
	_, err := syscall.Write(p.wakeUpEventFd, (*(*[8]byte)(unsafe.Pointer(&x)))[:])
	return err
}

func (p *epoll) Run(handler func(events []event)) {
	defer func() {
		syscall.Close(p.fd)
		syscall.Close(p.wakeUpEventFd)
		p.fd = -1
		p.wakeUpEventFd = -1
		p.finishChan <- struct{}{}
	}()
	events := make([]syscall.EpollEvent, maxEvents)
	for {
		select {
		case <-p.closeChan:
			return
		default:
			p.toAddMutex.Lock()
			for _, fd := range p.toAdd {
				syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, int(fd), &syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLRDHUP | syscall.EPOLLIN | syscall.EPOLLOUT | EPOLLET_})
			}
			p.toAdd = p.toAdd[:0]
			p.toAddMutex.Unlock()

			n, err := syscall.EpollWait(p.fd, events, -1)
			if err == syscall.EINTR {
				continue
			}
			if err != nil {
				return
			}
			var pe []event
			for i := 0; i < n; i++ {
				ev := &events[i]
				if int(ev.Fd) == p.wakeUpEventFd {
					syscall.Read(p.wakeUpEventFd, p.wakeUpEventFdSignalOut)
				} else {
					e := event{FD: int(ev.Fd)}

					// EPOLLRDHUP (since Linux 2.6.17)
					// Stream socket peer closed connection, or shut down writing
					// half of connection.  (This flag is especially useful for writ-
					// ing simple code to detect peer shutdown when using Edge Trig-
					// gered monitoring.)
					if ((events[i].Events & syscall.EPOLLHUP) != 0) && ((events[i].Events & syscall.EPOLLIN) == 0) {
						e.Event |= EventErr
					}
					if (events[i].Events&syscall.EPOLLERR != 0) || (events[i].Events&syscall.EPOLLOUT != 0) {
						e.Event |= EventWrite
					}
					if events[i].Events&(syscall.EPOLLIN|syscall.EPOLLPRI|syscall.EPOLLRDHUP) != 0 {
						e.Event |= EventRead
					}
					pe = append(pe, e)
				}
			}
			handler(pe)
		}
	}
}

// Test


func TestEpoll(t *testing.T)  {
	// TODO
}

