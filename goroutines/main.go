// package main

// import (
// 	"fmt"
// )

// func numbers(ch chan int) {
// 	for i := 0; i < 8; i++ {
// 		ch <- i
// 	}
// 	// close(ch)
// }

// func multiply(ch1 chan int, ch2 chan int) {
// 	num := <-ch1
// 	ch2 <- num * 2
// }

// func main() {
// 	ch1 := make(chan int)
// 	ch2 := make(chan int)
// 	go numbers(ch1)
// 	for i := 0; i < 8; i++ {
// 		go multiply(ch1, ch2)
// 		select {
// 		case m1 := <-ch1:
// 			fmt.Printf("channel 1: %v\n", m1)
// 		case m2 := <-ch2:
// 			fmt.Printf("channel 2: %v\n", m2)
// 		}
// 	}
// }

package main

import (
	"fmt"
	"sync"
)

func numbers(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 8; i++ {
		fmt.Println(i)
	}
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	fmt.Println("Main starting...")
	go numbers(wg)
	wg.Wait()
	fmt.Println("Main finished...")
}
