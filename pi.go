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
	four   = big.NewInt(4)
)

// calculate next increment
func (s *Service) calculate(den *big.Int, cIncrement chan<- *big.Float) {
	var next big.Float
	next.SetInt(den)
	next.Quo(one, &next)
	cIncrement <- &next
}

// Start pi-service
func (s *Service) Start() {
	amountWorkers := 10
	cIncrement := make(chan *big.Float, amountWorkers)
	cDen := make(chan *big.Int, amountWorkers)

	// denominator
	den := big.NewInt(-3)
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < amountWorkers; i++ {
			wg.Add(1)
			go func() {
				for d := range cDen {
					// Step : { 1/(-3) }
					s.calculate(d, cIncrement)

					// Step : { 1/(+5) }
					d.Neg(d)
					d.Add(d, two)
					s.calculate(d, cIncrement)
				}
				wg.Done()
			}()
		}
		go func() {
			wg.Wait()
			close(cIncrement)
		}()
		for {
			select {
			case <-s.cStop:
				close(cDen)
				return
			default:
			}
			copy := new(big.Int).Set(den)
			cDen <- den

			den = copy
			den.Sub(den, four)

			s.m.Lock()
			s.iter.Add(s.iter, two)
			s.m.Unlock()
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
	return fmt.Sprintf("%s", out.Text('f', 62))
}

// Stop pi-service
func (s *Service) Stop() {
	s.cStop <- struct{}{}
}
