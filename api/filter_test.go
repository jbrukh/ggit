package api

import (
	"github.com/jbrukh/ggit/util"
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
	util.Assert(t, numbers[0] == 1)
	util.Assert(t, numbers[1] == 3)
	util.Assert(t, numbers[2] == 5)
	util.Assert(t, numbers[3] == 7)
	util.Assert(t, numbers[4] == 9)
	util.Assert(t, len(numbers) == 5)
}

func Test_EvenNumbers(t *testing.T) {
	numbers := applyInt(numbers, evenNumbers)
	util.Assert(t, numbers[0] == 2)
	util.Assert(t, numbers[1] == 4)
	util.Assert(t, numbers[2] == 6)
	util.Assert(t, numbers[3] == 8)
	util.Assert(t, numbers[4] == 10)
	util.Assert(t, len(numbers) == 5)
}

func Test_NotOdd(t *testing.T) {
	numbers := applyInt(numbers, FilterNot(oddNumbers))
	util.Assert(t, numbers[0] == 2)
	util.Assert(t, numbers[1] == 4)
	util.Assert(t, numbers[2] == 6)
	util.Assert(t, numbers[3] == 8)
	util.Assert(t, numbers[4] == 10)
	util.Assert(t, len(numbers) == 5)
}

func Test_NotEven(t *testing.T) {
	numbers := applyInt(numbers, FilterNot(evenNumbers))
	util.Assert(t, numbers[0] == 1)
	util.Assert(t, numbers[1] == 3)
	util.Assert(t, numbers[2] == 5)
	util.Assert(t, numbers[3] == 7)
	util.Assert(t, numbers[4] == 9)
	util.Assert(t, len(numbers) == 5)
}

func Test_Or(t *testing.T) {
	numbers := applyInt(numbers, FilterOr(evenNumbers, oddNumbers))
	for i := 0; i < 10; i++ {
		util.Assert(t, numbers[i] == i+1)
	}
	util.Assert(t, len(numbers) == 10)
}
