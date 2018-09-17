package schedule

func (node EventNode) Name() string {
	return node.name
}

func (node EndNode) Name() string {
	return node.name
}

func (node StartNode) Name() string {
	return "start"
}
