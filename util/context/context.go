package context

import (
	"context"
	"movie-booking/core/model"
)

// Context keys
const (
	UserID     = "userID"
	Datastore  = "datastore"
	GormAccessor = "gormaccessor"
)

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserID).(uint)
	return userID, ok
}

// SetUserID sets user ID in context
func SetUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, UserID, userID)
}

// GetDataStore extracts datastore from context
func GetDataStore(ctx context.Context) (model.DataStore, bool) {
	ds, ok := ctx.Value(Datastore).(model.DataStore)
	return ds, ok
}

// SetDataStore sets datastore in context
func SetDataStore(ctx context.Context, ds model.DataStore) context.Context {
	return context.WithValue(ctx, Datastore, ds)
}
