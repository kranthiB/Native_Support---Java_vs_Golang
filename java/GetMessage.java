public class GetMessage{
 public native String getMessage(String name, int flags, int priority);
 public static void main(String[] args) {
  System.out.println(new GetMessage().getMessage("/test_queue",0,0));
 }
 static {
  System.loadLibrary("GetMessage");
 }
}
