package main

import (
	"context"
	"log-service/data"
	"time"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.Background(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Time{},
	})
	if err != nil {
		return err
	}

	*resp = "Processed payload via RPC: " + payload.Name

	return nil
}
