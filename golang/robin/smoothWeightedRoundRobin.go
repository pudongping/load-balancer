package robin

type SmoothWeightedRoundRobin struct {
	servers []*SmoothServer
}

type SmoothServer struct {
	weight          int    // 初始化权重（配置的权重）
	host            string // 节点地址
	currentWeight   int    // 当前权重（每一次挑选节点时都会发生变化）
	effectiveWeight int    // 有效权重，默认值与 weight 相同，对此值的影响直接可以影响到被挑选的概率。值越大，选中次数则越多，反之亦然。
}

func NewSmoothWeightedRoundRobin(servers []*SmoothServer) *SmoothWeightedRoundRobin {

	for _, server := range servers {
		// 初始值时，effectiveWeight = weight
		server.effectiveWeight = server.weight
		// 初始值时，currentWeight = 0
		server.currentWeight = 0
	}

	return &SmoothWeightedRoundRobin{
		servers: servers,
	}
}

// 算法流程：
// 1. 计算出所有节点的当前权重 currentWeight = currentWeight + effectiveWeight
// 2. 选中此时所有节点的 currentWeight 最大的那个节点。如果有多个最大，那么则选中在所有节点中位置最靠前的那个，也就是第一个 currentWeight 最大的节点。
// 3. 被选中的节点的 currentWeight = currentWeight - totalWeight （totalWeight = effectiveWeight1 + effectiveWeight2 + .... + effectiveWeight n）
// 4. 如果需要选择多次时，则重复 1、2、3 步骤
func (swr *SmoothWeightedRoundRobin) getPeer() *SmoothServer {
	if len(swr.servers) == 0 {
		return nil
	}

	// 记录所有节点有效权重之和
	total := 0
	var best *SmoothServer

	for _, peer := range swr.servers {
		// 此时比较的权重为 currentWeight += effectiveWeight
		peer.currentWeight += peer.effectiveWeight
		total += peer.effectiveWeight

		// if peer.effectiveWeight < peer.weight {
		// 	peer.effectiveWeight++
		// }

		// 找出此时所有节点中，权重最大的节点作为最佳节点
		if best == nil || peer.currentWeight > best.currentWeight {
			best = peer
		}
	}

	if best == nil {
		return nil
	}

	// 因为是直接操作指针，因此会直接影响到原切片
	best.currentWeight -= total

	return best
}

func (swr *SmoothWeightedRoundRobin) adjustEffectiveWeight(host string, step int) {
	for _, server := range swr.servers {
		if server.host == host {
			tmp := server.effectiveWeight + step
			if tmp <= server.weight && tmp >= 1 {
				// 会直接影响到原切片
				server.effectiveWeight = tmp
			}
			return
		}
	}
}
