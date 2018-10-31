package main

import (
	"log"
	"reflect"
	"strconv"
	"time"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: this will need adapting to whatever format we request for Sobriety_Date and Birth_Date
func timeConverter(value string) reflect.Value {
	tstamp, _ := strconv.ParseInt(value, 10, 64)
	return reflect.ValueOf(time.Unix(tstamp, 0))
}
