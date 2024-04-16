package main

import (
	"fmt"
	"time"
)

func test_date() {
	//2023-12-21 16:00:33 +0700 WIB
	a := 1703149233975 / 1000

	fmt.Println(time.Unix(int64(a), 0))
}
