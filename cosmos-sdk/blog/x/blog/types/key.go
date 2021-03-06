package types

const (
	// ModuleName is the name of the module
	ModuleName = "blog"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

const (
	PostPrefix    = "post-"
	QueryListPost = "list-post"
)

const (
	CommentPrefix = "comment-"
)
		