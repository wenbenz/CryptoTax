package main

import (
	"fmt"
	"time"

	"github.com/wenbenz/CryptoTax/tools"
)

func main() {
	fmt.Println(tools.GetValueAtTime(time.Now().Add(-10 * time.Hour), "NEXO"))
}
