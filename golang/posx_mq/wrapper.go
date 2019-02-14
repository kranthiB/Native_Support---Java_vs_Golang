package posx_mq

/*
#include <stdlib.h>
#include <signal.h>
#include <fcntl.h>
#include <mqueue.h>

// Expose non-variadic function requires 4 arguments.
mqd_t mq_open4(const char *name, int oflag, int mode, struct mq_attr *attr) {
	return mq_open(name, oflag, mode, attr);
}
*/
import "C"
import (
	"cadx/util"
	"fmt"
	"strings"
	"unsafe"
)

const (
	O_RDONLY = C.O_RDONLY
	O_WRONLY = C.O_WRONLY
	O_RDWR   = C.O_RDWR

	O_CLOEXEC  = C.O_CLOEXEC
	O_CREAT    = C.O_CREAT
	O_EXCL     = C.O_EXCL
	O_NONBLOCK = C.O_NONBLOCK

	S_IRUSR = C.S_IRUSR
	S_IRGRP = C.S_IRGRP
	S_IWOTH = C.S_IWOTH

	INT_MAX = 2147483647
)

var (
	MemoryAllocationError = fmt.Errorf("Memory Allocation Error")
)

type receiveBuffer struct {
	buf  *C.char
	size C.size_t
}

type timeSpec struct {
	tv_sec  float64
	tv_nsec float64
}

func newReceiveBuffer(bufSize int) (*receiveBuffer, error) {
	buf := (*C.char)(C.malloc(C.size_t(bufSize)))
	if buf == nil {
		return nil, MemoryAllocationError
	}

	return &receiveBuffer{
		buf:  buf,
		size: C.size_t(bufSize),
	}, nil
}

func (rb *receiveBuffer) free() {
	C.free(unsafe.Pointer(rb.buf))
}

//The definition of HARD_MSGMAX has changed across kernel
//	versions:
//	    *  Up to Linux 2.6.32: 131072 / sizeof(void *)
//      *  Linux 2.6.33 to 3.4: (32768 * sizeof(void *) / 4)
//	    *  Since Linux 3.5: 65,536
//The definition of HARD_MSGSIZEMAX has changed across kernel
//	versions:
//	    *  Before Linux 2.6.28, the upper limit is INT_MAX.
//      *  From Linux 2.6.28 to 3.4, the limit is 1,048,576.
//	    *  Since Linux 3.5, the limit is 16,777,216
//http://man7.org/linux/man-pages/man7/mq_overview.7.html
func getHardMsgMaxAndMsgSizeMaxUpperLimit(version string) (int, int) {
	versions := strings.Split(version, ".")
	kernelVersion := util.ConvertStringToInt(versions[0])
	majorVersion := util.ConvertStringToInt(versions[1])
	minorVersion := util.ConvertStringToInt(versions[2])
	var x interface{}
	if kernelVersion == 2 {
		if majorVersion == 6 {
			if minorVersion < 28 {
				return 131072 / int(unsafe.Sizeof(unsafe.Pointer(&x))), INT_MAX
			}
			if minorVersion >= 28 && minorVersion < 33 {
				return 131072 / int(unsafe.Sizeof(unsafe.Pointer(&x))), 1048576
			}
			if minorVersion >= 33 {
				return (32768 * int(unsafe.Sizeof(unsafe.Pointer(&x)))) / 4, 1048576
			}
		}
		if majorVersion > 6 {
			return (32768 * int(unsafe.Sizeof(unsafe.Pointer(&x)))) / 4, 1048576
		}
		return 131072 / int(unsafe.Sizeof(unsafe.Pointer(&x))), INT_MAX
	}
	if kernelVersion == 3 {
		if majorVersion <= 4 {
			return (32768 * int(unsafe.Sizeof(unsafe.Pointer(&x)))) / 4, 1048576
		}
		return 65536, INT_MAX
	}
	if kernelVersion > 3 {
		return 65536, 16777216
	}
	return 131072 / int(unsafe.Sizeof(unsafe.Pointer(&x))), INT_MAX
}

func mq_open(name string, oflag int, mode int, maxMessages int, maxMessageSize int) (int, error) {
	var cAttr *C.struct_mq_attr
	cAttr = &C.struct_mq_attr{
		mq_maxmsg:  C.long(maxMessages),
		mq_msgsize: C.long(maxMessageSize),
	}

	h, err := C.mq_open4(C.CString(name), C.int(oflag), C.int(mode), cAttr)
	if err != nil {
		return 0, err
	}

	return int(h), nil
}

func mq_send(h int, data []byte, priority uint) (int, error) {
	byteStr := *(*string)(unsafe.Pointer(&data))
	rv, err := C.mq_send(C.int(h), C.CString(byteStr), C.size_t(len(data)), C.uint(priority))
	return int(rv), err
}

func mq_receive(h int, recvBuf *receiveBuffer) ([]byte, uint, error) {
	var msgPrio C.uint

	size, err := C.mq_receive(C.int(h), recvBuf.buf, recvBuf.size, &msgPrio)
	if err != nil {
		return nil, 0, err
	}

	return C.GoBytes(unsafe.Pointer(recvBuf.buf), C.int(size)), uint(msgPrio), nil
}

func mq_timedreceive(h int, recvBuf *receiveBuffer) ([]byte, uint, error) {
	var msgPrio C.uint
	var absTimeOut C.struct_timespec
	absTimeOut.tv_sec = 100 / 1000
	absTimeOut.tv_nsec = (100 % 1000) * 1000000

	size, err := C.mq_timedreceive(C.int(h), recvBuf.buf, recvBuf.size, &msgPrio, &absTimeOut)
	if err != nil {
		return nil, 0, err
	}

	return C.GoBytes(unsafe.Pointer(recvBuf.buf), C.int(size)), uint(msgPrio), nil
}

func mq_notify(h int, sigNo int) (int, error) {
	sigEvent := &C.struct_sigevent{
		sigev_notify: C.SIGEV_SIGNAL, // posix_mq supports only signal.
		sigev_signo:  C.int(sigNo),
	}

	rv, err := C.mq_notify(C.int(h), sigEvent)
	return int(rv), err
}

func mq_close(h int) (int, error) {
	rv, err := C.mq_close(C.int(h))
	return int(rv), err
}

func mq_unlink(name string) (int, error) {
	rv, err := C.mq_unlink(C.CString(name))
	return int(rv), err
}
