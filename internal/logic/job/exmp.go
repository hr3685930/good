package job

import (
	"context"
	"fmt"
)

// Example Example
func Example(ctx context.Context, msg []byte) error {
	fmt.Println(string(msg))
	return nil
}

// Examples Examples
func Examples(ctx context.Context, msg []byte) error {
	fmt.Println(string(msg))
	return nil
}