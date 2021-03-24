package main

import "fmt"

/*

// 以太坊中判断是否为 32bit系统的代码

if 32<<(^uintptr(0)>>63) == 32 && mem.Total > 2*1024*1024*1024 {
			log.Warn("Lowering memory allowance on 32bit arch", "available", mem.Total/1024/1024, "addressable", 2*1024)
			mem.Total = 2 * 1024 * 1024 * 1024
		}
*/

func main() {

    fmt.Printf("%X\n", ^uintptr(0))
    fmt.Printf("%v\n", (^uintptr(0)>>63))
    fmt.Printf("%X\n", 32<<(^uintptr(0)>>63))
    if (32<<(^uintptr(0)>>63)) == 32 {
        fmt.Println("yes")
    }else{
        fmt.Println("no")
    }
}
