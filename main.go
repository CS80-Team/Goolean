package main

import (
	"Boolean-IR-System/internal/engine"
	"Boolean-IR-System/shell"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"

	"github.com/joho/godotenv"
)

func openFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return nil
	}

	return cmd.Start()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	engine := engine.NewEngine()
	engine.LoadDirectory(os.Getenv("DATASET_PATH"))

	s := shell.NewShell()

	s.SetInputStream(os.Stdin)
	s.SetOutputStream(os.Stdout)

	s.RegisterCommand(shell.Command{
		Name:        "open",
		Description: "Open a document by ID in the default editor",
		Handler: func(args []string) shell.Status {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				s.WriteOutput("Invalid document ID\n")
			}
			if id < 0 || id >= len(engine.GetDocuments()) {
				s.WriteOutput("Document ID out of range\n")
			}

			doc := engine.GetDocumentByID(id)
			err = openFile(doc.Path + doc.Name)
			if err != nil {
				s.WriteOutput("Error opening file\n")
			}

			return shell.OK
		},
		Usage: "open <document_id>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "exit",
		Description: "Exit the shell",
		Handler: func(args []string) shell.Status {
			return shell.EXIT
		},
	})

	s.RegisterCommand(shell.Command{
		Name:        "query",
		Description: "Query the engine for a keyword or a boolean expression. Use AND, OR and NOT to build complex queries.",
		Handler: func(args []string) shell.Status {
			res := engine.Query(args)

			if res == nil || res.GetLength() == 0 {
				s.WriteOutput("No results found\n")
				return shell.OK
			}

			for i := 0; i < res.GetLength(); i++ {
				s.WriteOutput(
					strconv.Itoa(engine.GetDocumentByID(res.At(i)).ID) +
						". " +
						engine.GetDocumentByID(res.At(i)).Name + "\n")
			}

			return shell.OK
		},
		Usage: "query <keyword> | <expression>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "list",
		Description: "List all documents, displayable by name or/and path or/and ID or/and all",
		Handler: func(args []string) shell.Status {
			if len(args) == 0 {
				s.WriteOutput("No arguments provided\n")
				return shell.FAIL
			}

			var validArgs = []string{"-name", "-path", "-id", "-all"}
			var seen = make(map[string]bool)

			for _, arg := range args {
				if !slices.Contains(validArgs, arg) {
					s.WriteOutput("Invalid argument: " + arg + "\n")
					return shell.FAIL
				}

				if seen[arg] {
					s.WriteOutput("Duplicate argument: " + arg + "\n")
					return shell.FAIL
				}

				seen[arg] = true
			}

			if args[0] == "-all" {
				for i := 0; i < len(engine.GetDocuments()); i++ {
					doc := engine.GetDocumentByID(i)
					s.WriteOutput(doc.Name + " " + doc.Path + " " + strconv.Itoa(doc.ID) + "\n")
				}
				return shell.OK
			}

			for i := 0; i < len(engine.GetDocuments()); i++ {
				doc := engine.GetDocumentByID(i)
				for j := 0; j < len(args); j++ {
					switch args[j] {
					case "-name":
						s.WriteOutput(doc.Name)
					case "-path":
						s.WriteOutput(doc.Path)
					case "-id":
						s.WriteOutput(strconv.Itoa(doc.ID))
					}
					s.WriteOutput(" ")
				}
				s.WriteOutput("\n")
			}

			return shell.OK
		},
		Usage: "list <-name | -path | -id | -all>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "load",
		Description: "Load a new document into the engine",
		Handler: func(args []string) shell.Status {
			if len(args) != 1 {
				s.WriteOutput("Invalid number of arguments\n")
				return shell.FAIL
			}

			engine.LoadDirectory(args[0])
			return shell.OK
		},
		Usage: "load <document_path>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "open",
		Description: "Open a document by ID in the default editor",
		Handler: func(args []string) shell.Status {
			if len(args) != 1 {
				s.WriteOutput("Invalid number of arguments\n")
				return shell.FAIL
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				s.WriteOutput("Invalid document ID\n")
				return shell.FAIL
			}

			if id < 0 || id >= len(engine.GetDocuments()) {
				s.WriteOutput("Document ID out of range\n")
				return shell.FAIL
			}

			err = openFile(engine.GetDocumentByID(id).Path + engine.GetDocumentByID(id).Name)
			if err != nil {
				s.WriteOutput("Error opening file\n")
				return shell.FAIL
			}

			return shell.OK
		},
		Usage: "open <document_id>",
	})

	s.Run()
}
