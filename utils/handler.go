package utils

import "log"

func Handle(funcs ...func() error) {
	for _, f := range funcs {
		err := f()
		if err != nil {
			log.Panic(err)
		}
	}
}
