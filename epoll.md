# epoll


## L0:

server 会创建一个 socket， 我们称为 listener，这个 socket 有一个 fd

将 listener fd 注册到 epoll ，并说明要监听写入事件， 然后 epoll 就会返回发生在 listener fd 上的写入事件

接收到 listener fd 的事件后，用 Accept 接收到 conn fd， 之后将 conn fd 也注册到 epoll [# L2]

后续 conn fd 上的读写事件应由对应的业务函数处理， 一般是开一个 goroutine 处理


## L1:

socket上定义了几个IO事件：
    状态改变事件、有数据可读事件、有发送缓存可写事件、有IO错误事件。

对于这些事件，socket中分别定义了相应的事件处理函数，也称回调函数。

Socket I/O事件的处理过程中，要使用到sock上的两个队列：
    等待队列 和 异步通知队列，
这两个队列中都保存着等待该Socket I/O事件的进程。

等待队列上的进程会睡眠，直到Socket I/O事件的发生，然后在事件处理函数中被唤醒。（阻塞）

异步通知队列上的进程则不需要睡眠，Socket I/O事件发时，事件处理函数会给它们发送到信号，这些进程事先注册的信号处理函数就能够被执行。

### 对于 connection

三次握手中，当客户端收到SYNACK、发出ACK后，连接就成功建立了。

此时连接的状态从TCP_SYN_SENT或TCP_SYN_RECV变为TCP_ESTABLISHED，sock的状态发生变化，
并进入睡眠，直到超时或收到信号，或者被I/O事件处理函数唤醒。

等待结束时，把等待进程从等待队列中删除，把当前进程的状态设为TASK_RUNNING，
进程被唤醒，connect() 就能成功返回了。


##L2:

三次握手中，当服务器端接收到ACK完成连接建立的时候，会把新的连接链入 全连接队列 中，
然后唤醒监听listener socket上的等待进程，accept()就能成功返回 conn 了

##L3:

使用 epoll 等API，一般需要一个 eventfd 用于中断 epoll wait， 这个 eventfd 我们暂时称为 `wakeUpEventFd`

关于 eventfd 的概念， 这里只需要知道可以对 eventfd 读和写

当 epoll 对多个 fd 监听时 （wait函数），主程序阻塞，CPU切换到内核态，直到一个fd发生了事件被监听到，主程序才会被唤醒，得到 epoll wait() 返回值

所以，当所有 fd 都没有发生事件时， epoll 就会一直阻塞
为了使 epoll wait 跳出阻塞， 我们就可以给 `wakeUpEventFd` 写数据，这样就可以令 epoll 监听到 wakeUpEventFd 的写事件，epoll wait 就可以返回


### 关于 eventfd

```c
    #include<sys/eventfd.h>
    int eventfd(unsigned int initval, int flags);
```

使用这个函数来创建一个事件对象。linux线程间通信为了提高效率，大多使用异步通信，采用事件监听和回调函数的方式来实现高效的任务处理方式

linux内核会为这个事件对象维护一个64位的计数器(`uint64_t`).并在初始化时用传进去的initval来初始化这个计数器，然后返回一个文件描述符来代表这个事件对象

第二个参数是描述这个事件对象的属性，可以设置为
    `EFD_NONBLOCK`,
    `EFD_CLOEXEC`

前面的是设置对象为非阻塞状态
    如果没有设置为非阻塞状态，read系统调用来读这个计数器，且计数器的值为0时，就会一直阻塞在read系统调用上
    反之如果设置了该标志位，就会返回`EAGAIN`错误

后面的`EFD_CLOEXEC`功能是在程序调用`exec()`函数族加载其他程序时自动关闭当前已有的文件描述符

通过此函数得到的文件描述符既然是一个计数器，我们就可以对它进行读和写：

使用write将缓冲区写入的8字节整形值加到内核计数器上。
使用read将内核计数的8字节值读取到缓冲区中，并把计数器重设为0，如果buffer的长度小于8字节则read会失败，错误码设为EINVAl。


## 参考

[https://blog.csdn.net/u012319493/article/details/99211464](https://blog.csdn.net/u012319493/article/details/99211464)

[https://www.cnblogs.com/burningTheStar/p/7064193.html](https://www.cnblogs.com/burningTheStar/p/7064193.html)

[https://www.jianshu.com/p/d7ebac8dc9f8](https://www.jianshu.com/p/d7ebac8dc9f8)

[https://blog.csdn.net/zhangskd/article/details/45770323](https://blog.csdn.net/zhangskd/article/details/45770323)

[https://www.jianshu.com/p/2704cd87200a](https://www.jianshu.com/p/2704cd87200a)

![image](https://raw.githubusercontent.com/fumeboy/go_practices/main/epoll.png)