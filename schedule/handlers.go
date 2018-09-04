package schedule

func (h EmailHandlerNode) Handle() error {
	return nil
}

func (h EmailHandlerNode) Name() string {
	return h.name
}
