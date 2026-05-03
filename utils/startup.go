package utils

import (
	"cbs-simulator/config"
	"cbs-simulator/version"
	"fmt"
	"time"
)

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	blue  = "\033[34m"
	cyan  = "\033[36m"
	green = "\033[32m"
	gray  = "\033[90m"
	white = "\033[97m"
)

func StartupScreen(addr string) {
	start := time.Now()
	env := config.AppConfig.Environment

	fmt.Println()
	fmt.Println(blue + "  ╔══════════════════════════════════════════╗" + reset)
	fmt.Printf(blue+"  ║  "+bold+white+"CBS-SIMULATOR"+reset+" "+green+"v%s"+reset+blue+"                    ║\n"+reset, version.Version)
	fmt.Printf(blue+"  ║  "+gray+"Build: %s"+blue+"                       ║\n"+reset, version.Build)
	fmt.Println(blue + "  ╚══════════════════════════════════════════╝" + reset)
	fmt.Println()

	items := []struct {
		label string
		value string
	}{
		{"Config    ", "loaded"},
		{"Database  ", "connected"},
		{"Migrations", "completed"},
		{"Router    ", "initialized"},
	}

	for _, item := range items {
		fmt.Printf("  "+green+"✔"+reset+"  "+gray+"%s"+reset+"  "+green+"%s"+reset+"\n", item.label, item.value)
	}

	fmt.Println()
	fmt.Println(gray + "  ───────────────────────────────────────────" + reset)
	fmt.Printf("  "+gray+"environment  "+reset+"%s\n", env)
	fmt.Printf("  "+gray+"server       "+reset+cyan+"%s"+reset+"\n", addr)
	fmt.Printf("  " + gray + "api base     " + reset + cyan + "http://localhost:8080/api/v1" + reset + "\n")
	fmt.Println(gray + "  ───────────────────────────────────────────" + reset)
	fmt.Println()
	fmt.Printf("  "+gray+"ready in %dms"+reset+"\n\n", time.Since(start).Milliseconds())
}
