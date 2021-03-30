package mdbwrapper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) (interface{}, error)
	UpdateOne(context.Context, interface{}, interface{}, *options.UpdateOptions) (UpdateResultHelper, error)
	DeleteOne(context.Context, interface{}) (int64, error)
}

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type UpdateResultHelper interface {
	UnmarshalBSON(b []byte) error
}

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect() error
	StartSession() (mongo.Session, error)
}

type MongoClient struct {
	c *mongo.Client
}

type MongoDatabase struct {
	db *mongo.Database
}

type MongoCollection struct {
	col *mongo.Collection
}

type MongoSingleResult struct {
	sr *mongo.SingleResult
}

type MongoUpdateResult struct {
	ur *mongo.UpdateResult
}

type MongoSession struct {
	mongo.Session
}

func NewClient(cfg *Config) (ClientHelper, error) {
	client, err := mongo.NewClient(options.Client().SetAuth(
		options.Credential{
			Username:   cfg.Username,
			Password:   cfg.Password,
			AuthSource: cfg.DatabaseName,
		}).ApplyURI(cfg.URI))

	return &MongoClient{c: client}, err
}

func NewDatabase(cfg *Config, client ClientHelper) DatabaseHelper {
	return client.Database(cfg.DatabaseName)
}

func (mc *MongoClient) Database(dbName string) DatabaseHelper {
	db := mc.c.Database(dbName)
	return &MongoDatabase{db: db}
}

func (mc *MongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.c.StartSession()
	return &MongoSession{session}, err
}

func (mc *MongoClient) Connect() error {
	return mc.c.Connect(nil)
}

func (md *MongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &MongoCollection{col: collection}
}

func (md *MongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &MongoClient{c: client}
}

func (mCol *MongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mCol.col.FindOne(ctx, filter)
	return &MongoSingleResult{sr: singleResult}
}

func (mCol *MongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mCol.col.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mCol *MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opt *options.UpdateOptions) (UpdateResultHelper, error) {
	r, err := mCol.col.UpdateOne(ctx, filter, update, opt)
	return &MongoUpdateResult{ur: r}, err
}

func (mCol *MongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mCol.col.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (sr *MongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (ur *MongoUpdateResult) UnmarshalBSON(b []byte) error {
	return ur.ur.UnmarshalBSON(b)
}
