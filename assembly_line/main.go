package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Employee struct {
	ID        int
	Processed int // 處理多少物品
}

func (e *Employee) Work(wg *sync.WaitGroup, jobs <-chan Item, logCh chan<- string) {
	defer wg.Done()

	for item := range jobs {
		start := time.Now()
		logCh <- fmt.Sprintf("員工 %d 開始處理 %s", e.ID, item.Name())

		item.Process()

		elapsed := time.Since(start)
		e.Processed++
		logCh <- fmt.Sprintf("員工 %d 完成處理 %s，耗時 %v", e.ID, item.Name(), elapsed)
	}
}

type Item1 struct{}

func (i Item1) Name() string { return "Item1" }
func (i Item1) Process()     { time.Sleep(200 * time.Millisecond) }

type Item2 struct{}

func (i Item2) Name() string { return "Item2" }
func (i Item2) Process()     { time.Sleep(400 * time.Millisecond) }

type Item3 struct{}

func (i Item3) Name() string { return "Item3" }
func (i Item3) Process()     { time.Sleep(600 * time.Millisecond) }

type Item interface {
	// Process 這是一個耗時操作
	Process()
	Name() string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 產生三種物品，各 10 件
	var items []Item
	for i := 0; i < 10; i++ {
		items = append(items, Item1{})
		items = append(items, Item2{})
		items = append(items, Item3{})
	}

	// 將物品順序隨機打亂
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	// 工作隊列 & log channel
	jobs := make(chan Item)
	logCh := make(chan string, 100)

	// 啟動一個 logger goroutine，負責印出所有紀錄
	var logWg sync.WaitGroup
	logWg.Add(1)
	go func() {
		defer logWg.Done()
		for msg := range logCh {
			fmt.Println(time.Now().Format("15:04:05.000"), msg)
		}
	}()

	// 建立 5 個員工
	employees := make([]*Employee, 0, 5)
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		e := &Employee{ID: i}
		employees = append(employees, e)
		wg.Add(1)
		go e.Work(&wg, jobs, logCh)
	}

	// 開始計時
	start := time.Now()

	// 把所有物品送進 jobs channel
	go func() {
		for _, it := range items {
			jobs <- it
		}
		close(jobs)
	}()

	// 等待所有員工處理完
	wg.Wait()

	// 關閉 log，並等待 logger 結束
	close(logCh)
	logWg.Wait()

	totalElapsed := time.Since(start)

	// 統計結果
	fmt.Println("====================================")
	fmt.Printf("全部物品處理完成，總處理時間：%v\n", totalElapsed)
	for _, e := range employees {
		fmt.Printf("員工 %d 共處理 %d 件物品\n", e.ID, e.Processed)
	}
}
