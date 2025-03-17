package shell

type Command struct {
	Name        string
	Description string
	Usage       string
	Args        []Argument
	Aliases     []string
	Handler     func(s *Shell, args []string) Status
	Validator   func(args []string) (bool, string)
}

type Argument struct {
	Name        string
	Description string
	Required    bool
	Type        ArgType
	Default     string
	Validator   func(arg string) (bool, string)
}

type ArgType string

type EarlyCommand struct {
	Name        string
	Description string
	Usage       string
	Priority    int
	Handler     func(s *Shell)
}
