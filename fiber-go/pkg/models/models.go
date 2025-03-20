package models

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name" bson:"name"`
	Salary float64 `json:"salary" bson:"salary"`
	Age    int     `json:"age" bson:"age"`
}
