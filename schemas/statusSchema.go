package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type BotStatus struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Status string             `bson:"status"`
}
