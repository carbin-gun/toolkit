package base36_test

import (
	"github.com/carbin-gun/toolkit/rand/base36"
	"regexp"
	"fmt"
)

func ExampleRand36() {
	result := base36.Rand()
	if len(result) != 16+3 {
		fmt.Println("result size should be 19")
	}
	pattern := regexp.MustCompile("[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}")
	if !pattern.MatchString(result) {
		fmt.Println("pattern not match")
	} else {
		fmt.Println("Success")
	}
	// Output:
	// Success
}
