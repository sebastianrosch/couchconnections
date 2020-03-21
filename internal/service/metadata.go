package service

import "context"

type methodInfoKey struct{}

// MethodInfo struct holds metadata about the current request
type MethodInfo struct {
	// FullName the full name of the method that is being invoked
	FullName string
}

// ServiceMetadata struct provides access to the metadata of the current request
type ServiceMetadata struct { // nolint
}

// GetMethodInfo returns the method information from the context
func (s *ServiceMetadata) GetMethodInfo(ctx context.Context) *MethodInfo {
	metadata := ctx.Value(methodInfoKey{})
	if metadata != nil {
		return metadata.(*MethodInfo)
	}
	return nil
}

// WithMethodInfo stores the method information into the context
func (s *ServiceMetadata) WithMethodInfo(ctx context.Context, info *MethodInfo) context.Context {
	return context.WithValue(ctx, methodInfoKey{}, info)
}

// NewMetadata returns a new Metadata instance
func NewMetadata() *ServiceMetadata {
	return &ServiceMetadata{}
}
