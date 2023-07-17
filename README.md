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
2. 每一个节点的当前权重 effective_weight 和上一个最佳节点的当前权重 effective_weight 做比较，如果权重值大于上一个最佳节点的当前权重时，则此节点为最佳节点（一直会比较到最后一个节点，如果有多个节点符合条件时，那么则会以最后一个符合条件的节点作为最佳节点）
3. 选出的最佳节点的当前权重 effective_weight 减去所有节点有效权重之和

> 其实不难发现，直接影响节点的选中的因素就是节点的 effective_weight 值，那么 weight 又是做啥的呢？其实 weight 就单纯用于作为初始权重，可以理解为就是固定配置。最初时，effective_weight 就为 weight 的值，作为初始值。算法运算过程中时，我们可以对 effective_weight 的值
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

> 下面所说的比较过程指的是代码中的 `$this->servers[$key]['current_weight'] > $this->best['current_weight']` 这一行。（决定了选择 best）


| 节点  | 配置权重 weight | 有效权重 effective_weight（初始值等于 weight） | 当前权重 current_weight （每次选择节点时 current_weight += effective_weight） |
|-----|----------------------------|-------------------------------------|------------------------------------------------------------------|
| a   | 5                          | 5                                   | 0                                                                | 
| b | 1                          | 1                                   | 0                                                                | 
| c | 1                          | 1                                   | 0                                                                |

第一轮：

初始权重 current_weight 都为 0  
0  0  0  
比较过程：  
key=0 ==> 5=0+5 >    
key=1 ==> 1=0+1 > 5  
key=2 ==> 1=0+1 > 5  

此时 `best = null` 因此此时第一轮的第一次循环就是第一个节点 a，且又因为 **1 > 5 和 1 > 5** 都不成立，因此此时 best 就为 a  
选择后各节点 current_weight 分别为：   
-2  1  1

---

第二轮：

key=0 ==> 3=-2+5 > -2  
key=1 ==> 2=1+1  > 3  
key=2 ==> 2=1+1  > 3    

此时只有一个条件符合 **3 > -2**，因此此时 best 为 a  
选择后各节点 current_weight 分别为：   
-4  2  2  

---

第三轮：

key=0 ==> 1=-4+5 > -4  
key=1 ==> 3=2+1  > 1  
key=2 ==> 3=2+1  > 3  

此时有两个条件 **1 > -4 和 3 > 1** 成立，但因为此时是循环，因此后面会覆盖前面的，因此此时 best 就为 b  
选择后各节点 current_weight 分别为：  
1  -4  3  

---

第四轮：

key=0 ==> 6=1+5   > -4  
key=1 ==> -3=-4+1 > 6  
key=2 ==> 4=3+1   > 6  

此时 **6 > -4** 成立，因此此时 best 就为 a  
选择后各节点 current_weight 分别为：  
-1  -3  4  

---

第五轮：

key=0 ==> 4=-1+5  > -1  
key=1 ==> -2=-3+1 > 4  
key=2 ==> 5=4+1   > 4  

此时有两个条件 **4 > -1 和 5 > 4** 成立，但因为此时是循环，因此后面会覆盖前面的，因此此时 best 就为 c  
选择后各节点 current_weight 分别为：   
4  -2  -2  

---

第六轮：

key=0 ==> 9=4+5 > -2  
key=1 ==> -1=-2+1 > 9  
key=2 ==> -1=-2+1 > 9  

此时 **9 > -2** 成立，因此此时 best 就为 a   
选择后各节点 current_weight 分别为：  
2  -1  -1  

---

第七轮：

key=0 ==> 7=2+5 > 2  
key=1 ==> 0=-1+1 > 7  
key=2 ==> 0=-1+1 > 7  

此时 **7 > 2** 成立，因此此时 best 就为 a  
选择后各节点 current_weight 分别为：  
0  0  0  

---

第八轮：

key=0 ==> 5=0+5 > 0  
key=1 ==> 1=0+1 > 5  
key=2 ==> 1=0+1 > 5  

此时为第 2 个轮询周期开始，当前权重 current_weight 恢复到初始状态，此时 **5 > 0** 成立，因此此时 best 就为 a  
选择后各节点 current_weight 分别为：  
-2  1  1   
