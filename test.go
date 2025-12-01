package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var name string
	var age int
	var sex bool
	var delay time.Duration
	flag.StringVar(&name, "name", "张三", "姓名")
	flag.IntVar(&age, "age", 18, "年龄")
	flag.BoolVar(&sex, "sex", true, "性别")
	flag.DurationVar(&delay, "delay", 0, "延迟")
	flag.Parse()

	fmt.Println(name, age, sex, delay)
}
