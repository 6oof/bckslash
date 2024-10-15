package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/6oof/bckslash/pkg/constants"
	"github.com/6oof/bckslash/pkg/helpers"
	"github.com/6oof/bckslash/pkg/views"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var cmdDashboard = &cobra.Command{
		Use:   "dashboard",
		Short: "Start the Bckslash dashboard",
		Long:  "Launches the Bckslash dashboard for managing your applications.",
		Run: func(cmd *cobra.Command, args []string) {
			runTui() // This is where your TUI function will be called
		},
	}

	var cmdDeploy = &cobra.Command{
		Use:   "deploy [message]",
		Short: "Deploy an application",
		Long: `Deploys the specified application.
This command is not yet implemented.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Deploy command received with message:", strings.Join(args, " "))
		},
	}

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "Bckslash CLI",
		Long:  "Welcome to the Bckslash CLI!\nThis application allows you to manage your applications.",
	}

	rootCmd.AddCommand(cmdDashboard, cmdDeploy)
	rootCmd.Execute()
}

func runTui() {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	err := helpers.OpenDb(constants.DatabaseFile)
	if err != nil {
		panic("Falid opening bcks database: " + err.Error())
	}
	defer helpers.CloseDb()

	p := tea.NewProgram(views.InitHomeModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	fmt.Println("__/\\\\\\_______________        ")
	fmt.Println(" _\\///\\\\\\_____________       ")
	fmt.Println("  ___\\///\\\\\\___________      ")
	fmt.Println("   _____\\///\\\\\\_________     ")
	fmt.Println("    _______\\///\\\\\\_______    ")
	fmt.Println("     _________\\///\\\\\\_____   ")
	fmt.Println("      ___________\\///\\\\\\___  ")
	fmt.Println("       _____________\\///\\\\\\_ ")
	fmt.Println("        _______________\\///__  ")
	fmt.Println("")
	fmt.Println("!!! Hey there, Caution Ahead !!!")
	fmt.Println("You're still SSH-ed into the server.")
	fmt.Println("If you're still busy doing your thing, keep going!")
	fmt.Println("But if you're done here, press Ctrl+D or type 'exit' to log off.")

}
