package main

import (
	"errors"
	"fmt"
)

func main() {
	// var ret []string
	// var count uint = 0
	// for i := uint(6); i < 9; i++ {
	// 	ret = append(ret, fmt.Sprintf("cdr_postbrother_%d", i))
	// 	count++
	// }
	// fmt.Println(ret)
	mounths, err := genTables(1, 8)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mounths)
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
