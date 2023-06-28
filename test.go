package main

import (
	"fmt"
	"starlink/utils"
	"time"
)

func main() {
	var expireSec int64 = 3
	m := utils.NewExpiredMap()
	for i := 0; i < 10; i++ {
		m.Set(fmt.Sprintf("satellite-%d", i), true, expireSec)
		time.Sleep((1) * time.Second)
		print(m)
	}
	time.Sleep(time.Duration(expireSec) * time.Second)
	print(m)
}

func print(m *utils.ExpiredMap) {
	for k, v := range m.GetAll() {
		fmt.Printf("key: %v, value: %d\n", k, v)
	}
	fmt.Println()
}
