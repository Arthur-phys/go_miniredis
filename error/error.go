package error

import (
	"fmt"
	"log/slog"
)

type Error struct {
	Content string
	Code    uint16
	From    error
}

func (e Error) Error() string {
	return fmt.Sprintf("[MiniRedisError-%d] %v\n", e.Code, e.Content)
}

func (e Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("MiniRedisError", e.Content),
		slog.Any("From", e.From),
		slog.Int("ErrorCode", int(e.Code)),
	)
}
