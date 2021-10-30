### go语言进阶第九周作业  
==================

处理socket粘包的三种方式：  
# 1. Fix length  
发送方和接收方规定固定大小的缓冲区，也就是发送和接收都使用固定大小的 byte[] 数组长度，当字符长度不够时使用空字符弥补.   
## 示例：  
```java
/**  
 * 服务器端
 */  
static class ServSocketV1 {  
    private static final int BYTE_LENGTH = 1024;  // 字节数组长度（收消息用）  
    public static void main(String[] args) throws IOException {  
        ServerSocket serverSocket = new ServerSocket(9091);  
        // 获取到连接  
        Socket clientSocket = serverSocket.accept();  
        try (InputStream inputStream = clientSocket.getInputStream()) {  
            while (true) {  
                byte[] bytes = new byte[BYTE_LENGTH];  
                // 读取客户端发送的信息  
                int count = inputStream.read(bytes, 0, BYTE_LENGTH);  
                if (count > 0) {  
                    // 接收到消息打印  
                    System.out.println("接收到客户端的信息是:" + new String(bytes).trim());  
                }  
                count = 0;  
            }  
        }  
    }  
}  
```
```java
/**
 * 客户端
 */
static class ClientSocketV1 {
    private static final int BYTE_LENGTH = 1024;  // 字节长度
    public static void main(String[] args) throws IOException {
        Socket socket = new Socket("127.0.0.1", 9091);
        final String message = "Hi,Java."; // 发送消息
        try (OutputStream outputStream = socket.getOutputStream()) {
            // 将数据组装成定长字节数组
            byte[] bytes = new byte[BYTE_LENGTH];
            int idx = 0;
            for (byte b : message.getBytes()) {
                bytes[idx] = b;
                idx++;
            }
            // 给服务器端发送 10 次消息
            for (int i = 0; i < 10; i++) {
                outputStream.write(bytes, 0, BYTE_LENGTH);
            }
        }
    }
}
```
# 2. Delimiter based  
以特殊的字符结尾，比如以“\n”结尾，这样我们就知道结束字符，从而避免了半包和粘包问题.   
## 示例：    
```java
/**
 * 服务器端
 */
static class ServSocketV3 {
    public static void main(String[] args) throws IOException {
        // 创建 Socket 服务器端
        ServerSocket serverSocket = new ServerSocket(9092);
        // 获取客户端连接
        Socket clientSocket = serverSocket.accept();
        // 使用线程池处理更多的客户端
        ThreadPoolExecutor threadPool = new ThreadPoolExecutor(100, 150, 100,
                TimeUnit.SECONDS, new LinkedBlockingQueue<>(1000));
        threadPool.submit(() -> {
            // 消息处理
            processMessage(clientSocket);
        });
    }
    /**
     * 消息处理
     * @param clientSocket
     */
    private static void processMessage(Socket clientSocket) {
        // 获取客户端发送的消息流对象
        try (BufferedReader bufferedReader = new BufferedReader(
                new InputStreamReader(clientSocket.getInputStream()))) {
            while (true) {
                // 按行读取客户端发送的消息
                String msg = bufferedReader.readLine();
                if (msg != null) {
                    // 成功接收到客户端的消息并打印
                    System.out.println("接收到客户端的信息:" + msg);
                }
            }
        } catch (IOException ioException) {
            ioException.printStackTrace();
        }
    }
}
```
```java
/**
 * 客户端
 */
static class ClientSocketV3 {
    public static void main(String[] args) throws IOException {
        // 启动 Socket 并尝试连接服务器
        Socket socket = new Socket("127.0.0.1", 9092);
        final String message = "Hi,Java."; // 发送消息
        try (BufferedWriter bufferedWriter = new BufferedWriter(
                new OutputStreamWriter(socket.getOutputStream()))) {
            // 给服务器端发送 10 次消息
            for (int i = 0; i < 10; i++) {
                // 注意:结尾的 \n 不能省略,它表示按行写入
                bufferedWriter.write(message + "\n");
                // 刷新缓冲区(此步骤不能省略)
                bufferedWriter.flush();
            }
        }
    }
}
```
# 3. Length field based frame decoder  
在TCP协议的基础上封装一层数据请求协议，既将数据包封装成数据头（存储数据正文大小）+ 数据正文的形式，这样在服务端就可以知道每个数据包的具体长度了，知道了发送数据的具体边界之后，就可以解决半包和粘包的问题了 
##  示例： 
本项目
