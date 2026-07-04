package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/CallMeJaja/zte-cli/router"
)

// Version is set by main before launching interactive mode.
var Version = "dev"

var banner string

func init() {
	banner = fmt.Sprintf(`
==========================================================
                 ZTE F609 Router Manager v%s
==========================================================
(1) Status                         (4) Application
(2) Network                        (5) Administration
(3) Security                       (0) Exit
==========================================================
`, Version)
}

// mainMenuItems defines the top-level modules.
var mainMenuItems = []MenuItem{
	{ID: 1, Name: "Status"},
	{ID: 2, Name: "Network"},
	{ID: 3, Name: "Security"},
	{ID: 4, Name: "Application"},
	{ID: 5, Name: "Administration"},
}

// RunInteractive launches the interactive TUI menu loop.
func RunInteractive(client *router.Client) {
	scanner := bufio.NewScanner(os.Stdin)

	// Login first
	fmt.Println("\n  Connecting to router...")
	ok, err := client.Login()
	if err != nil {
		fmt.Printf("\n  [ERROR] %s\n", err)
		return
	}
	if !ok {
		fmt.Println("\n  [ERROR] Login failed.")
		return
	}
	fmt.Println("  [OK] Logged in successfully.\n")

	for {
		fmt.Print(banner)
		fmt.Print("Pls enter command number: ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())
		if choice == "" {
			continue
		}

		switch choice {
		case "1":
			RunStatusModule(client, scanner)
		case "2":
			RunNetworkModule(client, scanner)
		case "3":
			runSecurityModule(client, scanner)
		case "4":
			runApplicationModule(client, scanner)
		case "5":
			runAdministrationModule(client, scanner)
		case "0":
			fmt.Println("\n  Goodbye!")
			return
		default:
			fmt.Println("\n  [!] Invalid option. Please try again.")
		}
	}
}

// runSecurityModule placeholder for Phase 4.
func runSecurityModule(client *router.Client, scanner *bufio.Scanner) {
	fmt.Println("\n  [!] Security module is not yet implemented.")
	fmt.Println("      Coming in v0.4.0!")
	fmt.Println()
}

// runApplicationModule placeholder for Phase 5.
func runApplicationModule(client *router.Client, scanner *bufio.Scanner) {
	fmt.Println("\n  [!] Application module is not yet implemented.")
	fmt.Println("      Coming in v0.5.0!")
	fmt.Println()
}

// runAdministrationModule placeholder for Phase 6.
func runAdministrationModule(client *router.Client, scanner *bufio.Scanner) {
	fmt.Println("\n  [!] Administration module is not yet implemented.")
	fmt.Println("      Coming in v0.6.0!")
	fmt.Println()
}
