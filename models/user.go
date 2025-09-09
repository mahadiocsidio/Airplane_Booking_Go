package models

import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
	

type User struct{
	ID			primitive.ObjectID	`bson:"_id,omitempty" json:"id,omitempty"`
	Name		string				`bson:"name" json:"name"`
	Email		string				`bson:"email" json:"email"`
	Password	string				`bson:"password" json:"password"`
	Phone		string				`bson:"phone" json:"phone"`
	Role		string				`bson:"role,omitempty" json:"role,omitempty"`
	CreatedAt	time.Time			`bson:"created_at" json:"created_at"`
	UpdatedAt	time.Time			`bson:"updated_at" json:"updated_at"`
}