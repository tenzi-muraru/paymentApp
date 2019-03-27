package payment

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/x/bsonx"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tenzi-muraru/paymentApp/payment/model"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbName         = "paymentdb"
	collectionName = "payment"
)

// ErrPaymentNotFound - represents the error returned when no payment was found
var ErrPaymentNotFound = errors.New("Payment not found")

// Repository - contains methods for performing CRUD operations with payments
type Repository interface {
	GetAll() ([]model.Payment, error)
	GetByID(paymentID string) (*model.Payment, error)
	Add(payment model.Payment) (*model.Payment, error)
	Delete(paymentID string) error
	Update(payment model.Payment) error
}

// DBPaymentRepository - implements Repository using MongoDB for persistance
type DBPaymentRepository struct {
	collection *mongo.Collection
}

// NewDBPaymentRepository - creates a new Repository that connects to MongoDB
func NewDBPaymentRepository(uri string) (*DBPaymentRepository, error) {
	// Connect to MongoDB
	options := options.Client()
	options.SetConnectTimeout(3000 * time.Millisecond)
	options.SetServerSelectionTimeout(3000 * time.Millisecond)
	options.ApplyURI("mongodb://" + uri)

	mongoClient, err := mongo.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to MongoDB: %v", err)
	}

	if err = mongoClient.Connect(context.Background()); err != nil {
		return nil, fmt.Errorf("Cannot connect to MongoDB: %v", err)
	}

	c := mongoClient.Database(dbName).Collection(collectionName)
	_, err = c.Indexes().CreateOne(context.Background(), createUUIDIndexModel())
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to MongoDB: %v", err)
	}

	return &DBPaymentRepository{
		collection: c,
	}, nil
}

func createUUIDIndexModel() mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	}
}

// GetByID - retrieves the payment with the provided ID if it exists or error otherwise
func (r *DBPaymentRepository) GetByID(paymentID string) (*model.Payment, error) {
	filter := bson.M{"id": paymentID}
	res := r.collection.FindOne(context.Background(), filter)

	payment := &model.Payment{}
	if err := res.Decode(payment); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPaymentNotFound
		}
		return nil, fmt.Errorf("Get payment by id failed due to DB error: %v", err)
	}

	return payment, nil
}

// Add - adds a new Payment
func (r *DBPaymentRepository) Add(payment model.Payment) (*model.Payment, error) {
	payment.ID = uuid.New().String()
	_, err := r.collection.InsertOne(context.Background(), payment)
	if err != nil {
		return nil, fmt.Errorf("Add payment failed due to DB error: %v", err)
	}

	return &payment, nil
}

// GetAll - retrieves all payments
func (r *DBPaymentRepository) GetAll() ([]model.Payment, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("Get payments failed due to DB error: %v", err)
	}
	defer cursor.Close(context.Background())

	var payments []model.Payment
	for cursor.Next(context.Background()) {
		payment := &model.Payment{}
		if err = cursor.Decode(payment); err != nil {
			return nil, errors.Wrapf(err, "Get payments failed due to DB decode error")
		}
		payments = append(payments, *payment)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("Get payments failed due to DB error: %v", err)
	}

	return payments, nil
}

// Delete - deletes the payment with the provided ID
func (r *DBPaymentRepository) Delete(paymentID string) error {
	filter := bson.M{"id": paymentID}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("Delete failed due to DB error: %v", err)
	}

	return nil
}

// Update - updates the provided payment
func (r *DBPaymentRepository) Update(payment model.Payment) error {
	_, err := r.GetByID(payment.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"id": payment.ID}
	_, err = r.collection.ReplaceOne(context.Background(), filter, payment)
	if err != nil {
		return fmt.Errorf("Update failed due to DB error: %v", err)
	}

	return nil
}
