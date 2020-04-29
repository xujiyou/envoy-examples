package main

import (
	"context"
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	sd "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"io"
	"io/ioutil"
	"log"
)

type MySDS struct{}

var typeUrl = "type.googleapis.com/envoy.api.v2.auth.Secret"

func (s *MySDS) DeltaSecrets(server sd.SecretDiscoveryService_DeltaSecretsServer) error {
	return nil
}

func (s *MySDS) StreamSecrets(server sd.SecretDiscoveryService_StreamSecretsServer) error {
	for {
		_, err := server.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		//log.Printf("server recv: %s", in)
		resp := getResp()
		_ = server.Send(&resp)
	}
}

func (s *MySDS) FetchSecrets(ctx context.Context, request *v2.DiscoveryRequest) (*v2.DiscoveryResponse, error) {
	log.Println(request)
	resp := getResp()
	return &resp, nil
}

func getResp() v2.DiscoveryResponse {
	buff := proto.NewBuffer(nil)
	buff.SetDeterministic(true)

	fileByte, _ := ioutil.ReadFile("/home/admin/k8s-cluster/envoy/ssl/cert/server.crt")
	log.Println(string(fileByte))

	resources := []types.Resource{
		&auth.Secret{
			Name: "my_secret",
			Type: &auth.Secret_TlsCertificate{
				TlsCertificate: &auth.TlsCertificate{
					CertificateChain: &core.DataSource{
						Specifier: &core.DataSource_InlineString{
							InlineString: string(fileByte),
						},
					},
					PrivateKey: &core.DataSource{
						Specifier: &core.DataSource_Filename{
							Filename: "/home/admin/k8s-cluster/envoy/ssl/cert/server.key",
						},
					},
				},
			},
		},
	}

	err := buff.Marshal(resources[0])

	if err != nil {
		log.Fatal(err)
	}

	return v2.DiscoveryResponse{
		VersionInfo: "1.0",
		Resources: []*any.Any{
			{
				TypeUrl: typeUrl,
				Value:   buff.Bytes(),
			},
		},
		TypeUrl: typeUrl,
	}
}
