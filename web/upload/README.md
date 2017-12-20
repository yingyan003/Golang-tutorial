uploadDemo:上传文件
==================

### HTTP方式
------------------

### Request
------------------

* Method: Post
* URL: http://localhost:8000/upload

### Response
------------------

* 客户端返回：fmt.Fprintf(w,"fileHeader.Header:\n%v",fileHeader.Header)
* 返回值示例：
fileHeader.Header:
map[Content-Disposition:[form-data; name="uploadFile"; filename="1.txt"] Content-Type:[text/plain]]


### 使用说明
-------------------

1. 启动服务端:run uploadServer.go
2. 打开html页面（可在系统文件夹下直接双击打开）
3. 选择文件，上传
