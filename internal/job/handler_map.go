package job

type HandlerMap map[string]func() JobHandler

func NewHandlerMap() HandlerMap {
	return HandlerMap{
		"ProcessExample": func() JobHandler { return new(ProcessExample) },
	}
}
