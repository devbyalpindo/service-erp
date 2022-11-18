package helper

import "log"

func PanicIfError(err error) {
	if err != nil {
		println(err)
		log.Panic(err)
	}
}
