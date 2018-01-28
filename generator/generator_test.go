package generator

import (
	"fmt"
	"os"
)

func ExampleGenerator() {
	// log.SetOutput(ioutil.Discard)
	if err := NewGenerator(os.Getenv("PWD")+"/test", "ignore").Generate("*"); err != nil {
		fmt.Println("generation aborted:", err)
	}
	fmt.Println("ok")
	// Output:
	// ok
}
