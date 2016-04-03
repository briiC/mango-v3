package mango

// Collection - get collection by ckey
func (app *Application) Collection(ckey string) *Collection {
	return app.collections[ckey]
}

// CollectionCount - total count of collections
func (app *Application) CollectionCount() int {
	return len(app.collections)
}
