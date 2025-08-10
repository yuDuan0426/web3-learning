package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// This is a placeholder main function.
	fmt.Println("-------------指针值+10----------")
	num := 5
	increaseByTen(&num)
	fmt.Println("指针值5加10=", num) // Output: 15

	fmt.Println("-------------切片元素*2----------")
	nums := []int{1, 2, 3, 4, 5}
	doubleSliceElements(&nums)
	fmt.Println("切片元素乘以2=", nums) // Output: [2, 4, 6, 8, 10]

	fmt.Println("-------------打印1-10的奇数和偶数----------")
	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(2)
	go printODD(&wg, ch)
	go printEven(&wg, ch)
	ch <- struct{}{} // 发送初始信号以开始打印
	wg.Wait()
	fmt.Println("打印完成")
	close(ch)

	fmt.Println("-------------生产者消费者模式----------")
	ch2 := make(chan int, 10) // 创建一个带缓冲的通道
	var wg2 sync.WaitGroup
	wg2.Add(2)
	go producer(ch2, &wg2) // 启动生产者协程
	go consumer(ch2, &wg2) // 启动消费者协程
	wg2.Wait()             // 等待所有协程完成
	fmt.Println("生产者消费者模式完成")
	fmt.Println("-------------程序结束----------")

	fmt.Println("-------------使用锁保护共享资源----------")
	counter := SafeCounter{}
	var wg3 sync.WaitGroup
	wg3.Add(10) // 启动10个协程
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				counter.Increment() // 每个协程对计数器进行1000次添加
			}
			wg3.Done()
		}()
	}
	wg3.Wait()                                // 等待所有协程完成
	fmt.Println("计数器的值:", counter.getCount()) // 输出计数器的值
	fmt.Println("-------------锁保护共享资源完成----------")

	fmt.Println("-------------使用原子操作保护共享资源----------")
	counterAtomic := SafeCounter{}
	var wg4 sync.WaitGroup
	wg4.Add(10) // 启动10个协程
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				counterAtomic.Increment() // 每个协程对计数器进行1000次添加
			}
			wg4.Done()
		}()
	}
	wg4.Wait()                                      // 等待所有协程完成
	fmt.Println("计数器的值:", counterAtomic.getCount()) // 输出计数器的值
	fmt.Println("-------------原子操作保护共享资源完成----------")
	fmt.Println("-------------程序结束----------")
}

// 编写一个Go程序，定义一个函数，接受一个整数指针作为参数，并将该整数的值增加10。然后在main函数中调用该函数，并打印增加后的值。
func increaseByTen(num *int) {
	*num += 10
}

// 实现一个函数，接受一个整数切点的指针，将切片中的每个元素乘以2。
func doubleSliceElements(nums *[]int) {
	for i := range *nums {
		(*nums)[i] *= 2
	}
}

// 使用go启动两个协成，一个打印1-10的奇数，另一个协成打印2-10的偶数
func printODD(wg *sync.WaitGroup, ch chan struct{}) {
	defer wg.Done()
	for i := 1; i <= 10; i += 2 {
		<-ch // 等待信号
		fmt.Println("奇数:", i)
		ch <- struct{}{} // 发送信号
	}
}

func printEven(wg *sync.WaitGroup, ch chan struct{}) {
	defer wg.Done()
	for i := 2; i <= 10; i += 2 {
		<-ch // 等待信号
		fmt.Println("偶数:", i)
		//ch <- struct{}{} // 发送信号

		if i < 10 { // 如果不是最后一个偶数，发送信号
			ch <- struct{}{}
		}
	}
}

// 实现一个带有缓冲的通道，生产者协程向通道发送100个整数，消费者从通道接收这些整数并打印。使用sync.WaitGroup来确保所有协程完成后再退出程序。
func producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		ch <- i
		fmt.Println("生产者发送:", i)
	}
	close(ch) // 关闭通道，表示不再发送数据
}

// 接收channel参数
func consumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Println("消费者接收:", num)
	}
}

// 实现一个锁，使用sync.Mutex来保护共享资源，确保在多协程环境下的安全访问。每个协程对计数器进行1000次添加，最后输出计数器的值
// type SafeCounter struct {
// 	mu	sync.Mutex // 使用互斥锁保护共享资源
//  count int
// }

// func (c *SafeCounter) Increment() {
// 	c.mu.Lock()         // 锁定
// 	defer c.mu.Unlock() // 解锁
// 	c.count++           // 计数器加1
// }

// func (c *SafeCounter) getCount() int {
// func (c *SafeCounter) getCount() int64 {
// 	c.mu.Lock()         // 锁定
// 	defer c.mu.Unlock() // 解锁
// 	return c.count      // 返回计数器的值
// }

// 实现一个锁，使用原子类来保护共享资源，确保在多协程环境下的安全访问。每个协程对计数器进行1000次添加，最后输出计数器的值
type SafeCounter struct {
	count int64
}

func (c *SafeCounter) Increment() {
	atomic.AddInt64(&c.count, 1) // 计数器加1
}

func (c *SafeCounter) getCount() int64 {
	return atomic.LoadInt64(&c.count) // 返回计数器的值
}
