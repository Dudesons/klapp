package api

//import (
//	consulapi "github.com/hashicorp/consul/api"
//	"testing"
//	"fmt"
//	"encoding/json"
//	"github.com/dudesons/klapp/pb"
//	"google.golang.org/grpc"
//	"golang.org/x/net/context"
//	"github.com/stretchr/testify/assert"
//	"time"
//	"github.com/dudesons/klapp/flip"
//)
//
//func newConsulClient(t *testing.T) *consulapi.KV{
//	config := consulapi.DefaultConfig()
//	config.Address = "127.0.0.1:8500"
//	consul, err := consulapi.NewClient(config)
//	if err != nil {
//		t.Error(err)
//	}
//
//	return consul.KV()
//}
//
//func setupConsul(t *testing.T, kv *consulapi.KV, fixtures string) {
//	flips := []*flip.Flip{}
//	json.Unmarshal([]byte(fixtures), &flips)
//
//	for _, i := range flips {
//		v, err := json.Marshal(i)
//		if err != nil {
//			t.Error(err)
//		}
//		kv.Put(&consulapi.KVPair{Key: fmt.Sprintf("klapp/flips/%s", *i.Name), Value: v}, nil)
//	}
//}
//
//func cleanConsul(t *testing.T, kv *consulapi.KV) {
//	_, err := kv.DeleteTree("klapp/flips/", nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//func TestFlipServer_IsFlip_Integration(t *testing.T) {
//	kv := newConsulClient(t)
//	address := "127.0.0.1:50051"
//	conn, err := grpc.Dial(address, grpc.WithInsecure())
//
//	if err != nil {
//		t.Errorf("did not connect: %v", err)
//	}
//
//	defer conn.Close()
//
//	c := pb.NewFlipClient(conn)
//
//	setupConsul(t, kv, flip_fixture)
//	time.Sleep(2)
//
//	for _, i := range flipNamesForBool {
//		r, err := c.IsFlip(context.Background(), &pb.FlipRequest{FeatureTag: i.Flip})
//		t.Log(r)
//		t.Log(err)
//		if i.Success {
//			assert.Equal(t, i.Success, r.Activated, fmt.Sprintf("should be equal, flip: %s, success: %b, activated: %b", i.Flip, i.Success, r.Activated))
//		}
//	}
//
//	cleanConsul(t, kv)
//}
