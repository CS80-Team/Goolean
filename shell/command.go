package shell

type Command struct {
	Name        string
	Description string
	Usage       string
	Handler     func(s *Shell, args []string) Status
}

type EarlyCommand struct {
	Name        string
	Description string
	Usage       string
	Handler     func(s *Shell)
}
