package tcp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

func readLen(r *bufio.Reader) (int, error) {
	tmp, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}
	length, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}
	return length, nil
}

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	// 对键进行检查，看该键是否应该由本节点处理
	// 如果该键不应由本节点处理，服务端回返回一个"redirect <新节点地址>"的错误，客户端就能意识到需要重新读取集群节点
	keyLength, err := readLen(r)
	if err != nil {
		return "", err
	}
	keyByte := make([]byte, keyLength)
	_, err = io.ReadFull(r, keyByte)
	if err != nil {
		return "", err
	}
	key := string(keyByte)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", errors.New("redirect " + addr)
	}

	return key, nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	keyLength, err := readLen(r)
	if err != nil {
		return "", nil, err
	}
	valueLength, err := readLen(r)
	if err != nil {
		return "", nil, err
	}

	// 对键进行检查，看该键是否应该由本节点处理
	// 如果该键不应由本节点处理，服务端回返回一个"redirect <新节点地址>"的错误，客户端就能意识到需要重新读取集群节点
	keyByte := make([]byte, keyLength)
	_, err = io.ReadFull(r, keyByte)
	if err != nil {
		return "", nil, err
	}
	key := string(keyByte)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", nil, errors.New("redirect " + addr)
	}

	value := make([]byte, valueLength)
	_, err = io.ReadFull(r, value)
	if err != nil {
		return "", nil, err
	}

	return key, value, nil
}
