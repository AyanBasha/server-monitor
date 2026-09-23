package main

import (
	"fmt"
	//"math"
	//"math/rand"
	"math/cmplx"
	"math/rand"
	"math/rand/v2"
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

func checker(age, height, shoeSize int) int {
	if v:=height/age; v == shoeSize {
		return v
	} else {

	}

	return shoeSize
}

func Sqrt (x float64) float64 {
	var z = rand.Float64(int(x)+1)
	z = float64(z)

	for ; float64(z*z) != x; {
		z -= (z*z - x) / (2*z)
	} 
	
}

func main() {
	var sum int = 0
	
	for i := 0; i < 10; i++ {
		sum += i
	}

	fmt.Println(sum)

	fmt.Println("w/o init and post statements:")

	for ; sum < 1000; {
		sum += sum
	}

	fmt.Println(sum)

	//for {}

	if sum < 0 {
		sum = -1*sum
	}

	/*var x, y int = 3, 5
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
	fmt.Println(math.Pi)*/
}
