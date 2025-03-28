package main

import (
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/CS80-Team/Goolean/internal/service"
	"github.com/CS80-Team/Goolean/internal/structures/factory"
	"github.com/CS80-Team/Goolean/internal/transport/file"
	"github.com/CS80-Team/Goolean/internal/transport/query"
	"github.com/CS80-Team/gshell/pkg/gshell"
	"github.com/chzyer/readline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/CS80-Team/Goolean/internal/engine"
	"github.com/CS80-Team/Goolean/internal/engine/tokenizer"
	"github.com/CS80-Team/Goolean/internal/textprocessing"
)

// -service ip:port
// -shell

func main() {
	ipPort := flag.String("service", "localhost:8888", "ip:port for the service to run on")
	shell := flag.Bool("shell", false, "start shell interface")

	flag.Parse()

	engine := engine.NewEngine(
		textprocessing.NewDefaultProcessor(
			textprocessing.NewNormalizer(),
			textprocessing.NewStemmer(),
			textprocessing.NewStopWordRemover(),
		),
		tokenizer.NewDelimiterManager(
			&map[rune]struct{}{
				' ': {},

				',':  {},
				'?':  {},
				'!':  {},
				'.':  {},
				';':  {},
				':':  {},
				'\\': {},

				'(': {},
				')': {},
				'[': {},
				']': {},
				'{': {},
				'}': {},

				'=': {},
				'+': {},
				'-': {},
				'*': {},
				'/': {},
				'%': {},
				'^': {},
			},
		),
		*engine.NewIndexManager(factory.NewSkipPointerListFactory[int]()),
	)

	engine.LoadDirectory(filepath.Join(filepath.Base("."), "dataset"))

	if *shell {
		startShell(engine)
	} else {
		startService(ipPort, engine)
	}
}

func startShell(engine *engine.Engine) {
	stdin, stdinW := readline.NewFillableStdin(os.Stdin)

	s := gshell.NewShell(
		stdin,
		stdinW,
		os.Stdout,
		os.Stdout,
		gshell.SHELL_PROMPT,
		".shell_history",
		gshell.NewLogger("shell.log"),
	)

	RegisterCommands(s, engine)

	s.Run("Welcome to the Goolean search engine shell, type `help` for list of commands\n")
}

func startService(ipPort *string, engine *engine.Engine) {
	lis, err := net.Listen("tcp", *ipPort)
	if err != nil {
		log.Fatal("Could not listen on port")
	}

	grpcServer := grpc.NewServer()
	queryServer := service.NewQueryServer(engine)
	fileServer := service.NewFileServer("recivedFiles", engine)

	reflection.Register(grpcServer)
	query.RegisterQueryServer(grpcServer, queryServer)
	file.RegisterFileServiceServer(grpcServer, fileServer)

	log.Println("Starting server...")
	_ = grpcServer.Serve(lis)
}
