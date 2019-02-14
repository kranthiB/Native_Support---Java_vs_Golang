# Native-Support---Java-Vs-Golang
Comparing Java vs Golang in context of native support by taking IPC(Inter Process Communication) through POSIX message queue as an usecase


The linux manual page - http://man7.org/linux/man-pages/man7/mq_overview.7.html, give us about how to interact/create posix message queues

To call any one of the methods in the manual page, the only possible way is through C language code.
Core logic must be written in C and then have to wrap this C language code in Java / Golang

Let's compare by executing the method mq_eceive(http://man7.org/linux/man-pages/man3/mq_receive.3.html)

## Java
  we have an option to use **javah**  which can link the native methods in java to methods in C language using **jni**
  
