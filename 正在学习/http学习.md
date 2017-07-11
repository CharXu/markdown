# http包的结构 #

<font size=4 face=K>

&emsp;http协议是用于从www服务器传输超文本到本地浏览器位于传送协议，基于TCP/IP通信协议传递数据，数据的传递格式是http包。

&emsp;http协议由请求和相应构成，总是由客户端向服务器发起请求，服务器收到请求想客户端返回响应。

## 1. http请求包的格式 ##

客户端向服务器发送http的请求有三个部分构成：

- 第一行：请求方法+URL+协议/版本

	常用请求方法：GET和POST

	URL：指定的请求的资源所在的服务器位置

	
- 中间：请求头部

	包含客户端环境和请求正文的相关信息


- 最后一行：请求正文


## 2. http的应答包 ##

服务器收到客户端的请求，向客户端发送应答包：

- 第一行：协议+状态代码+描述

	协议/版本：HTTP/1.0

	状态代码：200请求成功，404资源不存在，304网页未修改，503服务器暂时不可用，500服务器内部错误

	描述：

- 中间：应答包头部

- 最后一行：应答正文

	服务器返回的html页面


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



