<?php
/**
 * 简单轮询：
 * 每一个节点都会被一个一个的调度，如果所有的节点都被调度过一次时，那么则从头再开始调度
 * 不支持配置负载
 */
declare(strict_types=1);

class SimpleRoundRobin
{

    /**
     * 所有服务器节点
     *
     * @var array
     */
    private $servers;

    /**
     * 当前被轮询到的服务器索引位置
     *
     * @var int
     */
    private $currentIndex;

    public function __construct(array $servers)
    {
        $this->servers = $servers;
        $this->currentIndex = 0;
    }

    /**
     * 得到服务器节点
     *
     * @return string|null
     */
    public function getPeer(): ?string
    {
        if (! $this->servers) {
            return null;
        }

        $peer = $this->servers[$this->currentIndex];
        $this->currentIndex = ($this->currentIndex + 1) % count($this->servers);
        return $peer;
    }

}