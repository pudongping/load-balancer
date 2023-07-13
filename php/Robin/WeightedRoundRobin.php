<?php
/**
 * 加权轮询：
 * 每一个节点都会被调度
 * 支持配置负载，但是负载不太均衡
 *
 */
declare(strict_types=1);

class WeightedRoundRobin
{

    /**
     * 所有服务器节点
     *
     * @var array
     */
    private $servers;

    /**
     * 当前权重
     *
     * @var int
     */
    private $currentWeight;

    /**
     * 所有服务器节点中最大的权重
     *
     * @var int
     */
    private $maxWeight;

    /**
     * 当前被挑中的节点的索引位置
     *
     * @var int
     */
    private $currentIndex;

    /**
     * 总共有多少个节点服务器
     *
     * @var int
     */
    private $total;

    public function __construct(array $servers)
    {
        $this->servers = $servers;
        $this->currentWeight = 0;
        $this->maxWeight = (int)max(1, ...array_column($this->servers, 'weight'));
        $this->currentIndex = -1;
        $this->total = count($this->servers);
    }

    /**
     * 选出节点
     *
     * @return array|null
     */
    public function getPeer(): ?array
    {
        if (! $this->servers) {
            return null;
        }

        $i = $this->currentIndex;

        while (true) {
            $i = ($i + 1) % $this->total;

            if (0 === $i) {  // 表示已经走完一轮了
                $this->currentWeight -= $this->calculateGCD();

                if ($this->currentWeight <= 0) {
                    $this->currentWeight = $this->maxWeight;
                }
            }

            if ($this->servers[$i]['weight'] >= $this->currentWeight) {
                $this->currentIndex = $i;  // 记录当前的被选中 host 的索引位置

                return $this->servers[$this->currentIndex];  // 被挑选出的 host
            }
        }

    }

    /**
     * 获取最大公约数
     *
     * @return int
     */
    private function calculateGCD(): int
    {
        $weights = array_column($this->servers, 'weight');
        $gcd = $weights[0];

        foreach ($weights as $weight) {
            $gcd = $this->calculateGCDRecursive($gcd, $weight);
        }

        return $gcd;
    }

    /**
     * 求两数的最大公约数
     *
     * 欧几里德算法：
     * 1. 将较大的数除以较小的数，得到余数。
     * 2. 将较小的数除以余数，再次得到余数。
     * 3. 重复以上步骤，直到余数为 0。此时，被除数就是两个数的最大公约数。
     *
     * @param int $a
     * @param int $b
     * @return int
     */
    private function calculateGCDRecursive(int $a, int $b): int
    {
        // a=25 b=10
        // a=10 b= 25%10 = 5
        // a=5  b= 10%5 = 0
        // result = 5
        return ($b === 0) ? $a : $this->calculateGCDRecursive($b, $a % $b);
    }

}