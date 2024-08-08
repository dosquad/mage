package helper

import "context"

type ContextKey string

const (
	DockerLocalPlatform ContextKey = "docker-local-platform"
)

func ContextDefaultValue[T any](ctx context.Context, key ContextKey, defaultValue T) T {
	if v, ok := ctx.Value(key).(T); ok {
		return v
	}

	return defaultValue
}
