# load-balancer

几种负载均衡调度算法

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

## 平滑加权轮询

[示例代码](./php/Robin/SmoothWeightedRoundRobin.php)

优点：
- 能够配置权重
- 就算是权重偏差很大，也能够做到尽可能的负载均衡，比如 { a, a, b, a, c, a, a } 情况

参考文献：

- [nginx 平滑加权轮询源码](https://github.com/nginx/nginx/blob/master/src/http/ngx_http_upstream_round_robin.c#L522)
- [相关算法详见这个 commit](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35)
- [负载均衡算法 — 平滑加权轮询](https://www.fanhaobai.com/2018/12/load-balance-smooth-weighted-round-robin.html)
- [关于验证该算法权重合理性以及平滑性，可以参考这位网友的文章 "nginx平滑的基于权重轮询算法分析"](https://tenfy.cn/2018/11/12/smooth-weighted-round-robin/) （虽然没有看懂，哈哈😄）

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
