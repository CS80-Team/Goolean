package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/CS80-Team/BooleanEngine/internal"
	"github.com/CS80-Team/BooleanEngine/internal/engine"
	"github.com/CS80-Team/BooleanEngine/shell"
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
		Name: "list",
		Description: "List all documents, displayable by name or/and path or/and ID or/and extension\n" +
			"Use -sortby to sort results by name, path, id or extension\n" +
			"Use -n to limit the number of results\n" +
			"Default fields order: -id -name -path -ext\n" +
			"Default sortby: id\n" +
			"Default limit: all",
		Handler: listCommand(engine),
		Usage:   "list <-id | -name | -path | -ext> [-n <limit>] [-sortby <name | path | id | -ext>]",
	})

	s.RegisterCommand(shell.Command{
		Name:        "load",
		Description: "Load a new document into the engine",
		Handler:     loadCommand(engine),
		Usage:       "load <document_path>",
	})

	s.RegisterCommand(shell.Command{
		Name: "find",
		Description: "Find a document by name or id.\n" +
			"Display the document's id, name and path\n" +
			"Default search field: -name\n",
		Handler: findCommand(engine),
		Usage:   "find <-id | -name> <value> || find <document_name>",
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
			// Convert the path to windows path
			convertCmd := exec.Command("wslpath", "-w", path)
			path, err := convertCmd.Output()
			if err != nil {
				return err
			}
			cmd = exec.Command("explorer.exe", strings.TrimSpace(string(path)))
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

func getDocumentByField(engine *engine.Engine, search, value string) (internal.Document, string) {
	var doc internal.Document
	if search != "-id" && search != "-name" {
		return doc, "Invalid search field, must be one of: -id, -name\n"
	}

	switch search {
	case "-id":
		id, err := strconv.Atoi(value)
		if err != nil {
			return doc, "Invalid document ID\n"
		}

		if id < 0 || id >= engine.GetDocumentsSize() {
			return doc, "Document ID out of range, Docs Ids [0, " + strconv.Itoa(engine.GetDocumentsSize()-1) + "]\n"
		}

		doc = engine.GetDocumentByIDCopy(id)
	default: // -name
		doc = engine.GetDocumentByNameCopy(value)
	}

	return doc, ""
}

func findCommand(engine *engine.Engine) func(s *shell.Shell, args []string) shell.Status {
	return func(s *shell.Shell, args []string) shell.Status {
		var doc internal.Document
		if len(args) > 2 {
			s.Write("Invalid number of arguments\n")
			return shell.FAIL
		}

		if len(args) == 1 {
			if args[0] == "-id" || args[0] == "-name" {
				s.Write("Missing value for search field\n")
				return shell.FAIL
			}
			doc = engine.GetDocumentByNameCopy(args[0])

		} else {
			var err string
			doc, err = getDocumentByField(engine, args[0], args[1])
			if err != "" {
				s.Write(err)
				return shell.FAIL
			}
		}

		if doc == (internal.Document{}) {
			s.Write("Document not found\n")
		} else {
			s.Write("Document found:\n" + doc.String())
		}

		return shell.OK
	}
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
		doc := engine.GetDocumentByID(id)
		path := filepath.Join(doc.DirectoryPath, doc.Name+doc.Ext)
		s.Write("Opening document: " + doc.Name + ", path: " + path + "\n")
		err = openFile(path)
		if err != nil {
			s.Write("Error opening file\n")
			s.Write(err.Error())
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
		s.Write("Total documents: " + strconv.Itoa(engine.GetDocumentsSize()) + "\n")
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
	var validArgs = []string{"-name", "-path", "-id", "-ext", "-n", "-sortby"}
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

			if sortby != "name" && sortby != "path" && sortby != "id" && sortby != "ext" {
				return nil, -1, "", "Invalid value (" + sortby + ") for -sortby. Must be one of: name, path, id, ext."
			}

			i++
		} else {
			displayFields = append(displayFields, arg)
		}
	}

	if len(displayFields) == 0 {
		displayFields = []string{"-id", "-name", "-path", "-ext"}
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
			return strings.Compare(i.DirectoryPath, j.DirectoryPath)
		})
	} else if sortby == "-ext" {
		slices.SortFunc(docs, func(i, j *internal.Document) int {
			return strings.Compare(i.Ext, j.Ext)
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
				s.Write(docs[i].DirectoryPath + " ")
			case "-ext":
				s.Write(docs[i].Ext + " ")
			}
		}
		s.Write("\n")
	}

	s.Write("Total documents: " + strconv.Itoa(limit) + "\n")
	s.Write("Sorted by: " + sortby + "\n")
}
