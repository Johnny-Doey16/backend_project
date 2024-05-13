package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func runMicroservices(command string) {
	cmd := exec.Command("go", "run", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running microservices %s: %s\n", command, err)
	}
}

func main() {
	var wg sync.WaitGroup

	microservices := []string{
		"internal/app/authentication/main.go",
		"internal/app/church/main.go",
		"internal/app/posts/main.go",
		"internal/app/bible/main.go",
		"internal/app/prayer/main.go",
	}

	for _, ms := range microservices {
		wg.Add(1)
		go func(microservice string) {
			defer wg.Done()
			runMicroservices(microservice)
		}(ms)
	}

	wg.Wait()
}
