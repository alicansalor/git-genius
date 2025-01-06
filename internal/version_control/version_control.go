package versioncontrol

import "context"

type PRCreator interface {
	CreatePR(ctx context.Context, title, body, head, base string) (string, error)
}
