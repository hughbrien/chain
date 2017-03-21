// Code generated by protoc-gen-go.
// source: config.proto
// DO NOT EDIT!

/*
Package config is a generated protocol buffer package.

It is generated from these files:
	config.proto

It has these top-level messages:
	Config
	BlockSigner
*/
package config

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import bc "chain/protocol/bc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Config struct {
	Id                   string         `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	IsSigner             bool           `protobuf:"varint,2,opt,name=is_signer,json=isSigner" json:"is_signer,omitempty"`
	IsGenerator          bool           `protobuf:"varint,3,opt,name=is_generator,json=isGenerator" json:"is_generator,omitempty"`
	BlockchainId         *bc.ProtoHash  `protobuf:"bytes,4,opt,name=blockchain_id,json=blockchainId" json:"blockchain_id,omitempty"`
	GeneratorUrl         string         `protobuf:"bytes,5,opt,name=generator_url,json=generatorUrl" json:"generator_url,omitempty"`
	GeneratorAccessToken string         `protobuf:"bytes,6,opt,name=generator_access_token,json=generatorAccessToken" json:"generator_access_token,omitempty"`
	BlockHsmUrl          string         `protobuf:"bytes,7,opt,name=block_hsm_url,json=blockHsmUrl" json:"block_hsm_url,omitempty"`
	BlockHsmAccessToken  string         `protobuf:"bytes,8,opt,name=block_hsm_access_token,json=blockHsmAccessToken" json:"block_hsm_access_token,omitempty"`
	ConfiguredAt         uint64         `protobuf:"varint,9,opt,name=configured_at,json=configuredAt" json:"configured_at,omitempty"`
	BlockPub             string         `protobuf:"bytes,10,opt,name=block_pub,json=blockPub" json:"block_pub,omitempty"`
	Signers              []*BlockSigner `protobuf:"bytes,11,rep,name=signers" json:"signers,omitempty"`
	Quorum               uint32         `protobuf:"varint,12,opt,name=quorum" json:"quorum,omitempty"`
	MaxIssuanceWindowMs  uint64         `protobuf:"varint,13,opt,name=max_issuance_window_ms,json=maxIssuanceWindowMs" json:"max_issuance_window_ms,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Config) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Config) GetIsSigner() bool {
	if m != nil {
		return m.IsSigner
	}
	return false
}

func (m *Config) GetIsGenerator() bool {
	if m != nil {
		return m.IsGenerator
	}
	return false
}

func (m *Config) GetBlockchainId() *bc.ProtoHash {
	if m != nil {
		return m.BlockchainId
	}
	return nil
}

func (m *Config) GetGeneratorUrl() string {
	if m != nil {
		return m.GeneratorUrl
	}
	return ""
}

func (m *Config) GetGeneratorAccessToken() string {
	if m != nil {
		return m.GeneratorAccessToken
	}
	return ""
}

func (m *Config) GetBlockHsmUrl() string {
	if m != nil {
		return m.BlockHsmUrl
	}
	return ""
}

func (m *Config) GetBlockHsmAccessToken() string {
	if m != nil {
		return m.BlockHsmAccessToken
	}
	return ""
}

func (m *Config) GetConfiguredAt() uint64 {
	if m != nil {
		return m.ConfiguredAt
	}
	return 0
}

func (m *Config) GetBlockPub() string {
	if m != nil {
		return m.BlockPub
	}
	return ""
}

func (m *Config) GetSigners() []*BlockSigner {
	if m != nil {
		return m.Signers
	}
	return nil
}

func (m *Config) GetQuorum() uint32 {
	if m != nil {
		return m.Quorum
	}
	return 0
}

func (m *Config) GetMaxIssuanceWindowMs() uint64 {
	if m != nil {
		return m.MaxIssuanceWindowMs
	}
	return 0
}

type BlockSigner struct {
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token,omitempty"`
	Pubkey      []byte `protobuf:"bytes,2,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
	Url         string `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
}

func (m *BlockSigner) Reset()                    { *m = BlockSigner{} }
func (m *BlockSigner) String() string            { return proto.CompactTextString(m) }
func (*BlockSigner) ProtoMessage()               {}
func (*BlockSigner) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BlockSigner) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *BlockSigner) GetPubkey() []byte {
	if m != nil {
		return m.Pubkey
	}
	return nil
}

