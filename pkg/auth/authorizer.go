package auth

import (
	"context"
	"fmt"
	"strings"
)

const (
	serviceReadPermission     = "service:%s:read"
	serviceWritePermission    = "service:%s:write"
	serviceAdminPermission    = "service:%s:admin"
	capabilityReadPermission  = "capability:%s:read"
	capabilityWritePermission = "capability:%s:write"
	capabilityAdminPermission = "capability:%s:admin"
)

// AuthorizationError is the error structure used to store and return authorization errors
type AuthorizationError struct {
	missingPermissions []string
}

// Error parses the list of missing permissions into one string
func (ae AuthorizationError) Error() string {
	return fmt.Sprintf("Missing permissions %s", strings.Join(ae.missingPermissions, " "))
}

// GetMissingPermissions returns the minimum authorization permissions that are missing
func (ae AuthorizationError) GetMissingPermissions() []string {
	return ae.missingPermissions
}

// NewAuthorizationError returns a new instance of AuthorizationError
func NewAuthorizationError(missingPermissions []string) *AuthorizationError {
	return &AuthorizationError{missingPermissions: missingPermissions}
}

type permissionsKey struct{}

type getPermissionsFromContext func(context.Context) []string

// WithAuthorizationPermissions adds a slice of authorization permissions to the context
// If this function is used then the authorizer can use the default function to retrieve these permissions
// otherwise it needs to receive a custom function to retrieve permissions when created.
func WithAuthorizationPermissions(ctx context.Context, permissions []string) context.Context {
	return context.WithValue(ctx, permissionsKey{}, permissions)
}

// Authorizer is the struct used to perform authorization assertions
// Use NewAuthorizer to build.
type Authorizer struct {
	capability                string
	getPermissionsFromContext getPermissionsFromContext
}

// AssertServicesReaderOrCapabilityAdmin asserts if the context has at least read permissions for all the services sent.
// If the context has service:write, service:admin or capability:admin permissions this will work too.
func (a *Authorizer) AssertServicesReaderOrCapabilityAdmin(ctx context.Context, services []string, environment string) error {
	// If context has capability admin permission we don't check anything else
	if err := a.AssertCapabilityAdmin(ctx); err == nil {
		return nil
	}

	permissions := a.getPermissionsFromContext(ctx)

	missingPermissions := []string{}

	// For each service we check that context has at least read permissions.
	for _, service := range services {
		readPermission := fmt.Sprintf(serviceReadPermission, service)
		writePermission := fmt.Sprintf(serviceWritePermission, service)
		adminPermission := fmt.Sprintf(serviceAdminPermission, service)
		if !contains(permissions, readPermission) && !contains(permissions, writePermission) && !contains(permissions, adminPermission) {
			missingPermissions = append(missingPermissions, readPermission)
		}
	}

	if len(missingPermissions) > 0 {
		return NewAuthorizationError(missingPermissions)
	}

	return nil
}

// AssertServicesWriterOrCapabilityAdmin asserts if the context has at least write permissions for all the services sent.
// If the context has service:admin or capability:admin permissions this will work too.
func (a *Authorizer) AssertServicesWriterOrCapabilityAdmin(ctx context.Context, services []string, environment string) error {
	// If context has capability admin permission we don't check anything else
	if err := a.AssertCapabilityAdmin(ctx); err == nil {
		return nil
	}

	permissions := a.getPermissionsFromContext(ctx)

	missingPermissions := []string{}

	// For each service we check that context has at least write permissions.
	for _, service := range services {
		writePermission := fmt.Sprintf(serviceWritePermission, service)
		adminPermission := fmt.Sprintf(serviceAdminPermission, service)
		if !contains(permissions, writePermission) && !contains(permissions, adminPermission) {
			missingPermissions = append(missingPermissions, writePermission)
		}
	}

	if len(missingPermissions) > 0 {
		return NewAuthorizationError(missingPermissions)
	}

	return nil
}

// AssertServicesAdminOrCapabilityAdmin asserts if the context has admin permissions for all the services sent.
// If the context has capability:admin permissions this will work too.
func (a *Authorizer) AssertServicesAdminOrCapabilityAdmin(ctx context.Context, services []string) error {
	// If context has capability admin permission we don't check anything else
	if err := a.AssertCapabilityAdmin(ctx); err == nil {
		return nil
	}

	permissions := a.getPermissionsFromContext(ctx)

	missingPermissions := []string{}

	// For each service we check that context has admin permissions.
	for _, service := range services {
		adminPermission := fmt.Sprintf(serviceAdminPermission, service)
		if !contains(permissions, adminPermission) {
			missingPermissions = append(missingPermissions, adminPermission)
		}
	}

	if len(missingPermissions) > 0 {
		return NewAuthorizationError(missingPermissions)
	}

	return nil
}

// AssertCapabilityReaderOrCapabilityAdmin asserts if the context has at least read permissions for the capability
// If the context has capability:write or capability:admin permissions this will work too.
func (a *Authorizer) AssertCapabilityReaderOrCapabilityAdmin(ctx context.Context, environment string) error {
	// If context has capability admin permission we don't check anything else
	if err := a.AssertCapabilityAdmin(ctx); err == nil {
		return nil
	}

	permissions := a.getPermissionsFromContext(ctx)

	// We check context has at least read permissions.
	readPermission := fmt.Sprintf(capabilityReadPermission, a.capability)
	writePermission := fmt.Sprintf(capabilityWritePermission, a.capability)
	if !contains(permissions, readPermission) && !contains(permissions, writePermission) {
		return NewAuthorizationError([]string{readPermission})
	}

	return nil
}

// AssertCapabilityWriterOrCapabilityAdmin asserts if the context has at least write permissions for the capability
// If the context has capability:admin permissions this will work too.
func (a *Authorizer) AssertCapabilityWriterOrCapabilityAdmin(ctx context.Context, environment string) error {
	// If context has capability admin permission we don't check anything else
	if err := a.AssertCapabilityAdmin(ctx); err == nil {
		return nil
	}

	permissions := a.getPermissionsFromContext(ctx)

	writePermission := fmt.Sprintf(capabilityWritePermission, a.capability)
	if !contains(permissions, writePermission) {
		return NewAuthorizationError([]string{writePermission})
	}

	return nil
}

// Asserts if the context has admin permissions for the capability
func (a *Authorizer) AssertCapabilityAdmin(ctx context.Context) error {
	permissions := a.getPermissionsFromContext(ctx)

	adminPermission := fmt.Sprintf(capabilityAdminPermission, a.capability)
	if !contains(permissions, adminPermission) {
		return NewAuthorizationError([]string{adminPermission})
	}

	return nil
}

func contains(s []string, value string) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

func defaultGetPermissionsFromContext(ctx context.Context) []string {
	permissions, ok := ctx.Value(permissionsKey{}).([]string)

	if !ok {
		return []string{}
	}

	return permissions
}

// NewAuthorizer returns an instance of Authorizer for the desired capability
// customGetPermissionsFromContextFunc should only be used if the permissions weren't set with WithAuthorizationPermissions
// capability must be provided.
func NewAuthorizer(capability string, customGetPermissionsFromContextFunc getPermissionsFromContext) (*Authorizer, error) {
	getPermissionsFunc := defaultGetPermissionsFromContext

	if customGetPermissionsFromContextFunc != nil {
		getPermissionsFunc = customGetPermissionsFromContextFunc
	}

	if capability == "" {
		return nil, fmt.Errorf("Capability cannot be an empty string")
	}

	return &Authorizer{
		capability:                capability,
		getPermissionsFromContext: getPermissionsFunc,
	}, nil
}
