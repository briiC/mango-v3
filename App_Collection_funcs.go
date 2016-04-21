package mango

// Collection - get collection by ckey
func (app *Application) Collection(ckey string) *Collection {
	c := app.collections[ckey]

	// Try to get plural forms of given ckey
	if c == nil {
		switch {
		case ckey == "Category":
			ckey = "Categories"
		case ckey == "Keyword":
			ckey = "Keywords"
		case ckey == "Tag":
			ckey = "Tags"
		}
		c = app.collections[ckey]
	}

	return c
}

// CollectionCount - total count of collections
func (app *Application) CollectionCount() int {
	return len(app.collections)
}

// CollectionPages - shorthand to get collection subitems
// Avoiding errors. Without it need to use: app.Collection(ckey).Get(csubkey)
// But there could be no ckey collection and that results in go-lang error
func (app *Application) CollectionPages(ckey, csubkey string) PageList {
	if c := app.Collection(ckey); c != nil {
		return c.Get(csubkey)
	}
	return nil
}
