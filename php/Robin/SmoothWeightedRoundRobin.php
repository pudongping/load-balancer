<?php
/**
 * 平滑加权轮询：
 * 每一个节点都会被调度
 * 支持配置负载，且负载较均匀
 *
 * 其中有三个比较关键的点：
 * 1. weight（初始权重、权重配置）：配置文件中设置的权重值，是定值，在整个选择过程中是不会改变的
 * 2. current_weight（当前权重）：后端服务器节点的当前权重值，初始值等于 0，在每轮选择中，该值为 current_weight += effective_weight
 * 3. effective_weight（变化权重值、有效权重）：初始值等于 weight（配置权重），用于动态调整服务器被选择的概率，
 * 即当被选中的服务器出现了故障的时候，该服务器对应的 effective_weight 就会减小，如果故障得以恢复时，
 * 则可通过增加 effective_weight 增加权重，但是不要超过 weight，否则将超过了初始的权重配置，即 effective_weight 永远只能小于或者等于 weight
 *
 * 以下代码摘抄自 nginx 源码： https://github.com/nginx/nginx/blob/master/src/http/ngx_http_upstream_round_robin.c#L522
 * 相关算法详见： https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
 *
 * Created by PhpStorm
 * User: Alex
 * Date: 2023-07-13 00:27
 */
declare(strict_types=1);

class SmoothWeightedRoundRobin
{

    /**
     * 所有服务器节点
     *
     * @var array
     */
    public $servers;

    /**
     * 当前被挑中的节点的索引位置
     *
     * @var int
     */
    private $currentIndex;

    /**
     * 被挑中的服务器节点（其实已知 currentIndex 就可得 best，但是考虑到尽可能的和源码相似，故作保留）
     *
     * @var array
     */
    private $best;

    public function __construct(array $servers)
    {
        $this->servers = $servers;
    }

    /**
     * 挑选出最合适的节点
     *
     * @return array|null
     */
    public function getPeer(): ?array
    {
        if (! $this->servers) {
            return null;
        }

        $total = 0;  // 所有的有效权重之和

        foreach ($this->servers as $key => $peer) {
            // 当前节点的当前权重 = 当前权重 + 有效权重
            $this->servers[$key]['current_weight'] += $this->servers[$key]['effective_weight'];
            // 记录所有节点有效权重之和
            $total += $this->servers[$key]['effective_weight'];

            // 这里可以直接通过 effective_weight 去影响当前节点的当前权重，达到可降级或者升级的效果
            // if ($this->servers[$key]['effective_weight'] < $this->servers[$key]['weight']) {
            //     $this->servers[$key]['effective_weight']++;
            // }

            // echo("key={$key} ==> {$this->servers[$key]['current_weight']} > {$this->best['current_weight']} \r\n");

            if ($this->best === null || $this->servers[$key]['current_weight'] > $this->best['current_weight']) {
                $this->best = $this->servers[$key];
                $this->currentIndex = $key;
            }
        }

        if ($this->best === null) {
            return null;
        }

        $this->servers[$this->currentIndex]['current_weight'] -= $total;
        $this->best['current_weight'] -= $total;

        // echo "\r\n";
        return $this->best;
    }

    /**
     * 调整指定节点的有效权重
     *
     * @param string $host 需要调整有效权重的节点
     * @param int $step 调整步长，可正可负
     * @return void
     */
    public function adjustEffectiveWeight(string $host, int $step)
    {
        $index = null;
        foreach ($this->servers as $k => $server) {
            if ($server['host'] === $host) {
                $index = $k;
                break;
            }
        }

        ! is_null($index) && $this->servers[$index]['effective_weight'] += $step;
    }

}