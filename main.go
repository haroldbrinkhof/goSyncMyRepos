package main

import (
	"fmt"
	"os"
	"log"
	"github.com/go-git/go-git/v5"
)

func main(){
	workingDir, err := fmt.Println(os.Getwd())
	if err != nil {
		log.Fatal("could not determine working directory.")
	}


	fmt.Println(workingDir)

}
