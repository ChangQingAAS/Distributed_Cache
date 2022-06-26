package scanner

type Scanner interface {
	Scan() bool    // 返回TRUE意味着后续还有未遍历的键值对，FALSE意味着遍历结束
	Key() string   // 访问当前键值对的Key
	Value() []byte // 访问当前键值对的Value
	Close()        // 结束遍历
}
