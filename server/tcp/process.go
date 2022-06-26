package tcp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type result struct {
	value []byte
	err   error
}

func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d", len(errString)) + " " + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	// 这里的长度后面必须有一个空格，
	// 因为在cache-benchmark/tcpClient.recvResponse()里需要通过readString(' ')来读取长度
	// 需要这个空格做划分,帮助识别到value的长度
	valueLength := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(valueLength), value...)) // 这里...的用法：通过append合并两个slice
	return e
}

func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c

	// 从命令里读key
	key, err := s.readKey(r)
	if err != nil {
		c <- &result{
			nil,
			err,
		}
		return
	}
	// 根据key, 向cache里找value
	go func() {
		value, err := s.Get(key)
		c <- &result{
			value: value,
			err:   err,
		}
	}()
}

func (s *Server) set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	key, value, err := s.readKeyAndValue(r)
	if err != nil {
		c <- &result{
			nil,
			err,
		}
		return
	}

	// 匿名函数自动继承外部的所有可见变量
	// 当匿名函数在goroutine中启动后，上层函数就会返回
	// 这样Server.Process方法就能立即开始处理下一个请求
	go func() {
		c <- &result{
			nil,
			s.Set(key, value),
		}
	}()
}

func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	key, err := s.readKey(r)
	if err != nil {
		c <- &result{
			nil,
			err,
		}
		return
	}

	go func() {
		c <- &result{
			nil,
			s.Del(key),
		}
	}()
}

func (s *Server) process(conn net.Conn) {
	// 用来对客户端连接进行缓冲读取，
	// 因为来自网络的数据不稳定，在读取时，客户端的数据可能只传输了一半
	// 可以阻塞等待，直到需要的数据全部就位以后一次性返回给我们
	r := bufio.NewReader(conn)
	resultCh := make(chan chan *result, 5000)
	defer close(resultCh)
	go reply(conn, resultCh)

	for {
		//fmt.Println("等待请求到达！")
		op, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Println("close connection due to error: ", err)
			}
			return
		}

		// 每个操作对应一个chan *result，里面有一个*result
		if op == 'S' {
			s.set(resultCh, r)
		} else if op == 'G' {
			s.get(resultCh, r)
		} else if op == 'D' {
			s.del(resultCh, r)
		} else {
			log.Println("close connection due to invalid operation: ", op)
			return
		}
	}
}

func reply(conn net.Conn, resultCh chan chan *result) {
	defer conn.Close()

	for {
		c, open := <-resultCh
		if !open {
			return
		}
		// 每个chan *result里只有一个*result{}
		r := <-c
		err := sendResponse(r.value, r.err, conn)
		if err != nil {
			log.Println("close connection due to error: ", err)
			return
		}
	}
}
