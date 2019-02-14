package posix_mq

import (
	"bytes"
	"cadx/asset"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

var cadxQueue asset.CadxQueue

// Represents the message queue
type MessageQueue struct {
	handler int
	name    string
	recvBuf *receiveBuffer
}

type QueueStats struct {
	maxMessages         int64
	curMessages         int64
	queueUsedPercentage float64
}

// Represents the message queue attribute
type MessageQueueAttribute struct {
	flags   int
	maxMsg  int
	msgSize int
	curMsgs int
}

// NewMessageQueue returns an instance of the message queue given a QueueConfig.
func NewMessageQueue(filterConfigRootPath string, posixQueueMountPath string) (*MessageQueue, error) {
	cadxQueue = asset.GetCadxQueueFilterDefinition(filterConfigRootPath)
	unlinkQueue(cadxQueue)
	h, err := mq_open(cadxQueue.Name, O_CREAT|O_RDONLY, S_IRUSR|S_IRGRP|S_IWOTH, cadxQueue.MaxMessages, cadxQueue.MaxMessageSize)
	if err != nil {
		return nil, err
	}
	changeQueuePermission(cadxQueue.Name, posixQueueMountPath)
	msgSize := cadxQueue.MaxMessageSize
	recvBuf, err := newReceiveBuffer(int(msgSize))
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		handler: h,
		name:    cadxQueue.Name,
		recvBuf: recvBuf,
	}, nil
}

func getMessageQueue(cadxQueue asset.CadxQueue) (*MessageQueue, error) {
	h, err := mq_open(cadxQueue.Name, O_RDONLY, S_IRUSR|S_IRGRP|S_IWOTH, cadxQueue.MaxMessages, cadxQueue.MaxMessageSize)
	if err != nil {
		return nil, err
	}
	return &MessageQueue{
		handler: h,
		name:    cadxQueue.Name,
	}, nil
}

// GetCadxQueueDefinition - returns the posix queue definiton used while creating queue
func GetCadxQueueDefinition() asset.CadxQueue {
	return cadxQueue
}

func getKernelVersion() string {
	cmd := exec.Command("uname", "-r")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		log.Fatalf("I! error while executing uname command to get kernel version %s", err.Error())
	}
	kernelRelease := strings.TrimSpace(string(cmdOutput.Bytes()))
	return strings.Split(kernelRelease, "-")[0]
}

// Send sends message to the message queue.
func (mq *MessageQueue) Send(data []byte, priority uint) error {
	_, err := mq_send(mq.handler, data, priority)
	return err
}

// Receive receives message from the message queue.
func (mq *MessageQueue) Receive() ([]byte, uint, error) {
	return mq_receive(mq.handler, mq.recvBuf)
}

// TimedReceive receives message from the POSIX queue on a timely basis
func (mq *MessageQueue) TimedReceive() ([]byte, uint, error) {
	return mq_timedreceive(mq.handler, mq.recvBuf)
}

// FIXME Don't work because of signal portability.
// Notify set signal notification to handle new messages.
func (mq *MessageQueue) Notify(sigNo syscall.Signal) error {
	_, err := mq_notify(mq.handler, int(sigNo))
	return err
}

// Close closes the message queue.
func (mq *MessageQueue) Close() error {
	mq.recvBuf.free()

	_, err := mq_close(mq.handler)
	return err
}

// Unlink deletes the message queue.
func (mq *MessageQueue) Unlink() error {
	mq.Close()

	_, err := mq_unlink(mq.name)
	return err
}

func unlinkQueue(cadxQueue asset.CadxQueue) {
	shouldUnlinkPosixQueue := os.Getenv("SHOULD_UNLINK_POSIX_QUEUE")
	ok, err := strconv.ParseBool(shouldUnlinkPosixQueue)
	if err != nil {
		log.Printf("failed to parse SHOULD_UNLINK_POSIX_QUEUE flag value=%s\n", shouldUnlinkPosixQueue)
	}
	if ok {
		mq, err := getMessageQueue(cadxQueue)
		if err != nil {
			log.Printf("in error while RemovePosixMessageQueue, error : %s", err.Error())
		} else {
			mq_close(mq.handler)
			mq_unlink(mq.name)
		}
	}
}

func changeQueuePermission(queueName string, posixQueueMountPath string) {
	var queuePath strings.Builder
	queuePath.WriteString(posixQueueMountPath)
	queuePath.WriteString(queueName)
	err := os.Chmod(queuePath.String(), 0442)
	if err != nil {
		log.Printf("in error while changing the permisssion for queue - %s , error : %s", queuePath, err.Error())
	}
}
