package utils

import (
	"errors"
	log2 "github.com/ethereum/go-ethereum/log"
	"log"
	"os"
	"reflect"
)

var (
	ErrTooManyRecords     = errors.New("collector: too many records")
	ErrInvalidKey         = errors.New("etherscan: invalid api key")
	ErrEtherscanRateLimit = errors.New("etherscan: rate limit")
)

func IsEmpty(target any) bool {
	t := reflect.TypeOf(target)
	if t.Kind() == reflect.Ptr { //指针类型获取真正type需要调用Elem
		return (reflect.ValueOf(target).IsNil()) || (reflect.ValueOf(target).Elem().IsValid())
	}

	/*if reflect.ValueOf(target).IsValid() {
		t = reflect.ValueOf(target).Type()
	}
	newStruc := reflect.New(t)
	fmt.Println(newStruc)

	if (reflect.DeepEqual(target, newStruc)) ||
		target == (newStruc) {
		return true
	}*/
	println("param must be ptr")
	return false
}

func LogPrint(info, file_path string) {
	file, err := os.OpenFile(file_path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.Println(info)
}

func LogError(info, file_path string) {
	file, err := os.OpenFile(file_path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log2.Error(info)
	log2.Error(info)
}
