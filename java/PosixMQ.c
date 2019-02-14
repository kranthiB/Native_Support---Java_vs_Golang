#include "PosixMQ.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <mqueue.h>

JNIEXPORT jlong JNICALL Java_PosixMQ_mq_1open
  (JNIEnv *env, jobject obj, jstring name, jint flags){
    struct mq_attr attr;
    attr.mq_flags = 0;
    attr.mq_maxmsg = 10;
    attr.mq_msgsize = 1024;
    attr.mq_curmsgs = 0;
 const char *str = (*env)->GetStringUTFChars(env, name, 0);
 mqd_t desc;
 desc = mq_open(str, flags, 0644, &attr);
 (*env)->ReleaseStringUTFChars(env, name, str);
 return (intptr_t)desc;
}

JNIEXPORT jint JNICALL Java_PosixMQ_mq_1close
  (JNIEnv *env, jobject obj, jlong desc){
 return mq_close((mqd_t)(intptr_t)desc);
}
