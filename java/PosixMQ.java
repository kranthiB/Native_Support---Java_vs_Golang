public class PosixMQ {
 public native long mq_open(String name, int flags);
 public native int mq_close(long desc);
 public static void main(String[] args) {
  PosixMQ pmq = new PosixMQ();
  long desc = pmq.mq_open("/test_queue",64);
  if(desc!=-1L){
   System.out.println("We have MQ open.");
   System.out.println("Close returned " + pmq.mq_close(desc));
  }else{
   System.out.println("Open failed.");
  }
 }
 static {
  System.loadLibrary("PosixMQ");
 }
}
