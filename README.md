# Native-Support---Java-Vs-Golang
Comparing Java vs Golang in context of native support by taking IPC(Inter Process Communication) through POSIX message queue as an usecase


The linux manual page - http://man7.org/linux/man-pages/man7/mq_overview.7.html, give us about how to interact/create posix message queues

To call any one of the methods in the manual page, the only possible way is through C language code.
Core logic must be written in C and then have to wrap this C language code in Java / Golang

Let's compare by executing the methods
- [mq_open(3)](http://man7.org/linux/man-pages/man3/mq_open.3.html)
- [mq_send(3)](http://man7.org/linux/man-pages/man3/mq_send.3.html)
- [mq_receive(3)](http://man7.org/linux/man-pages/man3/mq_receive.3.html)
- [mq_close(3)](http://man7.org/linux/man-pages/man3/mq_close.3.html)

## Java
  we have an option to use **javah**  which can link the native methods in java to methods in C language using **jni**

- From [GetMessage.java](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/java/GetMessage.java), have to generate [GetMessage.h](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/java/GetMessage.h) and the command to be used **(path-to-java)/bin/javah -jni GetMessage**
- The above generated GetMessage.h is just like an interface in java where we need to write implementation in C which is [GetMessage.c](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/java/GetMessage.c). This is where we call **mq_receive()**
- Now its time to generate **libGetMessage.so**, the command to execute is *gcc -I (path-to-java)/include -I (path-to-java)/include/linux -o libGetMessage.so -shared -fPIC GetMessage.c -lrt*
- In the runtime, we just need **GetMessage.class** and **libGetMessage.so** and the command for receiving messages from queue is *java -Djava.library.path=(path-of-executable-files) GetMessage*
- For usage of [mq_open(3)](http://man7.org/linux/man-pages/man3/mq_open.3.html) and [mq_close(3)](http://man7.org/linux/man-pages/man3/mq_close.3.html), check the files [PosixMQ.java](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/java/PosixMQ.java) , [PosixMQ.c](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/java/PosixMQ.c)
- To send messages to the queue , have a look on [SeneMessage.java](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/java/SendMessage.java), [SendMessage.c](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/java/SendMessage.c)


## Golang
- Golang comes with in-built C support and they call it as "CGO" . Due to this we can write directly C code in the go source file itself Unlike java, we don't need special binaries like javah / jni.
- [wrapper.go](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/golang/posx_mq/wrapper.go) contains C language code which contains logic related to [mq_open(3)](http://man7.org/linux/man-pages/man3/mq_open.3.html), [mq_send(3)](http://man7.org/linux/man-pages/man3/mq_send.3.html), [mq_receive(3)](http://man7.org/linux/man-pages/man3/mq_receive.3.html), [mq_close(3)](http://man7.org/linux/man-pages/man3/mq_close.3.html)
- [receiver.go](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/golang/receiver.go) / [sender.go](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/golang/sender.go) contains logic to receive/send messages to/from Posix queue through [mq.go](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/golang/posx_mq/mq.go),[wrapper.go](https://github.com/kranthiB/Native_Support---Java_vs_Golang/blob/master/golang/posx_mq/wrapper.go)
