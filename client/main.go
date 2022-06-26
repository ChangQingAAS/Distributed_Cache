package main

import (
	"distribute_cache/cache-benchmark/cacheClient"
	"flag"
	"fmt"
)

func main() {
	// flag标准包：解析命令行参数
	// 第一个参数指定命令行相应参数的名字，第二个参数是变量的默认值，第三个参数是对该参数的一个描述字符串，
	// 返回值是该类型的指针
	server := flag.String("h", "localhost", "cache server address")
	op := flag.String("c", "get", "command, coule be get/set/del")
	key := flag.String("k", "", "key")
	value := flag.String("v", "", "value")
	flag.Parse()

	client := cacheClient.New("tcp", *server)
	cmd := &cacheClient.Cmd{
		Name:  *op,
		Key:   *key,
		Value: *value,
		Error: nil,
	}
	client.Run(cmd)
	if cmd.Error != nil {
		fmt.Println("client-bin happen error: ", cmd.Error)
	} else {
		fmt.Println(cmd.Value)
	}
	return
}
