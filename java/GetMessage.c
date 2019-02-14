#include "GetMessage.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <mqueue.h>

JNIEXPORT jstring JNICALL Java_GetMessage_getMessage
  (JNIEnv *env, jobject obj, jstring name, jint flags, jint priority){
 unsigned int up = priority;
 char ret[1025];
 const char *cstr = (*env)->GetStringUTFChars(env, name, 0);
 mqd_t mq;
 mq = mq_open(cstr, flags);
 if((mqd_t)-1 == mq){

  (*env)->ReleaseStringUTFChars(env, name, cstr);
  perror("mq_open");
  return (*env)->NewStringUTF(env,"dud");
 }
 ssize_t bytes_read = mq_receive(mq, ret, 1024, &up);
 if(bytes_read==-1){
  perror("mq_receive");
  return (*env)->NewStringUTF(env,"dud");
 }
 ret[bytes_read] = '\0';
 if((mqd_t)-1 == mq_close(mq)){
  perror("mq_close");
  return (*env)->NewStringUTF(env,"dud");
 }
 (*env)->ReleaseStringUTFChars(env, name, cstr);
 return (*env)->NewStringUTF(env,ret);
}
