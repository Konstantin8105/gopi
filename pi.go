// Calculate π by Leibniz formula
// See: https://en.wikipedia.org/wiki/Leibniz_formula_for_%CF%80
package pi

import (
	"fmt"
	"math/big"
	"sync"
)

type PiService struct {
	m     sync.Mutex
	cStop chan struct{}

	// amount of iterations
	iter *big.Int

	// result of calculation and for infinite iteration return {pi/4}.
	result *big.Float

	// denominator
	den *big.Int
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
	var (
		one    *big.Float = big.NewFloat(1)
		oneInt *big.Int   = big.NewInt(1)
		twoInt *big.Int   = big.NewInt(2)
	)
	go func() {
		var minus bool = true
		cIncrement := make(chan big.Float)
		go func() {
			defer close(cIncrement)
			for {
				select {
				case <-s.cStop:
					return
				default:
				}
				// calculate next increment
				var next big.Float
				next.Quo(one, new(big.Float).SetInt(s.den))
				if minus {
					next.Neg(&next)
				}

				// lock for next increment
				s.m.Lock()
				s.iter.Add(s.iter, oneInt)
				s.den.Add(s.den, twoInt)
				minus = !minus
				cIncrement <- next
				s.m.Unlock()
			}
		}()

		// add increment to result
		for i := range cIncrement {
			s.result.Add(s.result, &i)
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
