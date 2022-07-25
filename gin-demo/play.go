package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	// number := math.MaxInt64
	var wg sync.WaitGroup
	wg.Add(100)
	var rslt []string
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			fmt.Printf("go func: %d\n", i)
			rslt = append(rslt /* "#"+string(i) */, fmt.Sprint("#", i))
			// time.Sleep(time.Second)
		}(i)
	}
	wg.Wait()
	fmt.Println(rslt)
	/* mounths, err := genTables(1, 8)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mounths) */
}

func genTables(start uint, end uint) ([]string, error) {
	var ret []string

	if start <= 0 || end > 12 {
		return ret, errors.New("between 1 and 12")
	}

	var count uint = 0
	for i := start; i < end; i++ {
		ret = append(ret, fmt.Sprintf("cdr_postbrother_%d", i))
		count++
	}

	return ret, nil
}

func genMounths(start uint, end uint) ([12]uint, error) {
	var ret [12]uint
	if start <= 0 || end > 12 {
		return ret, errors.New("between 1 and 12")
	}

	var count uint = 0
	for i := uint(6); i < 9; i++ {
		ret[count] = i
		count++
	}
	return ret, nil
}
