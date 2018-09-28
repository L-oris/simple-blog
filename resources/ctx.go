package resources

import "context"

var CTX context.Context

func init() {
	CTX = context.Background()
}