func (m *BlockSigner) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func init() {
	proto.RegisterType((*Config)(nil), "config.Config")
	proto.RegisterType((*BlockSigner)(nil), "config.BlockSigner")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 408 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0x41, 0x8f, 0xd3, 0x30,
	0x10, 0x85, 0x95, 0x66, 0xc9, 0xb6, 0x93, 0x04, 0x21, 0x17, 0x55, 0x16, 0x70, 0x08, 0xdd, 0x4b,
	0x2e, 0xb4, 0x52, 0x97, 0x3f, 0xb0, 0x70, 0x60, 0xf7, 0x80, 0xb4, 0x0a, 0x20, 0x24, 0x2e, 0x96,
	0xe3, 0x84, 0xc6, 0x6a, 0x13, 0x97, 0x4c, 0xac, 0x5d, 0xfe, 0x3c, 0x42, 0x9e, 0xa4, 0x4d, 0x7b,
	0xcb, 0xbc, 0xf7, 0xe5, 0x79, 0x3c, 0x1e, 0x88, 0x94, 0x69, 0x7e, 0xeb, 0xed, 0xea, 0xd0, 0x9a,
	0xce, 0xb0, 0xa0, 0xaf, 0xde, 0xbc, 0x53, 0x95, 0xd4, 0xcd, 0x9a, 0x44, 0x65, 0xf6, 0xeb, 0x5c,
	0xad, 0x2b, 0x89, 0x55, 0x4f, 0x2d, 0xff, 0xf9, 0x10, 0x7c, 0x26, 0x90, 0xbd, 0x84, 0x89, 0x2e,
	0xb8, 0x97, 0x78, 0xe9, 0x2c, 0x9b, 0xe8, 0x82, 0xbd, 0x85, 0x99, 0x46, 0x81, 0x7a, 0xdb, 0x94,
	0x2d, 0x9f, 0x24, 0x5e, 0x3a, 0xcd, 0xa6, 0x1a, 0xbf, 0x51, 0xcd, 0xde, 0x43, 0xa4, 0x51, 0x6c,
	0xcb, 0xa6, 0x6c, 0x65, 0x67, 0x5a, 0xee, 0x93, 0x1f, 0x6a, 0xfc, 0x72, 0x94, 0xd8, 0x06, 0xe2,
	0x7c, 0x6f, 0xd4, 0x8e, 0xce, 0x17, 0xba, 0xe0, 0x57, 0x89, 0x97, 0x86, 0x9b, 0x78, 0x95, 0xab,
	0xd5, 0xa3, 0x3b, 0xfc, 0x5e, 0x62, 0x95, 0x45, 0x23, 0xf3, 0x50, 0xb0, 0x1b, 0x88, 0x4f, 0x99,
	0xc2, 0xb6, 0x7b, 0xfe, 0x82, 0xda, 0x89, 0x4e, 0xe2, 0x8f, 0x76, 0xcf, 0x3e, 0xc2, 0x62, 0x84,
	0xa4, 0x52, 0x25, 0xa2, 0xe8, 0xcc, 0xae, 0x6c, 0x78, 0x40, 0xf4, 0xeb, 0x93, 0x7b, 0x47, 0xe6,
	0x77, 0xe7, 0xb1, 0xe5, 0xd0, 0x8e, 0xa8, 0xb0, 0xa6, 0xe8, 0x6b, 0x82, 0x43, 0x12, 0xef, 0xb1,
	0x76, 0xc9, 0xb7, 0xb0, 0x18, 0x99, 0x8b, 0xe4, 0x29, 0xc1, 0xf3, 0x23, 0x7c, 0x1e, 0x7c, 0x03,
	0x71, 0x3f, 0x6a, 0xdb, 0x96, 0x85, 0x90, 0x1d, 0x9f, 0x25, 0x5e, 0x7a, 0x95, 0x45, 0xa3, 0x78,
	0xd7, 0xb9, 0x61, 0xf6, 0xc9, 0x07, 0x9b, 0x73, 0xa0, 0xb0, 0x29, 0x09, 0x8f, 0x36, 0x67, 0x1f,
	0xe0, 0xba, 0x1f, 0x33, 0xf2, 0x30, 0xf1, 0xd3, 0x70, 0x33, 0x5f, 0x0d, 0x4f, 0xf9, 0xc9, 0x21,
	0xfd, 0xc8, 0xb3, 0x23, 0xc3, 0x16, 0x10, 0xfc, 0xb1, 0xa6, 0xb5, 0x35, 0x8f, 0x12, 0x2f, 0x8d,
	0xb3, 0xa1, 0x72, 0xdd, 0xd7, 0xf2, 0x59, 0x68, 0x44, 0x2b, 0x1b, 0x55, 0x8a, 0x27, 0xdd, 0x14,
	0xe6, 0x49, 0xd4, 0xc8, 0x63, 0xea, 0x68, 0x5e, 0xcb, 0xe7, 0x87, 0xc1, 0xfc, 0x49, 0xde, 0x57,
	0x5c, 0xfe, 0x82, 0xf0, 0xec, 0x10, 0xf7, 0xae, 0x17, 0xf7, 0xee, 0xd7, 0x21, 0x94, 0x67, 0xf7,
	0x5d, 0x40, 0x70, 0xb0, 0xf9, 0xae, 0xfc, 0x4b, 0x4b, 0x11, 0x65, 0x43, 0xc5, 0x5e, 0x81, 0xef,
	0xc6, 0xea, 0xd3, 0x1f, 0xee, 0x33, 0x0f, 0x68, 0xc7, 0x6e, 0xff, 0x07, 0x00, 0x00, 0xff, 0xff,
	0x18, 0xa4, 0x3a, 0x48, 0x99, 0x02, 0x00, 0x00,
}