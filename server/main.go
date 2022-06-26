package main

import (
	"distribute_cache/cache/cache"
	"distribute_cache/server/cluster"
	"distribute_cache/server/http"
	"distribute_cache/server/tcp"
	"flag"
	"log"
)

func main() {
	typ := flag.String("type", "inmemory", "cache type")
	ttl := flag.Int("ttl", 30, "cache time to live")
	// 本节点地址
	// 本节点启动后会向集群节点发送消息通知自己的存在
	node := flag.String("node", "127.0.0.1", "node address")
	// 需要加入的集群的某个节点地址
	clusterNode := flag.String("cluster", "", "cluster address")
	/*
		由于gossip协议的特性，集群内的任何节点接收到新节点信息后都会逐渐扩散让整个集群知晓
		所以cluster参数具体选择那个节点无关紧要，只需要是集群内已经存在的一个节点即可
	*/
	flag.Parse()
	log.Println("type is ", *typ)
	log.Println("ttl is ", *ttl)
	log.Println("node is ", *node)
	log.Println("cluster is ", *clusterNode)
	c := cache.New(*typ, *ttl)
	n, err := cluster.New(*node, *clusterNode)
	if err != nil {
		panic(err)
	}
	go tcp.New(c, n).Listen()
	http.New(c, n).Listen()
}
