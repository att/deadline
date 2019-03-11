package schedule

// Handle is the email handler's implementation of the Handler interface.
func (handler EmailHandlerNode) Handle(ctx *Context) {
	// cfg := config.GetEmailConfig()

	// var client *smtp.Client

}

// Name is the email handlers implementation of the Node interface.
func (handler EmailHandlerNode) Name() string {
	return handler.name
}

// Next for this type is simply defined. There's no logic computed. It
// return nil context.
func (handler EmailHandlerNode) Next() ([]*NodeInstance, *Context) {
	var ret []*NodeInstance
	return append(ret, handler.to), nil
}
