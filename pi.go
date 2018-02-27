// Calculate π by Leibniz formula
// See: https://en.wikipedia.org/wiki/Leibniz_formula_for_%CF%80
package pi

import (
	"fmt"
	"math/big"
	"sync"
)

type PiService struct {
	m      sync.Mutex
	cStop  chan struct{}
	iter   *big.Int // amount of iterations
	result *big.Float
	den    *big.Int // denominator
}

// NewService create a new service for calculate number of π(Pi)
func NewService() *PiService {
	return &PiService{
		cStop:  make(chan struct{}),
		iter:   big.NewInt(0),
		result: big.NewFloat(1),
		den:    big.NewInt(3),
	}
}

// Start pi-service
func (s *PiService) Start() {
	one := big.NewFloat(1)
	oneInt := big.NewInt(1)
	go func() {
		var minus bool = true
		for {
			select {
			case <-s.cStop:
				return
			default:
			}
			s.m.Lock()
			var next big.Float
			next.Quo(one, new(big.Float).SetInt(s.den))
			if minus {
				next.Neg(&next)
			}
			s.result.Add(s.result, &next)
			s.den.Add(s.den, big.NewInt(2))
			minus = !minus
			s.m.Unlock()
			s.iter.Add(s.iter, oneInt)
		}
	}()
}

// Result return result of calculation Pi
func (s *PiService) Result() string {
	s.m.Lock()
	defer s.m.Unlock()
	var out big.Float
	out.Mul(s.result, big.NewFloat(4))
	return fmt.Sprintf("%s", out.String())
}

// Stop pi-service
func (s *PiService) Stop() {
	s.cStop <- struct{}{}
}
