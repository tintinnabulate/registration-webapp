package main

import (
	"html/template"
)

var funcMap = template.FuncMap{
	"inc": inc,
}

func inc(i int) int {
	return i + 1
}
