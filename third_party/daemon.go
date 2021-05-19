package third_party

type (
	Daemon interface {
		Start() error
		Stop() error
	}
)
