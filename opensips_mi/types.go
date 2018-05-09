package opensips_mi

// OpenSIPS MI Tree Node
type MINode struct {
	Name     string
	Value    string
	Attrs    map[string]string
	Children []*MINode
	ChildValues map[string]string
}

// OpenSIPS MI Client
type Client interface {
	Command(cmd string, args ... string) (*MINode, error)
	Close() error
}
