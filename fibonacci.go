
import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	first := 0
	second := 0
	return func() int {
		n := first + second
		fmt.Printf("\tfirst=%d second=%d n=%d\n", first, second, n)
		if n <= 0 {
			second = 1
			n = 0
		} else {
			first = second
			second = n
		}
		return n
	}
}


func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

