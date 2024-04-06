package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
)

// wg is used to wait for the program to finish.
var wg = sync.WaitGroup{}

// channel is used to send information between goroutines
var channel = make(chan *big.Float, iterations)

var desired_decimals, _ = strconv.Atoi(os.Args[1])
var iterations = desired_decimals / 14
var precision = iterations * 50
var const1 = big.NewInt(13_591_409)

func main() {
	sum := new(big.Float).SetPrec(uint(precision))

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			chudnovsky(i)
		}(i)
	}

	wg.Wait()
	for i := 0; i < iterations; i++ {
		sum.Add(sum, <-channel)
	}
	close(channel)

	numerator := big.NewFloat(10_005).SetPrec(uint(precision))
	numerator.Sqrt(numerator)
	numerator.Mul(numerator, big.NewFloat(426_880))
	pi := numerator.Quo(numerator, sum)

	fmt.Println(pi)
}

func chudnovsky(k int) {

	temp := big.NewInt(545_140_134)
	temp.Mul(temp, big.NewInt(int64(k)))
	temp.Add(temp, const1)

	int_numerator := new(big.Int)
	int_numerator = int_numerator.Mul(factorial(6*k), temp)

	int_denominator := new(big.Int)
	int_denominator.Mul(factorial(3*k), new(big.Int).Exp(factorial(k), big.NewInt(3), nil))
	int_denominator.Mul(int_denominator, new(big.Int).Exp(big.NewInt(640_320), big.NewInt(int64(3*k)), nil))

	denominator := new(big.Float).SetInt(int_denominator).SetPrec(uint(precision))
	numerator := new(big.Float).SetInt(int_numerator).SetPrec(uint(precision))

	denominator.Quo(numerator, denominator)

	if k%2 == 0 {
		channel <- denominator
	} else {
		channel <- denominator.Neg(denominator)
	}
}

func factorial(n int) *big.Int {
	if n <= 1 {
		return big.NewInt(1)
	}
	result := big.NewInt(int64(n))
	return result.Mul(result, factorial(n-1))
}
