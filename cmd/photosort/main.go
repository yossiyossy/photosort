package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/yossiyossy/photosort/internal"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("executable: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	fmt.Println("Target directory:", exeDir)

	if err := internal.OrganizeInPlace(exeDir); err != nil {
		log.Fatalf("organize: %v", err)
	}

	fmt.Println("Done.")
}
