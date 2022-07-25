package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	/* c := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			c <- i
			fmt.Printf("send %d\n", i)
			time.Sleep(time.Second)
		}
		fmt.Println("ready close channel")
		close(c)
	}()

	for i := range c {
		fmt.Printf("receive %d\n", i)
	}
	fmt.Println("quit for loop") */

	// 1. 设置长度  2. 多wg方式
	// c := make(chan []string)

	// var dbWg sync.WaitGroup

	// var rslt []string

	// for i := 0; i < 2; i++ {
	// 	dbWg.Add(1)
	// 	go func(i int, cInner chan<- []string) {
	// 		defer func() {
	// 			dbWg.Done()
	// 		}()

	// 		cInner <- []string{"aa" + fmt.Sprint(i)}
	// 	}(i, c)
	// }
	// go func() {
	// 	dbWg.Wait()
	// 	close(c)
	// 	for received := range c {
	// 		fmt.Println(received)
	// 	}
	// }()
	// // dbWg.Wait()

	// // time.Sleep(time.Second)

	// fmt.Println("quit for loop")
	// // for received := range c {
	// // 	fmt.Println(received)
	// // }
	// fmt.Println(rslt)

	dbList := []string{"db1", "db2"}
	tableList := []string{"table1", "table2", "table3"}
	cSize := len(dbList) * len(tableList)
	fmt.Println(cSize)
	c := make(chan []string, cSize)
	var wg sync.WaitGroup

	for _, dbName := range dbList {
		wg.Add(1)
		fmt.Println(dbName)
		go func(dbName string) {
			defer func() {
				wg.Done()
			}()

			queryDB(dbName, tableList, c, &wg)
		}(dbName)
	}
	wg.Wait()
	close(c)

	var rslt []string
	for received := range c {
		fmt.Println(received)
		rslt = append(rslt, received...)
	}
	fmt.Println("quit for loop.\n rslt =>", rslt)
}

func queryDB(dbName string, tables []string, c chan<- []string, wg *sync.WaitGroup) {
	for _, t := range tables {
		wg.Add(1)
		go func(t string) {
			defer func() {
				wg.Done()
			}()
			queryTable(dbName, t, c)
		}(t)
	}
}

func queryTable(dbName string, tableName string, c chan<- []string) {

	r := rand.Intn(10)
	time.Sleep(time.Duration(r) * time.Second)

	c <- []string{fmt.Sprintf("#%s_#%s_val1", dbName, tableName), fmt.Sprintf("#%s_#%s_val2", dbName, tableName)}
}

func mutipleDBs(c chan<- []string) {
	dbCount := 3
	// tableCount := 5
	// var dbWg sync.WaitGroup
	// dbWg.Add(dbCount)
	for i := dbCount; i < dbCount; i++ {
		go func(dbIdx int, c chan<- []string) {
			defer func() {
				// dbWg.Done()
			}()
			fmt.Println(dbIdx)

			c <- []string{"aa"}

			// multipleRoutines(dbIdx, c)
			// for i := 0; i < tableCount; i++ {
			// 	go multipleRoutines(dbIdx, c)
			// }
		}(i, c)
	}
	// dbWg.Wait()
}

func multipleRoutines(dbIdx int, c chan []string) {
	// var wg sync.WaitGroup
	// wg.Add(10)
	for tableIdx := 0; tableIdx < 10; tableIdx++ {
		go func(dbIdx int, tableIdx int) {
			// defer func() {
			// 	wg.Done()
			// }()
			list := []string{fmt.Sprintf("dbIdx#%d_tableIdx#%d_0", dbIdx, tableIdx), fmt.Sprintf("dbIdx#%d_tableIdx#%d_1", dbIdx, tableIdx)}
			c <- list
		}(dbIdx, tableIdx)
	}
	// wg.Wait()
}
