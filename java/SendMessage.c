#include "SendMessage.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <mqueue.h>

JNIEXPORT jint JNICALL Java_SendMessage_send
  (JNIEnv *env, jobject obj, jstring name, jint flags, jstring text){
 mqd_t mq;
 const char *cstr = (*env)->GetStringUTFChars(env, name, 0);
 mq = mq_open(cstr, flags);
 if((mqd_t)-1 == mq){
  (*env)->ReleaseStringUTFChars(env, name, cstr);
  perror("mq_open");
  return -1;
 }
 const char *ctxt = (*env)->GetStringUTFChars(env, text, 0);
 int res = mq_send(mq, ctxt, strlen(ctxt)+1, 0);
 if(res==-1){
  perror("mq_send");
  return -1;
 }
 if((mqd_t)-1 == mq_close(mq)){
  perror("mq_close");
  return -1;
 }
 (*env)->ReleaseStringUTFChars(env, name, cstr);
 (*env)->ReleaseStringUTFChars(env, text, ctxt);
 return 0;
}
