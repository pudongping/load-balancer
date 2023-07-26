// 一致性 hash 代码参考：https://juejin.cn/post/6871169933150486542
package consistentHash

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// hash 函数需要符合以下要求：
// 1 单调性（唯一） 2平衡性 (数据 目标元素均衡) 3分散性(散列)
type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int           { return len(s) }
func (s UInt32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s UInt32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type ConsistentHashBalance struct {
	mux      sync.RWMutex
	hash     Hash                // 哈希函数类型 Hash 的实例，用于计算节点的哈希值
	replicas int                 // 复制因子，表示每个节点在哈希环上的虚拟节点个数
	keys     UInt32Slice         // 已排序的节点哈希切片，用于快速查找节点
	hashMap  map[uint32]string   // 节点哈希和键的映射，键是节点的哈希值，值是节点的标识（如 IP 地址）
	nodes    map[string][]uint32 // 键和节点哈希之间的映射关系，键是节点的标识（如 IP 地址），值是节点的哈希值切片
}

func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	if replicas <= 0 {
		panic("replicas must greater than 0")
	}

	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
		nodes:    make(map[string][]uint32),
	}
	if m.hash == nil {
		// 最多 32 位，保证是一个 2^32-1 环
		// 默认使用 crc32.ChecksumIEEE 作为哈希函数
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (c *ConsistentHashBalance) IsEmpty() bool {
	// 这里只需要知道所有的虚拟环哈希值是否为空即可
	return len(c.keys) == 0
}

// AddNode 方法用来添加节点，参数为节点 key，比如使用 IP 地址
func (c *ConsistentHashBalance) AddNode(addr string) error {
	if "" == addr {
		return errors.New("the node cannot be null")
	}
	if _, ok := c.nodes[addr]; ok {
		return errors.New("the node has already exists")
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	// 结合复制因子计算所有虚拟节点的 hash 值，并存入 m.keys 中，
	// 同时在 m.hashMap 中保存哈希值和 key 的映射
	// 以及在 m.nodes 中保存 key 和哈希值之间的映射关系
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + "-" + addr))
		c.keys = append(c.keys, hash)               // 记录所有的虚拟节点哈希值
		c.hashMap[hash] = addr                      // 添加虚拟节点和真实节点之间的映射关系
		c.nodes[addr] = append(c.nodes[addr], hash) // 添加真实节点和虚拟哈希值之间的关系
	}

	// 对所有虚拟节点的哈希值进行升序排序，方便之后进行二分查找
	// 这里和 php 处理有点儿区别，php 可以直接对关联数组进行排序，比较感叹 php 的数组功能是真的强大，😄
	sort.Sort(c.keys)

	return nil
}

// RemoveNode 方法用于移除缓存节点，参数为节点key，比如使用IP
func (c *ConsistentHashBalance) RemoveNode(addr string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.nodes[addr]; !ok {
		// 根本就不存在改节点时，则不需要被移除
		return
	}

	// 需要移除的节点的哈希值
	removeHashes := c.nodes[addr]

	// 从 keys 和 hashMap 中移除对应哈希值的节点
	newKeys := make(UInt32Slice, 0, len(c.keys))
	for _, key := range c.keys {
		if !contains(removeHashes, key) {
			newKeys = append(newKeys, key)
		} else {
			// 需要一个一个的删除虚拟节点和真实节点之间的映射关系
			delete(c.hashMap, key)
		}
	}
	c.keys = newKeys

	// 删除掉真实节点和虚拟节点之间的映射关系
	delete(c.nodes, addr)
}

// contains 函数用于判断切片中是否包含某个元素
func contains(s []uint32, val uint32) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}

// Lookup 方法根据给定的对象获取最靠近它的那个节点
func (c *ConsistentHashBalance) Lookup(key string) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("node is empty")
	}

	hash := c.hash([]byte(key))

	// 通过二分查找函数获取最优节点，第一个"服务器hash"值大于等于"数据hash"值的就是最优"服务器节点"
	idx := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })

	// 如果查找结果 大于 服务器节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，那么放入第一个节点中
	if idx == len(c.keys) {
		idx = 0
	}
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.hashMap[c.keys[idx]], nil
}
