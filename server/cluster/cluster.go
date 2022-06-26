package cluster

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"io/ioutil"
	"stathat.com/c/consistent"
	"time"
)

type Node interface {
	// ShouldProcess 告诉节点该key是否应该由自己处理
	ShouldProcess(key string) (string, bool)
	// Members 提供整个集群的节点列表
	Members() []string
	// Addr 获取本节点的地址
	Addr() string
}

type node struct {
	*consistent.Consistent
	addr string
}

func (n *node) Addr() string {
	return n.addr
}

func (n *node) ShouldProcess(key string) (string, bool) {
	address, _ := n.Get(key)
	return address, address == n.addr
}

func New(addr, cluster string) (Node, error) {

	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard // 任何写入操作都会成功且内容会被直接丢弃，使得终端免于被memberlist的日志刷屏
	list, err := memberlist.Create(conf)
	if err != nil {
		fmt.Println("hw")
		return nil, err
	}
	if cluster == "" {
		cluster = addr
	}

	clu := []string{cluster}
	_, err = list.Join(clu)
	if err != nil {
		return nil, err
	}

	circle := consistent.New()
	// 每个节点的虚拟节点的数量，默认为20
	// 当节点数较少时，20个虚拟节点还不能做到较好的负载均衡，故设为256
	circle.NumberOfReplicas = 256

	go func() {
		// 每个1s，把memberlist.Memberlist.Members()提供的集群节点列表m更新到circle中
		for {
			m := list.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()

	return &node{circle, addr}, nil
}
