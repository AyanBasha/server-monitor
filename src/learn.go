package main

import (
	"fmt"
	//"runtime"

	//"math"
	//"math/rand"
	"math/cmplx"
	//"math/rand"
	//"math/rand/v2"
)

var stud, stuf, stug bool = true, false, false
var i int = 3

const Pi = 3.14

//:= not available outside functions bc a keyword is required

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	ToAst  complex128 = cmplx.Sqrt(-5 + 12i)
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
	x = sum / 2
	y = sum - x
	return
}

func checker(age, height, shoeSize int) int {
	if v := height / age; v == shoeSize {
		return v
	} else {

	}

	return shoeSize
}

type Vertex struct {
	X int
	Y int
}

type Coordinate struct {
	Lat, Long float64
}

var m map[string]Coordinate

var m1 = map[string]Coordinate{
	"Google": Coordinate{
		1, 2,
	},
	"russel": Coordinate{
		2, 3,
	},
}

var m2 = map[string]Coordinate{
	"People": {40.23456, -72.523},
}

var (
	v1 = Vertex{X: 0}
	v2 = Vertex{1, 1}
	v3 = Vertex{-1, 1}
	p  = Vertex{}
)

/*func Sqrt (x float64) float64 {
	//var z = rand.Float64(int(x)+1)
	z = float64(z)

	for ; float64(z*z) != x; {
		z -= (z*z - x) / (2*z)
	}

}*/

func main() {
	map1 := make(map[string]int)

	map1["Answer"] = 42
	map1["Answer"] = 48

	delete(map1, "Answer")

	//elem, ok := map1["Answer1"]

	m = make(map[string]Coordinate)

	m["House1"] = Coordinate{
		40.67857, 45.6875,
	}

	fmt.Println(m["House1"])

	fmt.Println(v1, v2, v3, p)

	fmt.Println(Vertex{1, 2})

	v1 := Vertex{1, 2}

	p := &v1

	(*p).X = 1e9

	fmt.Println(v1)

	//array

	var array [10]int

	array[0] = 1
	//array 1 is made 0
	array[2] = 2
	array[3] = 3
	fmt.Println(array[0], array[1], array[2])

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	//slice
	var s []int = primes[1:4]
	fmt.Println(s)

	s = s[:0]

	s = s[:4]

	s = s[2:]

	fmt.Println(len(s), cap(s))

	a := make([]int, 5)

	fmt.Println(cap(a))

	b := make([]int, 0, 5)

	fmt.Println(b)

	a = append(a, 83110)

	/*var sum int = 0

	for i := 0; i < 10; i++ {
		sum += i
	}

	fmt.Println(sum)

	fmt.Println("w/o init and post statements:")

	for sum < 1000 {
		sum += sum
	}

	fmt.Println(sum)

	//for {}

	if sum < 0 {
		sum = -1 * sum
	}

	defer fmt.Print("Go runs on")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("MacOS.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.\n", os)
	}

	for i := 0; i < 10; i += 3 {
		defer fmt.Println(i)
	}


	//pointers
	z := 4
	p := &z

	*p = *p/4


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
