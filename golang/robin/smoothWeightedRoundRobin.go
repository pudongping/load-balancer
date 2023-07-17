package robin

type SmoothWeightedRoundRobin struct {
	servers      []*SmoothServer
	currentIndex int
	best         *SmoothServer
}

type SmoothServer struct {
	weight          int
	host            string
	currentWeight   int
	effectiveWeight int
}

func NewSmoothWeightedRoundRobin(servers []*SmoothServer) *SmoothWeightedRoundRobin {

	for _, server := range servers {
		// 初始值时，effectiveWeight = weight
		server.effectiveWeight = server.weight
		// 初始值时，currentWeight = 0
		server.currentWeight = 0
	}

	return &SmoothWeightedRoundRobin{
		servers:      servers,
		currentIndex: -1,
		best:         nil,
	}
}

func (swr *SmoothWeightedRoundRobin) getPeer() *SmoothServer {
	if len(swr.servers) == 0 {
		return nil
	}

	// 记录所有节点有效权重之和
	total := 0

	for key, peer := range swr.servers {
		peer.currentWeight += peer.effectiveWeight
		total += peer.effectiveWeight

		// if peer.effectiveWeight < peer.weight {
		// 	peer.effectiveWeight++
		// }

		if swr.best == nil || peer.currentWeight > swr.best.currentWeight {
			swr.best = peer
			swr.currentIndex = key
		}
	}

	if swr.best == nil {
		return nil
	}

	swr.best.currentWeight -= total

	return swr.best
}

func (swr *SmoothWeightedRoundRobin) adjustEffectiveWeight(host string, step int) {
	for _, server := range swr.servers {
		if server.host == host {
			tmp := server.effectiveWeight + step
			if tmp <= server.weight && tmp >= 1 {
				server.effectiveWeight = tmp
			}
			return
		}
	}
}
