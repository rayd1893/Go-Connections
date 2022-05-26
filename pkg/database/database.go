package database

import (
	"context"
)

type Database interface {
	Connect(context.Context) error
	Close(context.Context) error
}
