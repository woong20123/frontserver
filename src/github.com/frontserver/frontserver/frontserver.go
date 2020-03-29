package main

// "factored" import statment
import (
	"fmt"
	"math"

	// 새로 정의한 testlib/mymath 라이브러리를 추가해서 수행합니다.
	"testlib/mymath"
)

func main() {
	fmt.Printf("Now you have %g problems.\n", math.Round(math.Sqrt(7)))
	fmt.Printf("mymath Plus(%d)\n", mymath.Plus(10, 10))
	fmt.Printf("mymath Multiple(%d)\n", mymath.Multiple(10, 10))
}
