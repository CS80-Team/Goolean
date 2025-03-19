package shell

type ArgType string

const (
    EMPTY_TAG = ""
)

type Command struct {
	Name         string
	Description  string
	Usage        string
	Args         []Argument
	Aliases      []string
	Handler      func(s *Shell, args []string) Status
	ValidateArgs func(args []string) (bool, string)
}

type EarlyCommand struct {
	Name        string
	Description string
	Usage       string
	Priority    int
	Handler     func(s *Shell)
}

type Argument struct {
	Name        string
    Tag         string // Tag is used for auto-completion
	Description string
	Required    bool
	Type        ArgType
	Default     string
}

func (arg *Argument) String() string {
	//return arg.Name
	return ""
}

func NewArgument(
	name string,
    tag string,
	description string,
	required bool,
	argType ArgType,
	defaultValue string,
) *Argument {
	return &Argument{
		Name:        name,
        Tag:         tag,
		Description: description,
		Required:    required,
		Type:        argType,
		Default:     defaultValue,
	}
}

func NewCommand(
	name string,
	description string,
	usage string,
	args []Argument,
	aliases []string,
	handler func(s *Shell, args []string) Status,
	validator func(args []string) (bool, string),
) *Command {
	return &Command{
		Name:         name,
		Description:  description,
		Usage:        usage,
		Args:         args,
		Aliases:      aliases,
		Handler:      handler,
		ValidateArgs: validator,
	}
}

func NewEarlyCommand(
	name string,
	description string,
	usage string,
	priority int, // Decide which to be displayed first (lower is first)
	handler func(s *Shell),
) EarlyCommand {
	return EarlyCommand{
		Name:        name,
		Description: description,
		Usage:       usage,
		Priority:    priority,
		Handler:     handler,
	}
}

func (cmd *Command) AddAlias(alias string) {
	cmd.Aliases = append(cmd.Aliases, alias)
}
