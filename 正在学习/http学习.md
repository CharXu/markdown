# http包的结构 #

<font size=4 face=K>

&emsp;http协议是用于从www服务器传输超文本到本地浏览器位于传送协议，基于TCP/IP通信协议传递数据，数据的传递格式是http包。

&emsp;http协议由请求和相应构成，总是由客户端向服务器发起请求，服务器收到请求想客户端返回响应。

## 1. http请求包的格式 ##

客户端向服务器发送http的请求有三个部分构成：

- 第一行：请求方法+URL+协议/版本
	
- 中间：请求头部

	包含客户端环境和请求正文的相关信息

- 最后一行：请求正文

## 2. http的应答包 ##

服务器收到客户端的请求，向客户端发送应答包：

- 第一行：协议+状态代码+描述

- 中间：应答包头部

- 最后一行：应答正文

## 3. http的请求方法 ##

- GET

&emsp;&emsp;最常用的方法，客户端请求服务器发送某个资源

- HEAD

&emsp;&emsp;和GET方法类似，但服务器只返回请求资源的头部，而不包含主体部分，使用HEAD可以在：捕获去资源的情况下了解资源、查看响应中的状态码判断某个对象是否存在、通过查看资源首部判断资源是否被修改

- PUT

&emsp;&emsp;与PUT方法相反，PUT请求会向服务器写入文档，请求服务器用数据包的主体部分创建或者替代（如果请求URL已经存在）一个由请求的URL命名的新文档

- POST

&emsp;&emsp;POST方法用来向服务器输入数据，在HTML表单中填好的数据会被发送给服务器，然后由服务器处理数据。

- TRACE

&emsp;&emsp;允许客户端查看请求的原始版本，判断数据包是否以及如何被修改过

- OPTIONS

&emsp;&emsp;请求web服务器告知其支持的各种功能。例如某些特殊资源支持哪些请求方法

- DELETE

&emsp;&emsp;请求服务器删除URL指定的资源，但是客户端无法宝恒删除操作一定会被执行。

## 4. 状态码 ##

&emsp;&emsp;状态码方便的为客户端提供了服务器处理请求的结果

- 100-199——信息状态码

&emsp;&emsp;100：Continue，服务器收到请求的初始部分，请客户端继续。

&emsp;&emsp;101：Switching Protocal，服务器正在根据客户端的指定，将协议切换成Update首部所列的协议

- 200-299——成功状态码

&emsp;&emsp;200：OK，请求成功

&emsp;&emsp;201：Created，用于创建服务器对象的请求，返回这个请求说明对象创建成功

&emsp;&emsp;202：Accepted，请求已经被服务器接收，但是服务器还没有对请求执行任何动作。

&emsp;&emsp;203：Non-Authoritative Information，表示实体首部包含的信息不是来源于服务器本身，而是中间节点上面的副本

&emsp;&emsp;204：No Content，服务器的响应包


# Fiddler抓包 #

&emsp;Fiddler专门捕获计算机与网络之间传送的http数据包，进行通过分析数据包可以查看接口是否调用正确，数据返回是否正确，还可以对http数据包进行重发、编辑和转存。


# wireshark #

&emsp;wireshark可以截取链路层、网络层、传输层和应用层的所有数据包，并且显示网络数据包的详细信息。由于安全原因，wireshark只能查看网络包，不能修改。

## 1. 捕获过滤 ##

&emsp;wireshark捕获到的数据包数量非常大，为了便于分析，需要对对捕获的数据包进行过滤：

- IP过滤

&emsp;&emsp;ip.src eq 192.168.1.1 显示特定IP发来的数据包

&emsp;&emsp;ip.dst eq [ip-addr] 显示特定IP接收到的数据包

&emsp;&emsp;ip.addr == [ip-addr] 显示特定IP接受和发送的数据包

- 端口过滤

&emsp;&emsp;tcp.port == [port] 

&emsp;&emsp;tcp.dstport == [port]

&emsp;&emsp;tcp.srcport == [port]

&emsp;&emsp;tcp.port >= [port] 过滤某范围的端口

- 协议过滤

&emsp;&emsp;直接输入协议名称： http/tcp/udp/ftp/icmp/ssl/dns/等

&emsp;&emsp;排除协议： !http 或者 not http

- 包长度过滤

&emsp;&emsp;tcp.len

&emsp;&emsp;udp.len

&emsp;&emsp;ip.len

&emsp;&emsp;frame.len

- http模式过滤

&emsp;&emsp; http.request.method == POST/GET/PUT... 按请求方法过滤

&emsp;&emsp; http.request.uri 按请求的资源标识符过滤

&emsp;&emsp; http contains "" 按http包中的内容过滤



