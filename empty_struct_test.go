package main

import (
	"fmt"
	"testing"
)

type a struct {

}

type b struct {

}

func (v *a) view(){
	fmt.Println(1)
}

func (v *b) view(){
	fmt.Println(2)
}

func TestEmptyS(t *testing.T){
	a := a{}
	b := b{}
	a.view()
	b.view()
}