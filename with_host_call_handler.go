package vm

// WithHostCallHandler sets the host call handler for the VM.
func WithHostCallHandler(handler HostCallHandler) Option {
	return func(v *VM) {
		v.hostCallHandler = handler
	}
}
