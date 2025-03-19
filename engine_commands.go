package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/CS80-Team/Goolean/internal"
	"github.com/CS80-Team/Goolean/internal/engine"
	"github.com/CS80-Team/gshell/pkg/gshell"
)

func RegisterCommands(s *gshell.Shell, engine *engine.Engine) {
	s.RegisterCommand(
		gshell.NewCommand(
			"open",
			"Open a document by ID in the default editor",
			"open <document_id>",
			[]gshell.Argument{},
			[]string{},
			openCommand(engine),
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	s.RegisterCommand(
		gshell.NewCommand(
			"query",
			"Query the engine for a keyword or a boolean expression",
			"query <keyword> | <expression>",
			[]gshell.Argument{},
			[]string{},
			queryCommand(engine),
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	s.RegisterCommand(
		gshell.NewCommand(
			"list",
			"List all documents, displayable by name or/and path or/and ID or/and extension\n"+
				"Use -sortby to sort results by name, path, id or extension\n"+
				"Use -n to limit the number of results\n"+
				"Default fields order: -id -name -path -ext\n"+
				"Default sortby: id\n"+
				"Default limit: all",
			"list <-id | -name | -path | -ext> [-n <limit>] [-sortby <name | path | id | -ext>]",
			[]gshell.Argument{ // till now for testing autoCompleteArg()
				{Tag: "-id"},
				{Tag: "-name"},
				{Tag: "-path"},
				{Tag: "-ext"},
				{Tag: "-sortby"},
				{Tag: "-n"},
			},
			[]string{"ls"},
			listCommand(engine),
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	s.RegisterCommand(
		gshell.NewCommand(
			"load",
			"Load a new document into the engine",
			"load <document_path>",
			[]gshell.Argument{},
			[]string{},
			loadCommand(engine),
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	s.RegisterCommand(
		gshell.NewCommand(
			"find",
			"Find a document by name or id.\n"+
				"Display the document's id, name and path\n"+
				"Default search field: -name\n",
			"find <-id | -name> <value> || find <document_name>",
			[]gshell.Argument{
				{Tag: "-id"},
				{Tag: "-name"},
			},
			[]string{},
			findCommand(engine),
			func(args []string) (bool, string) {
				return true, ""
			},
		),
	)

	s.RegisterEarlyExecCommand(
		gshell.NewEarlyCommand(
			"engine-stats",
			"Displays the total number of documents and keys in the engine",
			"engine-stats",
			0,
			engineStatsCommand(engine),
		),
	)
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

func findCommand(engine *engine.Engine) func(s *gshell.Shell, args []string) gshell.Status {
	return func(s *gshell.Shell, args []string) gshell.Status {
		var doc internal.Document
		if len(args) > 2 {
			s.Write("Invalid number of arguments\n")
			return gshell.FAIL
		}

		if len(args) == 1 {
			if args[0] == "-id" || args[0] == "-name" {
				s.Write("Missing value for search field\n")
				return gshell.FAIL
			}
			doc = engine.GetDocumentByNameCopy(args[0])

		} else {
			var err string
			doc, err = getDocumentByField(engine, args[0], args[1])
			if err != "" {
				s.Write(err)
				return gshell.FAIL
			}
		}

		if doc == (internal.Document{}) {
			s.Write("Document not found\n")
		} else {
			s.Write("Document found:\n" + doc.String())
		}

		return gshell.OK
	}
}

func openCommand(engine *engine.Engine) func(s *gshell.Shell, args []string) gshell.Status {
	return func(s *gshell.Shell, args []string) gshell.Status {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			s.Write("Invalid document ID\n")
			return gshell.FAIL
		}
		if id < 0 || id >= len(engine.GetDocuments()) {
			s.Write("Document ID out of range, Docs Ids [0, " + strconv.Itoa(len(engine.GetDocuments())-1) + "]\n")
			return gshell.FAIL
		}
		doc := engine.GetDocumentByID(id)
		path := filepath.Join(doc.DirectoryPath, doc.Name+doc.Ext)
		s.Write("Opening document: " + doc.Name + ", path: " + path + "\n")
		err = openFile(path)
		if err != nil {
			s.Write("Error opening file\n")
			s.Write(err.Error())
			return gshell.FAIL
		}

		return gshell.OK
	}
}

func queryCommand(engine *engine.Engine) func(s *gshell.Shell, args []string) gshell.Status {
	return func(s *gshell.Shell, args []string) gshell.Status {
		res, err := engine.Query(args)

		if err != nil {
			s.Error(gshell.COMMAND_PREFIX, err.Error())
		}

		if res == nil || res.GetLength() == 0 {
			s.Write("No results found\n")
			return gshell.OK
		}

		for i := 0; i < res.GetLength(); i++ {
			s.Write(
				strconv.Itoa(engine.GetDocumentByID(res.At(i)).ID) +
					". " +
					engine.GetDocumentByID(res.At(i)).Name + "\n")
		}

		s.Write("Total documents found: " + strconv.Itoa(res.GetLength()) + "\n")

		return gshell.OK
	}
}

func loadCommand(engine *engine.Engine) func(s *gshell.Shell, args []string) gshell.Status {
	return func(s *gshell.Shell, args []string) gshell.Status {
		if len(args) != 1 {
			s.Write("Invalid number of arguments\n")
			return gshell.FAIL
		}

		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			s.Write("Path does not exist\n")
			return gshell.FAIL
		}

		lastId := engine.GetDocumentsSize()
		engine.LoadDirectory(args[0])

		totalLoaded := engine.GetDocumentsSize() - lastId

		if totalLoaded == 0 {
			s.Info(gshell.COMMAND_PREFIX, "No documents loaded")
		} else {
			for i := lastId; i < engine.GetDocumentsSize(); i++ {
				s.Write("Document: " + engine.GetDocumentByID(i).Name + " loaded\n")
			}
			s.Success(gshell.COMMAND_PREFIX, "Loaded "+strconv.Itoa(totalLoaded)+" documents")
		}

		return gshell.OK
	}
}

func engineStatsCommand(engine *engine.Engine) func(s *gshell.Shell) {
	return func(s *gshell.Shell) {
		s.WriteColored(gshell.COLOR_CYAN, "Engine stats:\n")
		s.WriteColored(gshell.COLOR_GREEN, "Total documents: ")
		s.Write(strconv.Itoa(engine.GetDocumentsSize()) + "\n")
		s.WriteColored(gshell.COLOR_GREEN, "Total keys: ")
		s.Write(strconv.Itoa(engine.GetIndexSize()) + "\n")
	}
}

func listCommand(engine *engine.Engine) func(s *gshell.Shell, args []string) gshell.Status {
	return func(s *gshell.Shell, args []string) gshell.Status {
		displayFields, limit, sortby, err := parseListArgs(args)
		if err != "" {
			s.Write(err + "\n")
			return gshell.FAIL
		}

		docs := getSortedDocuments(engine, sortby)
		if limit == -1 || limit > len(docs) {
			limit = len(docs)
		}

		displayDocuments(s, docs, displayFields, limit, sortby)
		return gshell.OK
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
	docs := engine.GetDocumentsCopy()

	switch sortby {
	case "name":
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].Name < docs[j].Name
		})
	case "path":
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].DirectoryPath < docs[j].DirectoryPath
		})
	case "-ext":
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].Ext < docs[j].Ext
		})
	default:
		return engine.GetDocuments()
	}
	return docs
}

func displayDocuments(s *gshell.Shell, docs []*internal.Document, displayFields []string, limit int, sortby string) {
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
