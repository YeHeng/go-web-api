package main

import (
	"github.com/YeHeng/go-web-api/internal/pkg/core"
	"github.com/YeHeng/go-web-api/internal/pkg/factory"
)

func main() {

	for _, p := range factory.GetAllBeans() {
		p.Init()
	}

	_, e := core.Create()
	if e != nil {
		panic(e)
	}

}
