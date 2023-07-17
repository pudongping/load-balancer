package robin

type WeightedRoundRobin struct {
	servers       []Server
	currentWeight int
	maxWeight     int
	currentIndex  int
	total         int
}

type Server struct {
	weight int
	host   string
}

func NewWeightedRoundRobin(servers []Server) *WeightedRoundRobin {

	// 获取所有节点中最大的权重
	maxWeight := 1
	for _, s := range servers {
		if s.weight > maxWeight {
			maxWeight = s.weight
		}
	}

	return &WeightedRoundRobin{
		servers:       servers,
		currentWeight: 0,
		maxWeight:     maxWeight,
		currentIndex:  -1,
		total:         len(servers),
	}
}

func (rr *WeightedRoundRobin) GetPeer() *Server {
	if rr.total == 0 {
		return nil
	}

	i := rr.currentIndex

	for {
		i = (i + 1) % rr.total

		if i == 0 { // 表示已经走完一轮了
			rr.currentWeight -= rr.calculateGCD()

			if rr.currentWeight <= 0 {
				rr.currentWeight = rr.maxWeight
			}
		}

		if rr.servers[i].weight >= rr.currentWeight {
			rr.currentIndex = i // 记录当前的被选中 host 的索引位置

			return &rr.servers[rr.currentIndex] // 被挑选出的 host
		}
	}
}

func (rr *WeightedRoundRobin) calculateGCD() int {
	weights := make([]int, len(rr.servers))
	for i, s := range rr.servers {
		weights[i] = s.weight
	}

	gcd := weights[0]

	for _, weight := range weights {
		gcd = rr.calculateGCDRecursive(gcd, weight)
	}

	return gcd
}

func (rr *WeightedRoundRobin) calculateGCDRecursive(a, b int) int {
	if b == 0 {
		return a
	}
	return rr.calculateGCDRecursive(b, a%b)
}
