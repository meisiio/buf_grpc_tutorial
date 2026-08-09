package main

import (
	"testing"

	user "github.com/meisam/buf-tutorial/gen/go/proto/user/v1"
	"google.golang.org/protobuf/proto"
)

func TestProtobufSerialization(t *testing.T) {
	// 1. Create a new user message
	originalUser := &user.User{
		Id:    101,
		Name:  "Meisam",
		Email: "meisam@example.com",
	}

	// 2. Serialize (Marshal) the Go struct to Protobuf binary format
	data, err := proto.Marshal(originalUser)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}
	t.Logf("Serialized data size: %d bytes", len(data))

	// 3. Deserialize (Unmarshal) the binary format back into a Go struct
	newUser := &user.User{}
	err = proto.Unmarshal(data, newUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal data: %v", err)
	}

	// 4. Verify that data remained perfectly intact across serialization
	if newUser.GetId() != originalUser.GetId() || newUser.GetName() != originalUser.GetName() {
		t.Errorf("Deserialized data does not match original data")
	}
	
	t.Logf("Successfully unmarshaled user: ID=%d, Name=%s, Email=%s", newUser.GetId(), newUser.GetName(), newUser.GetEmail())
}
