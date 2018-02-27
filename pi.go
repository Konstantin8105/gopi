// Package pi calculate π by Leibniz formula
// See: https://en.wikipedia.org/wiki/Leibniz_formula_for_%CF%80
package pi

import (
	"fmt"
	"math/big"
	"sync"
)

// Service is base struct of calculation
type Service struct {
	m     sync.Mutex
	cStop chan struct{}

	// amount of iterations
	iter *big.Int

	// result of calculation and for infinite iteration return {pi/4}.
	result *big.Float
}

// NewService create a new service for calculate number of π(Pi)
func NewService() *Service {
	return &Service{
		cStop:  make(chan struct{}),
		iter:   big.NewInt(0),
		result: big.NewFloat(1),
	}
}

var (
	one    = big.NewFloat(1)
	oneInt = big.NewInt(1)
	two    = big.NewInt(2)
)

// calculate next increment
func (s *Service) calculate(den *big.Int, cIncrement chan<- *big.Float) {
	var next big.Float
	next.SetInt(den)
	next.Quo(one, &next)
	cIncrement <- &next
}

// Prepare for next iteration
func (s *Service) prepare(den *big.Int, minus *bool) {
	if (*minus && den.Sign() > 0) ||
		(!*minus && den.Sign() < 0) {
		den.Neg(den)
	}

	if den.Sign() > 0 {
		den.Add(den, two)
	} else {
		den.Sub(den, two)
	}

	*minus = !*minus
}

// Start pi-service
func (s *Service) Start() {
	minus := false
	cIncrement := make(chan *big.Float)
	den := *big.NewInt(-3)
	go func() {
		defer close(cIncrement)
		// denominator
		for {
			select {
			case <-s.cStop:
				return
			default:
			}

			// calculation
			s.calculate(&den, cIncrement)
			// prepare for next iteration
			s.prepare(&den, &minus)

			s.iter.Add(s.iter, oneInt)
		}
	}()

	go func() {
		// add increment to result
		for i := range cIncrement {
			s.m.Lock()
			s.result.Add(s.result, i)
			s.m.Unlock()
		}
	}()
}

// Result return result of calculation Pi
func (s *Service) Result() string {
	s.m.Lock()
	defer s.m.Unlock()
	var out big.Float
	out.Mul(s.result, big.NewFloat(4))
	return fmt.Sprintf("%s", out.String())
}

// Stop pi-service
func (s *Service) Stop() {
	s.cStop <- struct{}{}
}
