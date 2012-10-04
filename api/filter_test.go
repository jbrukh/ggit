package api

import (
	"testing"
)

func applyInt(s []int, f Filter) []int {
	r := make([]int, 0)
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

var oddNumbers Filter = func(i interface{}) bool {
	if i.(int)%2 == 1 {
		return true
	}
	return false
}

var evenNumbers Filter = func(i interface{}) bool {
	if i.(int)%2 == 0 {
		return true
	}
	return false
}
var numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func Test_OddNumbers(t *testing.T) {
	numbers := applyInt(numbers, oddNumbers)
	assert(t, numbers[0] == 1)
	assert(t, numbers[1] == 3)
	assert(t, numbers[2] == 5)
	assert(t, numbers[3] == 7)
	assert(t, numbers[4] == 9)
	assert(t, len(numbers) == 5)
}

func Test_EvenNumbers(t *testing.T) {
	numbers := applyInt(numbers, evenNumbers)
	assert(t, numbers[0] == 2)
	assert(t, numbers[1] == 4)
	assert(t, numbers[2] == 6)
	assert(t, numbers[3] == 8)
	assert(t, numbers[4] == 10)
	assert(t, len(numbers) == 5)
}

func Test_NotOdd(t *testing.T) {
	numbers := applyInt(numbers, FilterNot(oddNumbers))
	assert(t, numbers[0] == 2)
	assert(t, numbers[1] == 4)
	assert(t, numbers[2] == 6)
	assert(t, numbers[3] == 8)
	assert(t, numbers[4] == 10)
	assert(t, len(numbers) == 5)
}

func Test_NotEven(t *testing.T) {
	numbers := applyInt(numbers, FilterNot(evenNumbers))
	assert(t, numbers[0] == 1)
	assert(t, numbers[1] == 3)
	assert(t, numbers[2] == 5)
	assert(t, numbers[3] == 7)
	assert(t, numbers[4] == 9)
	assert(t, len(numbers) == 5)
}

func Test_Or(t *testing.T) {
	numbers := applyInt(numbers, FilterOr(evenNumbers, oddNumbers))
	for i := 0; i < 10; i++ {
		assert(t, numbers[i] == i+1)
	}
	assert(t, len(numbers) == 10)
}
