package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Bookmark : bookmark struct model
type Bookmark struct{
    ID string `json:"id,omitempty"`
    Name string `json:"name"`
    URL string `json:"url"`
    Tags []string `json:"tags"`
}

// BookmarkBSON : bookmark bson struct model
type BookmarkBSON struct{
    ID primitive.ObjectID `bson:"_id"`
    Name string `bson:"name"`
    URL string `bson:"url"`
    Tags []string `bson:"tags"`
}