package schedule

// Handle is the email handler's implementation of the Handler interface.
func (h EmailHandlerNode) Handle() error {
	return nil
}

// Name is the email handlers implementation of the Node interface.
func (h EmailHandlerNode) Name() string {
	return h.name
}

// Next defines what's after this node completes.
func (h EmailHandlerNode) Next() ([]*NodeInstance, error) {
	return nil, nil
}
