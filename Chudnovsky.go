package main

import (
	"fmt"
	"math/big"
	"sync"
)

var const1 = big.NewFloat(13_591_409)

// wg is used to wait for the program to finish.
var wg = sync.WaitGroup{}

// channel is used to send information between goroutines
var channel = make(chan *big.Float)

var iterations = 2500
var precision = 40000

func main() {
	pi := new(big.Float).SetPrec(uint(precision))
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go chudnovsky(i)
	}

	for i := 0; i < iterations; i++ {
		pi.Add(pi, <-channel)
	}

	wg.Wait()

	//fmt.Println(pi)

	numerator := big.NewFloat(10_005).SetPrec(uint(precision))
	numerator.Sqrt(numerator)
	numerator.Mul(numerator, big.NewFloat(426_880))
	pi = numerator.Quo(numerator, pi)
	
	fmt.Printf("%1.15000%d", pi)
}

func chudnovsky(i int) {
	defer wg.Done()

	temp := big.NewFloat(545_140_134).SetPrec(uint(precision))
	temp.Mul(temp, big.NewFloat(float64(i)))
	temp.Add(temp, const1)

	numerator := new(big.Float).SetPrec(uint(precision))
	numerator = numerator.Mul(factorial(6*i), temp)

	denominator := new(big.Float).SetPrec(uint(precision))
	denominator.Mul(factorial(3*i), cube(factorial(i)))
	denominator.Mul(denominator, power(640_320, 3*i))
	denominator.Quo(numerator, denominator)

	if i%2 == 0 {
		channel <- denominator
	} else {
		channel <- denominator.Neg(denominator)
	}
}

func factorial(n int) *big.Float {
	if n <= 1 {
		return big.NewFloat(1).SetPrec(uint(precision))
	}
	result := big.NewFloat(float64(n)).SetPrec(uint(precision))
	return result.Mul(result, factorial(n-1))
}

func cube(x *big.Float) *big.Float {
	// Create a new big.Float for the cube
	cubed := new(big.Float).SetPrec(uint(precision))

	// Cube the value by multiplying it by itself three times
	cubed.Mul(x, x)
	cubed.Mul(cubed, x)

	return cubed
}

func power(x float64, y int) *big.Float {
	z := big.NewFloat(x).SetPrec(uint(precision))
	result := big.NewFloat(1).SetPrec(uint(precision))
	for y > 0 {
		result.Mul(result, z)
		y--
	}
	return result
}
