<?php
/**
 *
 *
 * Created by PhpStorm
 * User: Alex
 * Date: 2023-07-25 21:11
 */
declare(strict_types=1);

class ConsistentHashing
{

    /**
     * 每个节点对应的虚拟节点个数
     *
     * @var int
     */
    protected $replicas;

    /**
     * 虚拟节点哈希和真实节点值之间的映射关系，key 为虚拟节点的哈希值，value 为节点的标识（如 IP 地址）
     *
     * @var array
     */
    protected $hashMap = [];

    /**
     * 真实节点和虚拟节点之间的映射关系（一个真实节点至少存在一个虚拟节点）
     *
     * @var array
     */
    protected $nodes = [];

    public function __construct(int $replicas)
    {
        if ($replicas <= 0) throw new RuntimeException('replicas must greater than 0');
        $this->replicas = $replicas;
    }

    /**
     * 把字符串转为 32 位无符号整数
     *
     * @param string $str
     * @return int
     */
    public function hash(string $str): int
    {
        // 获取十进制格式的无符号 crc32 校验和的字符串表示形式
        return (int)sprintf('%u', crc32($str));
    }

    public function isEmpty(): bool
    {
        return count($this->hashMap) == 0 || count($this->nodes) == 0;
    }

    /**
     * 添加节点
     *
     * @param string $addr
     * @return void
     */
    public function addNode(string $addr)
    {
        if (! $addr) throw new RuntimeException('the node cannot be null');
        if (isset($this->nodes[$addr])) throw new RuntimeException('the node has already exists');

        for ($i = 0; $i < $this->replicas; $i++) {
            $hash = $this->hash($i . '-' . $addr);
            $this->hashMap[$hash] = $addr;  // 添加虚拟节点和真实节点之间的映射
            $this->nodes[$addr][] = $hash;  // 添加真实节点和虚拟哈希值之间的关系
        }

        // 重新升序排序虚拟节点哈希 key
        ksort($this->hashMap);
    }

    /**
     * 删除节点
     *
     * @param string $addr
     * @return void
     */
    public function removeNode(string $addr)
    {
        if (! isset($this->nodes[$addr])) return;

        // 循环删除虚拟节点
        foreach ($this->nodes[$addr] as $hashVal) {
            unset($this->hashMap[$hashVal]);
        }

        // 删除真实节点
        unset($this->nodes[$addr]);
    }

    /**
     * 根据给定的对象获取最靠近它的那个节点。即为最佳节点。
     *
     * @param string $key
     * @return string
     */
    public function lookup(string $key): string
    {
        if ($this->isEmpty()) {
            throw new RuntimeException('node is empty');
        }

        $hash = $this->hash($key);

        // 因为每添加一个节点时，都会使用 ksort 对虚拟节点哈希做升序排序且每次也都会使用 reset 进行重新指向
        // 因此，这里默认取的就是圆环上最小的一个节点，当作默认值
        $node = current($this->hashMap);

        // 循环获取相近的最优节点：第一个“服务器hash”值大于等于“数据hash”值的就是最优“服务器节点”
        // 因为此时 $this->hashMap 已经是通过虚拟哈希值升序排列，
        // 又因为 php 的 array 遍历均有顺序性，因此这里的遍历顺序也是按照从小到大的顺序进行遍历
        foreach ($this->hashMap as $k => $v) {
            if ($k >= $hash) {
                $node = $v;
                break;
            }
        }

        // 将数组的内部指针指向第一个单元，方便下一次使用 current 时，刚好指向环上最小的一个节点
        reset($this->hashMap);

        return $node;
    }

    public function toArray(): array
    {
        return [
            'replicas' => $this->replicas,
            'hash_map' => $this->hashMap,
            'nodes' => $this->nodes
        ];
    }

}