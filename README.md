# load-balancer

几种负载均衡调度算法

> 本示例提供了 [php 版本](./php/Robin) 和 [golang 版本](./golang/robin)

## 简单轮询

[示例代码](./php/Robin/SimpleRoundRobin.php)

优点：
- 简单明了

缺点：
- 不支持配置权重

## 加权轮询

[示例代码](./php/Robin/WeightedRoundRobin.php)

优点：
- 能够配置权重  

缺点：
- 虽然支持了权重，但是负载不均衡。当权重偏差越大时，偏移越严重，比如会出现 {c, b, a, a, a, a, a} 情况

参考文献：

- [负载均衡算法 — 轮询](https://www.fanhaobai.com/2018/11/load-balance-round-robin.html)

## 平滑加权轮询（加权动态优先级算法）

算法特点：  
1. 当前节点的当前权重 current_weight = 上一次当前节点的当前权重 current_weight + 当前节点的有效权重 effective_weight
2. 在当前所有节点中，选出此时当前权重 current_weight 最大的那个节点作为最佳节点。如果出现有多个节点 current_weight 同时最大，那么则以第一个 current_weight 最大节点为最佳节点。
3. 选出的最佳节点的当前权重 current_weight 减去所有节点有效权重 effective_weight 之和

> 其实不难发现，直接影响节点选中的因素就是节点的 effective_weight 值，那么 weight 又是做啥的呢？其实 weight 就单纯用于作为初始权重，可以理解为就是固定配置。最初时，effective_weight 就为 weight 的值，作为初始值。算法运算过程中时，我们可以对 effective_weight 的值
> 做变更，达到动态控制权重的效果，减少则降低权重，增加则加大权重，但是增加 effective_weight 时，不要超过 weight 的值，否则则有可能超过了最初的所有节点的总权重之和，这样貌似之前设置的初始权重 weight 就没有太大的意义了。比如说，初始值 weight 为 {a:5, b:1, c:1} 此时 a 节点
> 的权重控制在 5/7 （轮询 7 次时，最多 5 次命中 a 节点），如果一下子 a 的 effective_weight 变成了 6 那么则权重就可能为 6/8 打破了原来初始权重的比例，就不太符合初始权重的设想了。当然，这个也跟自己的实际业务场景有关。

[示例代码](./php/Robin/SmoothWeightedRoundRobin.php)

优点：
- 能够配置权重
- 就算是权重偏差很大，也能够做到尽可能的负载均衡，比如 { a, a, b, a, c, a, a } 情况，均匀程度提升的非常显著

参考文献：

- [nginx 平滑加权轮询源码](https://github.com/nginx/nginx/blob/master/src/http/ngx_http_upstream_round_robin.c#L522)
- [相关算法详见这个 commit](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35)
- [负载均衡算法 — 平滑加权轮询](https://www.fanhaobai.com/2018/12/load-balance-smooth-weighted-round-robin.html)
- [关于验证该算法权重合理性以及平滑性，可以参考这位网友的文章 "nginx平滑的基于权重轮询算法分析"](https://tenfy.cn/2018/11/12/smooth-weighted-round-robin/) （虽然没有看懂证明过程，哈哈😄）
- [Golang实现四种负载均衡算法](https://juejin.cn/post/6871169933150486542)

以下是个人的理解：  

> 下面所说的比较过程指的是代码中的 `$this->servers[$key]['current_weight'] > $best['current_weight']` 这一行。（决定了最终选择哪个节点作为最佳节点 best）


| 节点  | 配置权重 weight | 有效权重 effective_weight（初始值等于 weight） | 当前权重 current_weight （初始值都为 0，往后每次选择节点时 current_weight += effective_weight） |
|-----|----------------------------|-------------------------------------|----------------------------------------------------------------------------|
| a   | 5                          | 5                                   | 0                                                                          | 
| b | 1                          | 1                                   | 0                                                                          | 
| c | 1                          | 1                                   | 0                                                                          |


假设轮询 8 次，因为权重分别为 `{a:5, b:1, c:1}` 那么也就意味着 5+1+1=7 次循环为一个周期，且在一个周期内 a 被选中 5 次，b 被选中 1 次，c 被选中 1 次，第 8 次循环时，又是一个新的周期开始，且选中顺序和第一个周期的选中顺序一致。
代码比较过程如下：

| 轮询次数 | current_weight 变化前（current_weight += effective_weight） | 被选中的节点 best                                 | current_weight 变化后（被选中节点的当前权重 current_weight 要减去所有节点 effective_weight 之和） |
|------|--------------------------------------------------------|---------------------------------------------|---------------------------------------------------------------------------|
| 1    | {a:5=0+5, b:1=0+1, c:1=0+1}                            | 5 最大，则选 a                                   | {a:-2=5-7, b:1, c:1}                                                      |
| 2    | {a:3=-2+5, b:2=1+1, c:2=1+1}                           | 3 最大，则选 a                                   | {a:-4=3-7, b:2, c:2}                                                      |
| 3    | {a:1=-4+5, b:3=2+1, c:3=2+1}                           | 3 最大，同时有两个最大 b 和 c，但是 b 比 c 在所有节点中位置靠前，则选 b | {a:1, b:-4=3-7, c:3}                                                      |
| 4    | {a:6=1+5, b:-3=-4+1, c:4=3+1}                          | 6 最大，则选 a                                   | {a:-1=6-7, b:-3, c:4}                                                     |
| 5    | {a:4=-1+5, b:-2=-3+1, c:5=4+1}                         | 5 最大，则选 c                                   | {a:4, b:-2, c:-2=5-7}                                                     |
| 6    | {a:9=4+5, b:-1=-2+1, c:-1=-2+1}                        | 9 最大，则选 a                                   | {a:2=9-7, b:-1, c:-1}                                                     |
| 7    | {a:7=2+5, b:0=-1+1, c:0=-1+1}                          | 7 最大，则选 a                                   | {a:0=7-7, b:0, c:0}                                                       |
| 8    | {a:5=0+5, b:1=0+1, c:1=0+1}                            | 5 最大，则选 a                                   | {a:-2=5-7, b:1, c:1}                                                      |
