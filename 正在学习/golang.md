# Golang 技巧 #

<br>

## channel ##

<font size=4 face=K>

	ch chan type
    从管道接受：<-ch
	从管道接收并保存变量，需要等号:var v type = <- ch


## 类型转换 ##

- int to string

    `strconv.Itoa`

- float - int

    `int(float), float(int)`

## append函数 ##

append函数必须有变量接收

## Time包 ##


## init函数 ##

- 程序执行前的包的初始化函数，在main（）之前自动执行

- 不能被其他函数调用

- 同一个文件可以拥有多个init（）函数



## 编程技巧 ##

- 结构体变量使用指针传递

- 不要想着使用一个数组变量存取数据

- 进行参数检测

- 使用多重赋值



