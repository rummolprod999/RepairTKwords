package main

import (
	"flag"
	"fmt"
	"strconv"
)

func init() {
	flag.Parse()
	ArgS := flag.Arg(0)
	i, err := strconv.Atoi(ArgS)
	Typefz = i
	if err != nil {
		fmt.Println("argument has not been converted to int")
		panic("")
	}
	CreateEnv()
}

func main() {
	defer SaveStack()
	Logging("Start work")
	p := Eis{}
	p.run()
	Logging(fmt.Sprintf("Count purchases %d", p.Count))
	Logging(fmt.Sprintf("Update purchases %d", p.Updated))
	Logging("End work")
}
