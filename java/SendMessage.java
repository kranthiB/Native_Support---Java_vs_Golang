public class SendMessage{
 public native int send(String name, int flags, String text);
 public static void main(String[] args) {
  int status = new SendMessage().send("/test_queue",1,"Hello World!");
  if(status!=-1)
   System.out.println("Message dispatched.");
  else
   System.out.println("Sending message failed.");
 }
 static {
  System.loadLibrary("SendMessage");
 }
}
