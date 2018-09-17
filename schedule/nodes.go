package schedule

// Name returns the name of the event node.
func (node EventNode) Name() string {
	return node.name
}

// Name returns the name of the end node.
func (node EndNode) Name() string {
	return node.name
}

// Name returns the name of the start node.
func (node StartNode) Name() string {
	return "start"
}
