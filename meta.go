package google_cloud_trace

import (
	"context"
	"os"

	"cloud.google.com/go/compute/metadata"
)

func isCloudRun() bool {
	ks := os.Getenv("K_SERVICE")
	return len(ks) > 0
}

func projectIDWithContext(ctx context.Context) string {
	if !isCloudRun() {
		return ""
	}
	v, _ := metadata.ProjectIDWithContext(ctx)
	return v
}
