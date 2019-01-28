package main

import "golang.org/x/tour/tree"
import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	/*fmt.Printf("\tTree=%#v\n", t)
	fmt.Printf("New Value=%d\n", t.Value)*/
	ch <- t.Value
	if t.Right != nil {
		go Walk(t.Right, ch)
	}
	if t.Left != nil {
		go Walk(t.Left, ch)
	}
	if t.Right == nil && t.Left == nil {
		/*fmt.Printf("Right=%#v Left=%#v\n", t.Right, t.Left)
		fmt.Println("closing channel")*/
		close(ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	fmt.Printf("t1.Value=%d t2.Value=%d\n", t1.Value, t2.Value)
	if t1.Value == t2.Value {
		if t1.Right != nil && t2.Right != nil {
			fmt.Println("right not nil")
			return Same(t1.Right, t2.Right)
		}
		if t1.Left != nil && t2.Left != nil {
			fmt.Println("left not nil")
			return Same(t1.Left, t2.Left)
		}
		return true
	} else {
		return false
	}
}

func main() {
	t1 := tree.New(1)
	t2 := tree.New(2)
	treesAreSame := Same(t1, t1)
	fmt.Printf("treesAreSame=%t\n", treesAreSame)
	fmt.Printf("t1=%#v\n", t1)
	fmt.Printf("t2=%#v\n", t2)

	ch := make(chan int)
	go Walk(t1, ch)
	for i := range ch {
		fmt.Println(i)
	}
}
