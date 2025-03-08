package commands

import (
	"Boolean-IR-System/internal/engine"
	"Boolean-IR-System/shell"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
)

func RegisterCommands(s *shell.Shell, engine *engine.Engine) {
	registerBuiltInCommands(s)

	s.RegisterCommand(shell.Command{
		Name:        "open",
		Description: "Open a document by ID in the default editor",
		Handler:     openCommand(engine),
		Usage:       "open <document_id>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "query",
		Description: "Query the engine for a keyword or a boolean expression",
		Handler:     queryCommand(engine),
		Usage:       "query <keyword> | <expression>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "list",
		Description: "List all documents, displayable by name or/and path or/and ID. Use -n to limit the number of documents.",
		Handler:     listCommand(engine),
		Usage:       "list <-id | -name | -path>",
	})

	s.RegisterCommand(shell.Command{
		Name:        "load",
		Description: "Load a new document into the engine",
		Handler:     loadCommand(engine),
		Usage:       "load <document_path>",
	})

	s.RegisterEarlyExecCommand(shell.EarlyCommand{
		Name:        "engine-stats",
		Description: "Displays the total number of documents and keys in the engine",
		Handler:     engineStatsCommand(engine),
		Usage:       "engine-stats",
	})
}

func registerBuiltInCommands(sh *shell.Shell) {
	sh.RegisterCommand(shell.Command{
		Name:        "exit",
		Description: "Exit the shell",
		Handler: func(s *shell.Shell, args []string) shell.Status {
			return shell.EXIT
		},
		Usage: "exit",
	})

	sh.RegisterCommand(shell.Command{
		Name:        "help",
		Description: "List all available commands",
		Handler: func(s *shell.Shell, args []string) shell.Status {
			for _, cmd := range sh.GetCommands() {
				sh.Write(cmd.Name + ": " + cmd.Description + "\n")
				sh.Write("    Usage: " + cmd.Usage + "\n")
			}
			return shell.OK
		},
		Usage: "help",
	})

	sh.RegisterCommand(shell.Command{
		Name:        "clear",
		Description: "Clear the screen",
		Handler: func(s *shell.Shell, args []string) shell.Status {
			switch runtime.GOOS {
			case "windows":
				cmd := exec.Command("cmd", "/c", "cls")
				cmd.Stdout = os.Stdout
				_ = cmd.Run()
			default:
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				_ = cmd.Run()
			}

			return shell.OK
		},
		Usage: "clear",
	})
}

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

func openCommand(engine *engine.Engine) func(s *shell.Shell, args []string) shell.Status {
	return func(s *shell.Shell, args []string) shell.Status {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			s.Write("Invalid document ID\n")
			return shell.FAIL
		}
		if id < 0 || id >= len(engine.GetDocuments()) {
			s.Write("Document ID out of range, Docs Ids [0, " + strconv.Itoa(len(engine.GetDocuments())-1) + "]\n")
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
	}
}

func queryCommand(engine *engine.Engine) func(s *shell.Shell, args []string) shell.Status {
	return func(s *shell.Shell, args []string) shell.Status {
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

		s.Write("Total documents found: " + strconv.Itoa(res.GetLength()) + "\n")

		return shell.OK
	}
}

func listCommand(engine *engine.Engine) func(s *shell.Shell, args []string) shell.Status {
	return func(s *shell.Shell, args []string) shell.Status {
		var stat = shell.OK

		if len(args) == 0 {
			for i := 0; i < len(engine.GetDocuments()); i++ {
				doc := engine.GetDocumentByID(i)
				s.Write(strconv.Itoa(doc.ID) + " " + doc.Name + " " + doc.Path + "\n")
			}
			stat = shell.OK
		} else {

			var validArgs = []string{"-name", "-path", "-id", "-n"}
			var displayFields []string
			var seen = make(map[string]bool)
			var limit = -1

			for i := 0; i < len(args); i++ {
				arg := args[i]
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

				if arg == "-n" {
					if i+1 >= len(args) {
						s.Write("Missing value for -n\n")
						return shell.FAIL
					}

					n, err := strconv.Atoi(args[i+1])
					if err != nil || n < 1 {
						s.Write("Invalid value for -n. Must be a positive integer.\n")
						return shell.FAIL
					}

					limit = n
					i++
				} else {
					displayFields = append(displayFields, arg)
				}
			}

			if stat == shell.FAIL {
				goto end
			}

			if len(displayFields) == 0 {
				displayFields = []string{"-id", "-name", "-path"}
			}

			totalDocs := len(engine.GetDocuments())
			if limit == -1 || limit > totalDocs {
				limit = totalDocs
			}

			for i := 0; i < limit; i++ {
				doc := engine.GetDocumentByID(i)
				for _, field := range displayFields {
					switch field {
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
	}
}

func loadCommand(engine *engine.Engine) func(s *shell.Shell, args []string) shell.Status {
	return func(s *shell.Shell, args []string) shell.Status {
		if len(args) != 1 {
			s.Write("Invalid number of arguments\n")
			return shell.FAIL
		}

		engine.LoadDirectory(args[0])
		return shell.OK
	}
}

func engineStatsCommand(engine *engine.Engine) func(s *shell.Shell) {
	return func(s *shell.Shell) {
		s.Write("Engine stats:\n")
		s.Write("Total documents: " + strconv.Itoa(len(engine.GetDocuments())) + "\n")
		s.Write("Total keys: " + strconv.Itoa(engine.GetIndexSize()) + "\n")
	}
}
