# protobuf学习 #

<font size=4 face=d>

## 1. protobuf ##

高效、灵活以及可扩展的序列化数据结构，用于跨平台、跨语言的数据交换和数据存储

## 使用protobuf ##

- 执行protoc --go-out=. *.proto解析proto文件，生成一个go文件

## 2. protobuf解析为golang ##

- message对应struct

- field的名字aa_bb对应struct中的AcBb

- message中的field存在getter方法

- 非repeated的的字段都被实现为一个指针

- repeated的字段被实现为slice

## proto3的语法 ##

- 枚举成员的值必须从0开始

- 允许定义多个message

- 嵌套结构的message实例名和消息名称不能相同

## 编译proto文件 ##

    protoc -I=SRC＿Dir --go_out=DST＿DIR　SRC_DIR/FILENAEM.proto

##  ##