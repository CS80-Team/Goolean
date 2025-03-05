package textprocessing

type Processor interface {
	Process(token string) string
}

type DefaultProcessor struct {
	steps []Processor
}

func NewDefaultProcessor(steps ...Processor) *DefaultProcessor {
	return &DefaultProcessor{steps: steps}
}

func (p *DefaultProcessor) Process(token string) string {
	for _, step := range p.steps {
		token = step.Process(token)
	}
	return token
}
