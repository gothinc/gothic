package gothic

type HookPoint int

const (
	BeforeCreateApp HookPoint = iota
	BeforeConfingLoad
	AfterConfigLoad
	BeforLoggerLoad
	AfterLoggerLoad
)

type hookFunc func() error
