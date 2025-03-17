package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/CS80-Team/Boolean-IR-System/internal"
	"github.com/CS80-Team/Boolean-IR-System/internal/engine"
	"github.com/CS80-Team/Boolean-IR-System/shell"
)

func RegisterCommands(s *shell.Shell, engine *engine.Engine) {
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

func openFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		if isWSL() {
			cmd = exec.Command("wslview", path)
		} else {
			cmd = exec.Command("xdg-open", path)
		}
	default:
		return nil
	}

	return cmd.Start()
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
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

func listCommand(engine *engine.Engine) func(s *shell.Shell, args []string) shell.Status {
	return func(s *shell.Shell, args []string) shell.Status {
		displayFields, limit, sortby, err := parseListArgs(args)
		if err != "" {
			s.Write(err + "\n")
			return shell.FAIL
		}

		docs := getSortedDocuments(engine, sortby)
		if limit == -1 || limit > len(docs) {
			limit = len(docs)
		}

		displayDocuments(s, docs, displayFields, limit, sortby)
		return shell.OK
	}
}

func parseListArgs(args []string) ([]string, int, string, string) {
	var validArgs = []string{"-name", "-path", "-id", "-n", "-sortby"}
	var displayFields []string
	var seen = make(map[string]bool)
	var limit = -1
	var sortby = "id"

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !slices.Contains(validArgs, arg) {
			return nil, -1, "", "Invalid argument: " + arg
		}

		if seen[arg] {
			return nil, -1, "", "Duplicate argument: " + arg
		}

		seen[arg] = true

		if arg == "-n" {
			if i+1 >= len(args) {
				return nil, -1, "", "Missing value for -n"
			}

			n, err := strconv.Atoi((args)[i+1])
			if err != nil || n < 1 {
				return nil, -1, "", "Invalid value for -n. Must be a positive integer."
			}

			limit = n
			i++
		} else if arg == "-sortby" {
			if i+1 >= len(args) {
				return nil, -1, "", "Missing value for -sortby"
			}

			sortby = args[i+1]

			if sortby != "name" && sortby != "path" && sortby != "id" {
				return nil, -1, "", "Invalid value (" + sortby + ") for -sortby. Must be one of: name, path, id."
			}

			i++
		} else {
			displayFields = append(displayFields, arg)
		}
	}

	if len(displayFields) == 0 {
		displayFields = []string{"-id", "-name", "-path"}
	}

	return displayFields, limit, sortby, ""
}

func getSortedDocuments(engine *engine.Engine, sortby string) []*internal.Document {
	if sortby == "id" {
		return engine.GetDocuments()
	}

	docs := engine.GetDocumentsCopy()
	if sortby == "name" {
		slices.SortFunc(docs, func(i, j *internal.Document) int {
			return strings.Compare(i.Name, j.Name)
		})
	} else if sortby == "path" {
		slices.SortFunc(docs, func(i, j *internal.Document) int {
			return strings.Compare(i.Path, j.Path)
		})
	}
	return docs
}

func displayDocuments(s *shell.Shell, docs []*internal.Document, displayFields []string, limit int, sortby string) {
	for i := 0; i < limit; i++ {
		for _, field := range displayFields {
			switch field {
			case "-id":
				s.Write(strconv.Itoa(docs[i].ID) + " ")
			case "-name":
				s.Write(docs[i].Name + " ")
			case "-path":
				s.Write(docs[i].Path + " ")
			}
		}
		s.Write("\n")
	}

	s.Write("Total documents: " + strconv.Itoa(limit) + "\n")
	s.Write("Sorted by: " + sortby + "\n")
}
