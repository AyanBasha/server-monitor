package main

import (
	"fmt"
	"math"
	"math/rand"
	"math/cmplx"
)

var stud, stuf, stug bool = true, false, false
var i int = 3
const Pi = 3.14
//:= not available outside functions bc a keyword is required

var (
	ToBe bool = false
	MaxInt uint64 = 1 << 64 -1
	ToAst complex128 = cmplx.Sqrt(-5 + 12i)
)

func add(x int, y int) int {
	return x + y
}

func subtract(x, y int) int {
	return x - y
}

func swap(x, y int) (int, int) {
	return y, x
}

func split(sum int) (x, y int) {
	x = sum/2
	y = sum - x
	return
}

func main() {
	var x, y int = 3, 5
	var f float64 = math.Sqrt(float64(x*x+y*y))
	var z uint = uint(f)
	fmt.Println(x, y, z)
	fmt.Printf("Type: %T Value: %v \n", ToAst, ToAst)
	fmt.Println(stud, stuf, stug, i)
	a, b := split(17)
	fmt.Println(a, b)
	a, b = swap(54, 27)
	fmt.Println(a, b)
	fmt.Println(swap(786, 2836))
	fmt.Println(subtract(786, 2836))
	fmt.Println(add(4, 2))
	fmt.Println("My favorite number is", rand.Intn(10))
	fmt.Printf("Now you have %g problems.\n", math.Sqrt(7))
	fmt.Println(math.Pi)
}
