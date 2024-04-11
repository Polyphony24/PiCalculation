package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"
)

// wg is used to wait for the program to finish.
var wg = sync.WaitGroup{}

// channel is used to send information between goroutines
var channel = make(chan *big.Float, iterations)

// input is desired_decimals i.e. ./Chudnovsky 100 will print out 100 decimals
var desired_decimals, _ = strconv.Atoi(os.Args[1])

// each iteration makes minimum 14 decimals
var iterations = desired_decimals / 14

// this is just a number that seems to work
// may be able to lower it for speedup in some situations (not sure)
var precision = iterations * 50

// const1 is just a value used in each sum
// global variable so we don't have to calculate it everytime
var const1 = big.NewInt(13_591_409)

var sum = new(big.Float).SetPrec(uint(precision))

func main() {

	start := time.Now()
	// sum is the sum which we are parallelizing

	for i := 0; i < iterations; i++ {
		// add a goroutine to the waitgroup
		// and call the go routine
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			chudnovsky(i)
		}(i)
	}

	for i := 0; i < iterations; i++ {
		sum.Add(sum, <-channel)
	}

	// wait until all iterations are done
	wg.Wait()
	// close the channel
	close(channel)

	numerator := big.NewFloat(10_005).SetPrec(uint(precision))
	numerator.Sqrt(numerator)
	numerator.Mul(numerator, big.NewFloat(426_880))

	pi := numerator.Quo(numerator, sum)
	pi.Neg(pi)
	endTime := time.Now().Sub(start)

	//fmt.Println(pi)
	fmt.Println(endTime)
}

func chudnovsky(k int) {
	// start with big.Ints because ints are faster and they are all integers in the fraction
	temp := big.NewInt(545_140_134)
	temp.Mul(temp, big.NewInt(int64(k)))
	temp.Add(temp, const1)

	// numerator of the sum
	// (6k)!(545,140,134k + 13,591,409)
	numerator := new(big.Int)
	numerator = numerator.Mul(factorial(6*k), temp)

	// denominator of the sum
	// (3k)!(k!)^3(640320)^3k
	denominator := new(big.Int)
	denominator.Mul(factorial(3*k), new(big.Int).Exp(factorial(k), big.NewInt(3), nil))
	denominator.Mul(denominator, new(big.Int).Exp(big.NewInt(640_320), big.NewInt(int64(3*k)), nil))

	fraction := new(big.Float).Quo(new(big.Float).SetPrec(uint(precision)).SetInt(numerator), new(big.Float).SetPrec(uint(precision)).SetInt(denominator))

	// return the result to the channel
	if k%2 == 1 {
		fraction.Neg(fraction)
	}

	channel <- fraction
}

// pretty fuckin self explanatory
func factorial(n int) *big.Int {
	result := big.NewInt(1)
	for i := 1; i < n+1; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}
