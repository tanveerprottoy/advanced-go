package logex

import (
	"context"
	"log/slog"
	"os"
)

func Ex() {
	slog.Info("hello, world")

	// slog’s top-level functions use the default logger. We can get this logger explicitly,
	// and call its methods:
	logger := slog.Default()
	logger.Info("hello, world", "user", "txp")

	// Initially, slog’s output goes through the default log.Logger, producing the
	// output we’ve seen above. We can change the output by changing the handler used by the logger. slog comes with two built-in handlers. A TextHandler emits all log information in the form key=value. This program creates a new logger using a TextHandler and makes the same call to the Info method:
	// Everything has been turned into a key-value pair, with strings quoted as needed to preserve structure.
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("hello, world", "user", "txp")

	// For JSON output, install the built-in JSONHandler instead:
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("hello, world", "user", "txp")

	// for frequently executed log statements it may be more efficient to use
	// the Attr type and call the LogAttrs method. These work together to minimize 
	// memory allocations. There are functions for building Attrs out of strings, numbers, and other common types. This call to LogAttrs produces the same output as above, but does it faster:
	slog.LogAttrs(context.Background(), slog.LevelInfo, "hello, world",
		slog.String("user", "txp"))
}
