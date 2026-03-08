package utils

import (
	"cbs-simulator/version"
	"fmt"
)

func StartupScreen(addr string) {

	fmt.Println("")
	fmt.Printf("CBS-SIMULATOR v%s\n", version.Version)
	fmt.Printf("Build: %s\n", version.Build)
	fmt.Println("")

	fmt.Println("- Config loaded")
	fmt.Println("- Database connected")
	fmt.Println("- Migrations completed")
	fmt.Println("- Router initialized")
	fmt.Println("")

	fmt.Printf("Server running on : %s\n", addr)
	fmt.Println("API Base URL      : http://localhost:8080/api/v1")
	fmt.Println("")
}
