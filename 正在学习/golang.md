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

## comma-ok断言检测 ##

    value, ok := element.(Type)

- value是接口实际存储变量的返回值

- ok表示是否为真

		v, ok := <- ch chan

- 检测通道是否关闭

## 指针接收器 ##

- 当结构体通过指针接收器实现一个接口的函数的时候，实现接口不是结构体本身，而是指向接口的指针

- 接口变量只能存储指针类型的变量



    	type men interface {
			think()
		}
		type student struct{
			
		}
		func (s *student) think () {
			
		}

		var m men
		var s struct

		m = &s


## 函数的局部变量 ##

- 函数内部通过：=初始化的变量，只在函数内部有效，函数外无法访问函数内部的变量

- 但是return会自动把函数内部变量转化成全局变量
		



