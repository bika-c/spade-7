package main

import (
	"fmt"
	"reflect"
	"runtime"
	"spade-7/Game"
	"spade-7/Spade7"
	"strings"
)

func abc() {

}

func main() {
	var g Game.Game = &Spade7.Spade7{}

	str := runtime.FuncForPC(reflect.ValueOf(g.RemovePlayers).Pointer()).Name()
	str = str[strings.LastIndex(str, ".")+1 : len(str)-3]
	fmt.Println(str)
	str2 := runtime.FuncForPC(reflect.ValueOf(abc).Pointer()).Name()
	str2 = str2[strings.LastIndex(str2, ".")+1 : len(str2)]
	fmt.Println(str2)
}
