package shell

import (
	"fmt"
	"io/ioutil"
	"log"
)

func Welcome() {
	file, err := ioutil.ReadFile("banner.txt")
	if err != nil {
		log.Printf("Failure when reading banner file: %v ", err)
	}
	banner := string(file)
	fmt.Printf("\n")
	fmt.Printf("%s\n", banner)
	fmt.Printf("Welcome to BabelDB CLI! - Version 0.1\n\n")
}
