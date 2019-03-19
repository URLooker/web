package test

import (
	"testing"
	"time"
	"fmt"
)

func Test_Zero(t * testing.T) {
	date_time := time.Unix(1552639137, 0)
	fmt.Printf("timestamp to datetime:%v\n", date_time.Format("2006-01-02 15:04:05"))
}
