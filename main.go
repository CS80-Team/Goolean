package main

import (
	"Boolean-IR-System/internal/engine"
	"Boolean-IR-System/shell"
	"os"
	"os/exec"
	"path/filepath"
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
				s.Write("Invalid document ID\n")
				return shell.FAIL
			}
			if id < 0 || id >= len(engine.GetDocuments()) {
				s.Write("Document ID out of range, Docs Ids [0" + strconv.Itoa(len(engine.GetDocuments())-1) + "]\n")
				return shell.FAIL
			}
			s.Write("Opening document: " + engine.GetDocumentByID(id).Name)
			doc := engine.GetDocumentByID(id)
			err = openFile(filepath.Join(doc.Path, doc.Name))
			if err != nil {
				s.Write("Error opening file\n")
				return shell.FAIL
			}

			return shell.OK
		},
		Usage: "open <document_id>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "query",
		Description: "Query the engine for a keyword or a boolean expression. Use AND, OR and NOT to build complex queries.",
		Handler: func(args []string) shell.Status {
			res := engine.Query(args)

			if res == nil || res.GetLength() == 0 {
				s.Write("No results found\n")
				return shell.OK
			}

			for i := 0; i < res.GetLength(); i++ {
				s.Write(
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
		Description: "List all documents, displayable by name or/and path or/and ID",
		Handler: func(args []string) shell.Status {
			var stat shell.Status = shell.OK

			if len(args) == 0 {
				for i := 0; i < len(engine.GetDocuments()); i++ {
					doc := engine.GetDocumentByID(i)
					s.Write(doc.Name + " " + doc.Path + " " + strconv.Itoa(doc.ID) + "\n")
				}
				stat = shell.OK
			} else {

				var validArgs = []string{"-name", "-path", "-id"}
				var seen = make(map[string]bool)

				for _, arg := range args {
					if !slices.Contains(validArgs, arg) {
						s.Write("Invalid argument: " + arg + "\n")
						stat = shell.FAIL
						break
					}

					if seen[arg] {
						s.Write("Duplicate argument: " + arg + "\n")
						stat = shell.FAIL
						break
					}

					seen[arg] = true
				}

				if stat == shell.FAIL {
					goto end
				}

				for i := 0; i < len(engine.GetDocuments()); i++ {
					doc := engine.GetDocumentByID(i)
					for j := 0; j < len(args); j++ {
						switch args[j] {
						case "-name":
							s.Write(doc.Name)
						case "-path":
							s.Write(doc.Path)
						case "-id":
							s.Write(strconv.Itoa(doc.ID))
						}
						s.Write(" ")
					}
					s.Write("\n")
				}
			}

		end:

			if stat == shell.OK {
				s.Write("Total documents: " + strconv.Itoa(len(engine.GetDocuments())) + "\n")
			}

			return stat
		},
		Usage: "list <-name | -path | -id>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "load",
		Description: "Load a new document into the engine",
		Handler: func(args []string) shell.Status {
			if len(args) != 1 {
				s.Write("Invalid number of arguments\n")
				return shell.FAIL
			}

			engine.LoadDirectory(args[0])
			return shell.OK
		},
		Usage: "load <document_path>",
	})

	s.RegisterEarlyExecCommand(shell.EarlyCommand{
		Name:        "engine-stats",
		Description: "Displays the total number of documents and keys in the engine",
		Handler: func() {
			s.Write("Engine stats:\n")
			s.Write("Total documents: " + strconv.Itoa(len(engine.GetDocuments())) + "\n")
			s.Write("Total keys: " + strconv.Itoa(engine.GetIndexSize()) + "\n")
		},
		Usage: "engine-stats",
	})

	s.Run("Welcome to the Boolean IR System shell, type `help` for list of commands\n")
}
