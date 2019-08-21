package gql

import (
	"encoding/json"
    "log"
    "errors"
    "github.com/graphql-go/graphql"
    "bookmarks/pkg/database"
    "bookmarks/pkg/model"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/bson"
)

// Resolver : struct that contains a mongo client object
type Resolver struct{
    m *database.MongoClient
}

// GetBookmarksResolver : graphql resolver that returns all the bookmarks
func (r *Resolver) GetBookmarksResolver(p graphql.ResolveParams) (interface{}, error) {
    log.Println("=> Get all bookmarks")

    queryDocument := bson.D{{}}
    bookmarks, err := r.m.GetBookmarks("bookmarks", "bookmarks", queryDocument)

    if err != nil{
        log.Println(err)
        return nil, err
    }
    
    return bookmarks, nil
}

// GetBookmarksByNameResolver : graphql resolver that returns bookmarks by name
func (r *Resolver) GetBookmarksByNameResolver(p graphql.ResolveParams) (interface{}, error) {
    log.Println("=> Get bookmarks by name")

    name := p.Args["name"].(string)

    log.Printf("=> Bookmark name: %s", name)

    queryDocument := bson.D{ {"name", primitive.Regex{Pattern: name, Options: "i"} } }

    bookmarks, err := r.m.GetBookmarks("bookmarks", "bookmarks", queryDocument)

    if err != nil{
        log.Println(err)
        return nil, err
    }
    
    return bookmarks, nil
}

// GetBookmarksByTagsResolver : graphql resolver that returns bookmarks by tags
func (r *Resolver) GetBookmarksByTagsResolver(p graphql.ResolveParams) (interface{}, error) {
    log.Println("=> Get bookmarks by tags")

    var tags = []string{}
    encodedArgs, err := json.Marshal(p.Args["tags"])
    
    if err != nil{
        log.Println(err)
        return nil, err
    }

    err = json.Unmarshal(encodedArgs, &tags)

    if err != nil{
        log.Println(err)
        return nil, err
    }

    var tagsBSON = []bson.D{}

	for _, tag := range tags{
        log.Printf("=> Bookmark tag: %s", tag)
		tagBSON := bson.D{{"tags", primitive.Regex{Pattern: tag, Options: "i"}}}
		tagsBSON = append(tagsBSON, tagBSON)
	}

	queryDocument := bson.D{
		{"$and", 
			tagsBSON,
		},
	}

    bookmarks, err := r.m.GetBookmarks("bookmarks", "bookmarks", queryDocument)

    if err != nil{
        log.Println(err)
        return nil, err
    }
    
    return bookmarks, nil
}

// GetBookmarkByIDResolver : graphql resolver that returns a bookmark by id
func (r *Resolver) GetBookmarkByIDResolver(p graphql.ResolveParams) (interface{}, error){

    id, err := primitive.ObjectIDFromHex(p.Args["id"].(string))
    log.Printf("=> Get bookmark with id: %s", id.Hex())

    if err != nil{
        log.Println(err)
        return nil, err
    }

    idDocument := bson.D{{"_id", id}}

    bookmarkDocument, err := r.m.GetBookmarkByID("bookmarks", "bookmarks", idDocument)   

    if err != nil{
        log.Println(err)
        return nil, err
    }

    bookmark := model.Bookmark {
		ID: bookmarkDocument.ID.Hex(),
		Name: bookmarkDocument.Name,
		URL: bookmarkDocument.URL,
		Tags: bookmarkDocument.Tags,
	}
        
    return bookmark, nil
}

// DeleteBookmarkResolver : graphql resolver that deletes a bookmark by id
func (r *Resolver) DeleteBookmarkResolver(p graphql.ResolveParams) (interface{}, error) {
    log.Println("=> Delete bookmark")
    
    id := p.Args["id"].(string)

    log.Printf("=> Bookmark id: %s", id)

    objectID, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        log.Println(err)
        return nil, err
    }

    idDocument := bson.D{{"_id", objectID}}

    deletedCount, err := r.m.DeleteBookmark("bookmarks", "bookmarks", idDocument)
    
    if err != nil {
        log.Println(err)
        return nil, err
    }

    if deletedCount == 0{
        return nil, errors.New("mongo: deleted count is 0. Bookmark was not found") 
    }

    return id, nil
}

// CreateBookmarkResolver : graphql resolver that creates a new bookmark
func (r *Resolver) CreateBookmarkResolver(p graphql.ResolveParams) (interface{}, error) {
    log.Println("=> Create new bookmark")

    encodedArgs, err := json.Marshal(p.Args)

    if err != nil{
        log.Println(err)
        return nil, err
    }

    var b model.BookmarkBSON
    err = json.Unmarshal(encodedArgs, &b)

    if err != nil{
        log.Println(err)
        return nil, err
    }

    b.ID = primitive.NewObjectID()

    log.Printf("=> Bookmark id: %s", b.ID.Hex())

    bookmarkDocument, err := r.m.CreateBookmark("bookmarks", "bookmarks", &b)

    if err != nil{
        log.Println(err)
        return nil, err
    }

    bookmark := model.Bookmark {
		ID: bookmarkDocument.ID.Hex(),
		Name: bookmarkDocument.Name,
		URL: bookmarkDocument.URL,
		Tags: bookmarkDocument.Tags,
	}

    return bookmark, nil
}

// UpdateBookmarkResolver : graphql resolver that updates a bookmark
func (r *Resolver) UpdateBookmarkResolver(p graphql.ResolveParams) (interface{}, error) {
    log.Println("=> Update bookmark")

    id, err := primitive.ObjectIDFromHex(p.Args["id"].(string))
    
    if err != nil{
        log.Println(err)
        return nil, err
    }

    log.Printf("=> Bookmark id: %s", id.Hex())

    var fieldsQuery bson.D
    name, ok := p.Args["name"].(string)
    if ok {
        fieldsQuery = append(fieldsQuery, bson.E{"name", name})
    }

    url, ok := p.Args["url"].(string)
    if ok {
        fieldsQuery = append(fieldsQuery, bson.E{"url", url})
    }

    tagsInterface, ok := p.Args["tags"]
    if ok {

        var tags = []string{}
        encodedArgs, err := json.Marshal(tagsInterface)

        if err != nil{
            log.Println(err)
            return nil, err
        }

        err = json.Unmarshal(encodedArgs, &tags)

        if err != nil{
            log.Println(err)
            return nil, err
        }

        fieldsQuery = append(fieldsQuery, bson.E{"tags", tags})  
    }
    
    idDocument := bson.D{{"_id", id}}
    queryDocument := bson.D{{ "$set", fieldsQuery }}

    bookmarkDocument, err := r.m.UpdateBookmark("bookmarks", "bookmarks", idDocument, queryDocument)

    if err != nil{
        log.Println(err)
        return nil, err
    }

    bookmark := model.Bookmark {
		ID: bookmarkDocument.ID.Hex(),
		Name: bookmarkDocument.Name,
		URL: bookmarkDocument.URL,
		Tags: bookmarkDocument.Tags,
	}

    return bookmark, nil
}