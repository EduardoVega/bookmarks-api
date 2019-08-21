package database

import (
	"log"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"bookmarks/pkg/model"
)

// MongoClient : creates a struct with a mongo client
type MongoClient struct{
	*mongo.Client
}

// Connect : creates a connection to mongo
func Connect(URI string) (*MongoClient, error){
	log.Println("=> Start DB connection")

	// Create connection
	ctxConnect, cancelConnect := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConnect()
	client, err := mongo.Connect(ctxConnect, options.Client().ApplyURI(URI))

	if err != nil {
		return nil, err
	}

	// Test connection
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()
	err = client.Ping(ctxPing, readpref.Primary())

	if err != nil {
		return nil, err
	}

	log.Println("=> DB connection successful")

	return &MongoClient{client}, nil
}

// GetBookmarks : gets all bookmarks in mongo
func (m *MongoClient) GetBookmarks(dbName, colName string, queryDocument bson.D) ([] *model.Bookmark, error){

	collection := m.Database(dbName).Collection(colName)

	cur, err := collection.Find(context.TODO(), queryDocument)

	if err != nil{
		return nil, err
	}

	var results []*model.Bookmark

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		
		// create a value into which the single document can be decoded
		var bb model.BookmarkBSON
		err := cur.Decode(&bb)
		if err != nil {
			return nil, err
		}

		b := model.Bookmark{
			ID: bb.ID.Hex(),
			Name: bb.Name,
			URL: bb.URL,
			Tags: bb.Tags,
		}

		results = append(results, &b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return results, nil
}

// GetBookmarkByID : get bookmark by id in mongo
func (m *MongoClient) GetBookmarkByID(dbName, colName string, idDocument bson.D) (*model.BookmarkBSON, error){

	collection := m.Database(dbName).Collection(colName)

	var bb model.BookmarkBSON

	err := collection.FindOne(context.TODO(), idDocument).Decode(&bb)
	
	if err != nil{
		log.Println(err.Error())

		if err  == mongo.ErrNoDocuments {
			 return nil, nil
		}
		
		return nil, err
	}

	return &bb, nil
}

// DeleteBookmark : delete a bookmark by id in mongo
func (m *MongoClient) DeleteBookmark(dbName, colName string, idDocument bson.D) (int64, error){

	collection := m.Database(dbName).Collection(colName)

	result, err := collection.DeleteOne(context.TODO(), idDocument)

	if err != nil {
		log.Println(err)
		return 0, err
	}

	log.Printf("=> Deleted count: %v", result.DeletedCount)
	
	return result.DeletedCount, nil
}

// CreateBookmark : create a new bookmark in mongo
func (m *MongoClient) CreateBookmark(dbName, colName string, bookmark *model.BookmarkBSON) (*model.BookmarkBSON, error){

	collection := m.Database(dbName).Collection(colName)

	result, err := collection.InsertOne(context.TODO(), bookmark)

	if err != nil{
		return nil, err
	}

	bookmarkDocument, err := m.GetBookmarkByID(dbName, colName, bson.D{{"_id", result.InsertedID}})

	if err != nil{
		return nil, err
	}

	return bookmarkDocument, nil
}

// UpdateBookmark : update a bookmark in mongo
func (m *MongoClient) UpdateBookmark(dbName, colName string, idDocument bson.D, queryDocument bson.D) (*model.BookmarkBSON, error){

	collection := m.Database(dbName).Collection(colName)

	_, err := collection.UpdateOne(
		context.TODO(), 
		idDocument, 
		queryDocument,
	)

	if err != nil{
		return nil, err
	}

	bookmarkDocument, err := m.GetBookmarkByID(dbName, colName, idDocument)

	if err != nil{
		return nil, err
	}

	return bookmarkDocument, nil
}

