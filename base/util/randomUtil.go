package util

// GetRandomNums 随机范围 [1,max]
func GetRandomNums(max int64, length int, except uint64) []uint64 {
	var res []uint64
	var has bool
	var loop int

	for len(res) < length {
		if loop > 10000 {
			break
		}
		loop++
		t := uint64(randSource.Int63n(max)) + 1
		if t == except {
			continue
		}
		for _, v := range res {
			if v == t {
				has = true
			}
		}
		if !has {
			res = append(res, t)
		}
		has = false
	}
	return res
}

// RandRange
func RandRange(start int64, end int64) int64 {
	if start > end {
		return 0
	}

	return start + randSource.Int63n(end-start)
}

// WeightRandom 加权随机算法
func WeightRandom(list map[interface{}]float32) interface{} {
	if len(list) == 0 {
		return nil
	}
	var weight float32
	for _, v := range list {
		weight += v
	}

	index := randSource.Int63n(int64(weight))

	var curWeight float32
	for k, v := range list {
		curWeight += v
		if index < int64(curWeight) {
			return k
		}
	}

	return nil
}

// WeightRandomUint32 加权随机算法
func WeightRandomUint32(list map[uint32]float32) uint32 {
	if len(list) == 0 {
		return 0
	}
	var weight float32
	for _, v := range list {
		weight += v
	}

	index := randSource.Int63n(int64(weight))

	var curWeight float32
	for k, v := range list {
		curWeight += v
		if index < int64(curWeight) {
			return k
		}
	}

	return 0
}
