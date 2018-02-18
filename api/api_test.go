package api

import (
	"testing"
	"github.com/dudesons/klapp/pb"
	"fmt"
	"github.com/dudesons/klapp/flip"
	"github.com/dudesons/klapp/config"
	"google.golang.org/grpc"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/stretchr/testify/assert"
	"github.com/kelseyhightower/envconfig"
)




func setupKV(t *testing.T, kv flip.FlipStore, fixtures string) {
	flips := []*flip.Flip{}
	json.Unmarshal([]byte(fixtures), &flips)

	for _, i := range flips {
		err := kv.Put(fmt.Sprintf("klapp/flips/%s", *i.Name), i)
		if err != nil {
			t.Error(err)
		}
	}
}

func cleanKV(t *testing.T, kv flip.FlipStore) {
	err := kv.Delete("klapp/flips/", true)
	if err != nil {
		fmt.Println(err)
	}
}

// TODO(Rework this test is an E2E)
func TestFlipServer_IsFlip_Integration(t *testing.T) {
	var conf config.KlappConfig
	err := envconfig.Process("klapp", &conf)
	if err != nil {
		t.Fatal(err)
	}

	kv, err := flip.NewFlipStore(&conf)
	address := "127.0.0.1:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		t.Errorf("did not connect: %v", err)
	}

	defer conn.Close()
	defer cleanKV(t, kv)

	c := pb.NewFlipClient(conn)

	setupKV(t, kv, flip.Flip_fixture)

	for _, i := range flip.FlipNamesForBool {
		r, err := c.IsFlip(context.Background(), &pb.FlipRequest{FeatureTag: i.Flip})
		t.Log(r)
		t.Log(err)
		if i.Success {
			assert.Equal(t, i.Success, r.Activated, fmt.Sprintf("should be equal, flip: %s, success: %b, activated: %b", i.Flip, i.Success, r.Activated))
		}
	}
}
