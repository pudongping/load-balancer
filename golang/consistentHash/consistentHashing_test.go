package consistentHash

import (
	"strconv"
	"testing"
)

func TestConsistentHashBalance(t *testing.T) {
	// 创建一个具有3个复制因子和默认哈希函数（crc32.ChecksumIEEE）的 ConsistentHashBalance 实例
	ch := NewConsistentHashBalance(3, nil)

	// 向一致性哈希环中添加一些节点（IP地址）
	// 由于复制因子的作用，这些节点应该在哈希环中均匀分布
	nodes := []string{"192.168.0.1", "192.168.0.2", "192.168.0.3", "192.168.0.4", "192.168.0.5"}
	for _, node := range nodes {
		if err := ch.AddNode(node); err != nil {
			t.Errorf("添加节点出错 %v \n", err)
		}
	}

	// 使用一些键测试 Get() 方法，并检查它们是否被正确映射到相应的节点
	keyToNodeMapping := map[string]string{
		"apple":  "192.168.0.1",
		"banana": "192.168.0.5",
		"cherry": "192.168.0.5", // 由于一致性哈希的作用，这个键应该与 "banana" 映射到同一个节点
		"grape":  "192.168.0.4",
		"orange": "192.168.0.2",
		"pear":   "192.168.0.3",
	}

	for key, expectedNode := range keyToNodeMapping {
		node, err := ch.Lookup(key)
		if err != nil {
			t.Errorf("获取键 %s 的节点失败：%v", key, err)
		}
		if node != expectedNode {
			t.Errorf("键 %s 映射到了错误的节点。预期节点：%s，实际节点：%s", key, expectedNode, node)
		}
	}

	// 使用一个空环测试 Lookup() 方法是否返回错误
	emptyCh := NewConsistentHashBalance(3, nil)
	isMustEmpty := emptyCh.IsEmpty()
	if !isMustEmpty {
		t.Error("预期应该是一个空环，但是目前不是")
	}
	_, err := emptyCh.Lookup("test")
	if err == nil {
		t.Error("预期在空环中返回错误，但未得到错误")
	}
}

func TestConsistentHashBalance_RemoveNode(t *testing.T) {
	// 创建一个具有3个复制因子和默认哈希函数（crc32.ChecksumIEEE）的 ConsistentHashBalance 实例
	ch := NewConsistentHashBalance(3, nil)

	// 向一致性哈希环中添加一些节点（IP地址）
	nodes := []string{"192.168.0.1", "192.168.0.2", "192.168.0.3", "192.168.0.4", "192.168.0.5"}
	for _, node := range nodes {
		if err := ch.AddNode(node); err != nil {
			t.Errorf("添加节点出错 %v \n", err)
		}
	}

	// 移除一个节点并验证它是否从哈希环中正确地被移除
	removedNode := "192.168.0.2"
	ch.RemoveNode(removedNode)

	// 期望移除的节点不存在于哈希环中
	for i := 0; i < ch.replicas; i++ {
		hash := ch.hash([]byte(strconv.Itoa(i) + "-" + removedNode))
		node, err := ch.Lookup(strconv.Itoa(i))
		if err == nil && node == removedNode {
			t.Errorf("期望节点 %s 已被移除，但仍存在于虚拟节点中。虚拟节点哈希值: %d", removedNode, hash)
		}
	}
}

func TestConsistentHashBalance_RemoveNonExistentNode(t *testing.T) {
	// 创建一个具有3个复制因子和默认哈希函数（crc32.ChecksumIEEE）的 ConsistentHashBalance 实例
	ch := NewConsistentHashBalance(3, nil)

	// 向一致性哈希环中添加一些节点（IP地址）
	nodes := []string{"192.168.0.1", "192.168.0.2", "192.168.0.3", "192.168.0.4", "192.168.0.5"}
	for _, node := range nodes {
		if err := ch.AddNode(node); err != nil {
			t.Errorf("添加节点出错 %v \n", err)
		}
	}

	// 尝试移除一个不存在的节点，预期应该不会有任何影响
	nonExistentNode := "192.168.0.6"
	ch.RemoveNode(nonExistentNode)

	// 验证不存在的节点确实不存在于哈希环中
	for i := 0; i < ch.replicas; i++ {
		hash := ch.hash([]byte(strconv.Itoa(i) + "-" + nonExistentNode))
		node, err := ch.Lookup(strconv.Itoa(i))
		if err == nil && node == nonExistentNode {
			t.Errorf("不存在的节点 %s 不应存在于虚拟节点中，但找到了。虚拟节点哈希值: %d", nonExistentNode, hash)
		}
	}
}
