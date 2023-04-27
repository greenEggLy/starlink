package globaldata

type Cmdline uint32

const (
	update Cmdline = 0
	getsys Cmdline = 1
	get    Cmdline = 2
)

func (cmd Cmdline) String() string {
	switch cmd {
	case update:
		return "update"
	case getsys:
		return "getsys"
	case get:
		return "get"
	default:
		return "unknown"
	}
}
