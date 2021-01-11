package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkRoundRobin(rb []int, weights []int, gcd int) bool {
	ret := make(map[int]int)
	for _, node := range rb {
		ret[node]++
	}

	invalidNode := 0
	for node, weight := range weights {
		if weight == 0 {
			invalidNode++
			continue
		}
		if v, ok := ret[node]; !ok {
			return false
		} else if v != weight/gcd {
			return false
		}
	}

	if len(ret) != len(weights)-invalidNode {
		return false
	}
	return true
}

func TestBalancer(t *testing.T) {
	weight := []int{2, 1, 4}
	gcd := gcd(weight)
	assert.Equal(t, 1, gcd)

	b := newBalancer(weight, 4)
	assert.Equal(t, true, checkRoundRobin(b.roundRobinQ, weight, gcd))
}

func TestBalancerA(t *testing.T) {
	weight := []int{2, 2, 2, 4}
	gcd := gcd(weight)
	assert.Equal(t, 2, gcd)

	b := newBalancer(weight, 4)
	assert.Equal(t, true, checkRoundRobin(b.roundRobinQ, weight, gcd))
}

func TestBalancerB(t *testing.T) {
	weight := []int{1, 2}
	gcd := gcd(weight)
	assert.Equal(t, 1, gcd)

	b := newBalancer(weight, 3)
	assert.Equal(t, true, checkRoundRobin(b.roundRobinQ, weight, gcd))
	assert.Equal(t, 0, b.lastIndex)

	for i := 0; i < 10; i++ {
		_, err := b.next()
		assert.Equal(t, nil, err)
		assert.Equal(t, (i+1)%len(b.roundRobinQ), b.lastIndex)
	}
}
