// Copyright © 2022 Weald Technology Trading.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validatorcredentialsset

import (
	"context"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	capella "github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/ethdo/beacon"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

func TestGenerateOperationFromMnemonicAndPath(t *testing.T) {
	ctx := context.Background()

	require.NoError(t, e2types.InitBLS())

	chainInfo := &beacon.ChainInfo{
		Version: 1,
		Validators: []*beacon.ValidatorInfo{
			{
				Index:                 0,
				Pubkey:                phase0.BLSPubKey{0xb3, 0x84, 0xf7, 0x67, 0xd9, 0x64, 0xe1, 0x00, 0xc8, 0xa9, 0xb2, 0x10, 0x18, 0xd0, 0x8c, 0x25, 0xff, 0xeb, 0xae, 0x26, 0x8b, 0x3a, 0xb6, 0xd6, 0x10, 0x35, 0x38, 0x97, 0x54, 0x19, 0x71, 0x72, 0x6d, 0xbf, 0xc3, 0xc7, 0x46, 0x38, 0x84, 0xc6, 0x8a, 0x53, 0x15, 0x15, 0xaa, 0xb9, 0x4c, 0x87},
				WithdrawalCredentials: []byte{0x00, 0x8b, 0xa1, 0xcc, 0x4b, 0x09, 0x1b, 0x91, 0xc1, 0x20, 0x2b, 0xba, 0x3f, 0x50, 0x80, 0x75, 0xd6, 0xff, 0x56, 0x5c, 0x77, 0xe5, 0x59, 0xf0, 0x80, 0x3c, 0x07, 0x92, 0xe0, 0x30, 0x2b, 0xf1},
			},
			{
				Index:                 1,
				Pubkey:                phase0.BLSPubKey{0xb3, 0xd8, 0x9e, 0x2f, 0x29, 0xc7, 0x12, 0xc6, 0xa9, 0xf8, 0xe5, 0xa2, 0x69, 0xb9, 0x76, 0x17, 0xc4, 0xa9, 0x4d, 0xd6, 0xf6, 0x66, 0x2a, 0xb3, 0xb0, 0x7c, 0xe9, 0xe5, 0x43, 0x45, 0x73, 0xf1, 0x5b, 0x5c, 0x98, 0x8c, 0xd1, 0x4b, 0xbd, 0x58, 0x04, 0xf7, 0x71, 0x56, 0xa8, 0xaf, 0x1c, 0xfa},
				WithdrawalCredentials: []byte{0x00, 0x78, 0x6c, 0xb0, 0x2e, 0xd2, 0x8e, 0x5f, 0xbb, 0x1f, 0x7f, 0x9e, 0x93, 0x1a, 0x2b, 0x72, 0x69, 0x29, 0x06, 0xe6, 0xb1, 0x2c, 0xe4, 0x64, 0x39, 0x75, 0xe3, 0x2b, 0x51, 0x76, 0x91, 0xf2},
			},
		},
		GenesisValidatorsRoot: phase0.Root{},
		Epoch:                 1,
		CurrentForkVersion:    phase0.Version{},
	}

	tests := []struct {
		name     string
		command  *command
		expected []*capella.SignedBLSToExecutionChange
		err      string
	}{
		{
			name: "MnemonicInvalid",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon",
				path:                 "m/12381/3600/0/0/0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			err: "mnemonic is invalid",
		},
		{
			name: "PathInvalidNoIndex",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				path:                 "m/12381/3600/0/0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			err: "path m/12381/3600/0/0 does not match EIP-2334 format for a validator",
		},
		{
			name: "NoPathProvided",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			err: "no validator path provided",
		},
		{
			name: "PathInvlaidIndexNot2334Format",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				path:                 "1",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			err: "path 1 does not match EIP-2334 format for a validator",
		},
		{
			name: "WithdrawalAddressNo0xPrefix",
			command: &command{mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				path:                 "m/12381/3600/0/0/0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			err: "failed to generate operation from seed and path: invalid withdrawal address: withdrawal address 8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15 does not contain a 0x prefix",
		},
		{
			name: "WithdrawalAddressInvalid",
			command: &command{mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				path:                 "m/12381/3600/0/0/0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac",
			},
			err: "failed to generate operation from seed and path: invalid withdrawal address: withdrawal address must be exactly 20 bytes in length",
		},
		{
			name: "WithdrawalAddressMissing",
			command: &command{mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				path:             "m/12381/3600/0/0/0",
				chainInfo:        chainInfo,
				signedOperations: make([]*capella.SignedBLSToExecutionChange, 0),
			},
			err: "failed to generate operation from seed and path: invalid withdrawal address: no withdrawal address provided",
		},
		{
			name: "InvalidWithdrawalAddressNotHex",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				path:                 "m/12381/3600/0/0/0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0xrc1Ff978036F2e9d7CC382Eff7B4c8c53C22acaa",
			},
			err: "failed to generate operation from seed and path: invalid withdrawal address: failed to obtain execution address: encoding/hex: invalid byte: U+0072 'r'",
		},
		{
			name: "Good",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				path:                 "m/12381/3600/0/0/0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			expected: []*capella.SignedBLSToExecutionChange{
				{
					Message: &capella.BLSToExecutionChange{
						ValidatorIndex:     0,
						FromBLSPubkey:      phase0.BLSPubKey{0x99, 0xb1, 0xf1, 0xd8, 0x4d, 0x76, 0x18, 0x54, 0x66, 0xd8, 0x6c, 0x34, 0xbd, 0xe1, 0x10, 0x13, 0x16, 0xaf, 0xdd, 0xae, 0x76, 0x21, 0x7a, 0xa8, 0x6c, 0xd0, 0x66, 0x97, 0x9b, 0x19, 0x85, 0x8c, 0x2c, 0x9d, 0x9e, 0x56, 0xee, 0xbc, 0x1e, 0x06, 0x7a, 0xc5, 0x42, 0x77, 0xa6, 0x17, 0x90, 0xdb},
						ToExecutionAddress: bellatrix.ExecutionAddress{0x8c, 0x1f, 0xf9, 0x78, 0x03, 0x6f, 0x2e, 0x9d, 0x7c, 0xc3, 0x82, 0xef, 0xf7, 0xb4, 0xc8, 0xc5, 0x3c, 0x22, 0xac, 0x15},
					},
					Signature: phase0.BLSSignature{0xb7, 0x8a, 0x05, 0xba, 0xd9, 0x27, 0xfc, 0x89, 0x6f, 0x14, 0x06, 0xb3, 0x2d, 0x64, 0x4a, 0xe1, 0x69, 0xce, 0xcd, 0x89, 0x86, 0xc1, 0xef, 0x8c, 0x0d, 0x03, 0x7d, 0x70, 0x86, 0xf8, 0x5f, 0x13, 0xe1, 0xe1, 0x88, 0xb4, 0x30, 0x96, 0x43, 0xa2, 0xc1, 0x3f, 0xfe, 0xfb, 0x0a, 0xe8, 0x05, 0x11, 0x09, 0x98, 0x53, 0xa0, 0x58, 0x1f, 0x4b, 0x2b, 0xd2, 0xe1, 0x45, 0x41, 0x04, 0x79, 0x01, 0xe2, 0x2a, 0x94, 0x0a, 0x9c, 0x7e, 0x3a, 0xc0, 0xa8, 0x82, 0xd1, 0xa8, 0xaf, 0x6b, 0xfa, 0xea, 0x81, 0x3a, 0x6a, 0x6b, 0xe7, 0x21, 0xf9, 0x26, 0x22, 0x04, 0xaa, 0x9d, 0xa4, 0xe4, 0x77, 0x27, 0xd0},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.command.generateOperationFromMnemonicAndPath(ctx)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				// fmt.Printf("%v\n", test.command.signedOperations)
				require.Equal(t, test.expected, test.command.signedOperations)
			}
		})
	}
}

func TestGenerateOperationFromMnemonicAndValidator(t *testing.T) {
	ctx := context.Background()

	require.NoError(t, e2types.InitBLS())

	chainInfo := &beacon.ChainInfo{
		Version: 1,
		Validators: []*beacon.ValidatorInfo{
			{
				Index:                 0,
				Pubkey:                phase0.BLSPubKey{0xb3, 0x84, 0xf7, 0x67, 0xd9, 0x64, 0xe1, 0x00, 0xc8, 0xa9, 0xb2, 0x10, 0x18, 0xd0, 0x8c, 0x25, 0xff, 0xeb, 0xae, 0x26, 0x8b, 0x3a, 0xb6, 0xd6, 0x10, 0x35, 0x38, 0x97, 0x54, 0x19, 0x71, 0x72, 0x6d, 0xbf, 0xc3, 0xc7, 0x46, 0x38, 0x84, 0xc6, 0x8a, 0x53, 0x15, 0x15, 0xaa, 0xb9, 0x4c, 0x87},
				WithdrawalCredentials: []byte{0x00, 0x8b, 0xa1, 0xcc, 0x4b, 0x09, 0x1b, 0x91, 0xc1, 0x20, 0x2b, 0xba, 0x3f, 0x50, 0x80, 0x75, 0xd6, 0xff, 0x56, 0x5c, 0x77, 0xe5, 0x59, 0xf0, 0x80, 0x3c, 0x07, 0x92, 0xe0, 0x30, 0x2b, 0xf1},
			},
			{
				Index:                 1,
				Pubkey:                phase0.BLSPubKey{0xb3, 0xd8, 0x9e, 0x2f, 0x29, 0xc7, 0x12, 0xc6, 0xa9, 0xf8, 0xe5, 0xa2, 0x69, 0xb9, 0x76, 0x17, 0xc4, 0xa9, 0x4d, 0xd6, 0xf6, 0x66, 0x2a, 0xb3, 0xb0, 0x7c, 0xe9, 0xe5, 0x43, 0x45, 0x73, 0xf1, 0x5b, 0x5c, 0x98, 0x8c, 0xd1, 0x4b, 0xbd, 0x58, 0x04, 0xf7, 0x71, 0x56, 0xa8, 0xaf, 0x1c, 0xfa},
				WithdrawalCredentials: []byte{0x00, 0x78, 0x6c, 0xb0, 0x2e, 0xd2, 0x8e, 0x5f, 0xbb, 0x1f, 0x7f, 0x9e, 0x93, 0x1a, 0x2b, 0x72, 0x69, 0x29, 0x06, 0xe6, 0xb1, 0x2c, 0xe4, 0x64, 0x39, 0x75, 0xe3, 0x2b, 0x51, 0x76, 0x91, 0xf2},
			},
		},
		GenesisValidatorsRoot: phase0.Root{},
		Epoch:                 1,
		CurrentForkVersion:    phase0.Version{},
	}

	tests := []struct {
		name     string
		command  *command
		expected []*capella.SignedBLSToExecutionChange
		err      string
	}{
		{
			name: "MnemonicInvalid",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon",
				validator:            "0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			err: "mnemonic is invalid",
		},
		{
			name: "ValidatorMissing",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			err: "no validator specified",
		},
		{
			name: "WithdrawalAddressMissing",
			command: &command{
				mnemonic:         "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:        "0",
				chainInfo:        chainInfo,
				signedOperations: make([]*capella.SignedBLSToExecutionChange, 0),
			},
			err: "invalid withdrawal address: no withdrawal address provided",
		},
		{
			name: "InvalidWithdrawalAddressLen",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:            "0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac",
			},
			err: "invalid withdrawal address: withdrawal address must be exactly 20 bytes in length",
		},
		{
			name: "InvalidWithdrawalAddressPrefix",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:            "0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "8c1Ff978036F2e9d7CC382Eff7B4c8c53C22acaa",
			},
			err: "invalid withdrawal address: withdrawal address 8c1Ff978036F2e9d7CC382Eff7B4c8c53C22acaa does not contain a 0x prefix",
		},
		{
			name: "InvalidWithdrawalAddressNotHex",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:            "0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0xrc1Ff978036F2e9d7CC382Eff7B4c8c53C22acaa",
			},
			err: "invalid withdrawal address: failed to obtain execution address: encoding/hex: invalid byte: U+0072 'r'",
		},
		{
			name: "UnknownValidatorPubKey",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:            "0xb384f767d964e100c8a9b21018d08c25ffebae268b3ab6d610353897541971726dbfc3c7463884c68a531515aab94c80",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},

			err: "unknown validator",
		},
		{
			name: "UnknownValidatorIndex",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:            "10",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},

			err: "unknown validator",
		},
		{
			name: "InvalidPubkeyLength",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:            "0xb384f767d964e100c8a9b21018d08c25ffebae268b3ab6d610353897541971726dbfc3c7463884c68a531515aab94c",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},

			err: "invalid public key: incorrect length",
		},
		{
			name: "Good",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:            "0",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			expected: []*capella.SignedBLSToExecutionChange{
				{
					Message: &capella.BLSToExecutionChange{
						ValidatorIndex:     0,
						FromBLSPubkey:      phase0.BLSPubKey{0x99, 0xb1, 0xf1, 0xd8, 0x4d, 0x76, 0x18, 0x54, 0x66, 0xd8, 0x6c, 0x34, 0xbd, 0xe1, 0x10, 0x13, 0x16, 0xaf, 0xdd, 0xae, 0x76, 0x21, 0x7a, 0xa8, 0x6c, 0xd0, 0x66, 0x97, 0x9b, 0x19, 0x85, 0x8c, 0x2c, 0x9d, 0x9e, 0x56, 0xee, 0xbc, 0x1e, 0x06, 0x7a, 0xc5, 0x42, 0x77, 0xa6, 0x17, 0x90, 0xdb},
						ToExecutionAddress: bellatrix.ExecutionAddress{0x8c, 0x1f, 0xf9, 0x78, 0x03, 0x6f, 0x2e, 0x9d, 0x7c, 0xc3, 0x82, 0xef, 0xf7, 0xb4, 0xc8, 0xc5, 0x3c, 0x22, 0xac, 0x15},
					},
					Signature: phase0.BLSSignature{0xb7, 0x8a, 0x05, 0xba, 0xd9, 0x27, 0xfc, 0x89, 0x6f, 0x14, 0x06, 0xb3, 0x2d, 0x64, 0x4a, 0xe1, 0x69, 0xce, 0xcd, 0x89, 0x86, 0xc1, 0xef, 0x8c, 0x0d, 0x03, 0x7d, 0x70, 0x86, 0xf8, 0x5f, 0x13, 0xe1, 0xe1, 0x88, 0xb4, 0x30, 0x96, 0x43, 0xa2, 0xc1, 0x3f, 0xfe, 0xfb, 0x0a, 0xe8, 0x05, 0x11, 0x09, 0x98, 0x53, 0xa0, 0x58, 0x1f, 0x4b, 0x2b, 0xd2, 0xe1, 0x45, 0x41, 0x04, 0x79, 0x01, 0xe2, 0x2a, 0x94, 0x0a, 0x9c, 0x7e, 0x3a, 0xc0, 0xa8, 0x82, 0xd1, 0xa8, 0xaf, 0x6b, 0xfa, 0xea, 0x81, 0x3a, 0x6a, 0x6b, 0xe7, 0x21, 0xf9, 0x26, 0x22, 0x04, 0xaa, 0x9d, 0xa4, 0xe4, 0x77, 0x27, 0xd0},
				},
			},
		},
		{
			name: "GoodPubkey",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				validator:            "0xb384f767d964e100c8a9b21018d08c25ffebae268b3ab6d610353897541971726dbfc3c7463884c68a531515aab94c87",
				chainInfo:            chainInfo,
				signedOperations:     make([]*capella.SignedBLSToExecutionChange, 0),
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			expected: []*capella.SignedBLSToExecutionChange{
				{
					Message: &capella.BLSToExecutionChange{
						ValidatorIndex:     0,
						FromBLSPubkey:      phase0.BLSPubKey{0x99, 0xb1, 0xf1, 0xd8, 0x4d, 0x76, 0x18, 0x54, 0x66, 0xd8, 0x6c, 0x34, 0xbd, 0xe1, 0x10, 0x13, 0x16, 0xaf, 0xdd, 0xae, 0x76, 0x21, 0x7a, 0xa8, 0x6c, 0xd0, 0x66, 0x97, 0x9b, 0x19, 0x85, 0x8c, 0x2c, 0x9d, 0x9e, 0x56, 0xee, 0xbc, 0x1e, 0x06, 0x7a, 0xc5, 0x42, 0x77, 0xa6, 0x17, 0x90, 0xdb},
						ToExecutionAddress: bellatrix.ExecutionAddress{0x8c, 0x1f, 0xf9, 0x78, 0x03, 0x6f, 0x2e, 0x9d, 0x7c, 0xc3, 0x82, 0xef, 0xf7, 0xb4, 0xc8, 0xc5, 0x3c, 0x22, 0xac, 0x15},
					},
					Signature: phase0.BLSSignature{0xb7, 0x8a, 0x05, 0xba, 0xd9, 0x27, 0xfc, 0x89, 0x6f, 0x14, 0x06, 0xb3, 0x2d, 0x64, 0x4a, 0xe1, 0x69, 0xce, 0xcd, 0x89, 0x86, 0xc1, 0xef, 0x8c, 0x0d, 0x03, 0x7d, 0x70, 0x86, 0xf8, 0x5f, 0x13, 0xe1, 0xe1, 0x88, 0xb4, 0x30, 0x96, 0x43, 0xa2, 0xc1, 0x3f, 0xfe, 0xfb, 0x0a, 0xe8, 0x05, 0x11, 0x09, 0x98, 0x53, 0xa0, 0x58, 0x1f, 0x4b, 0x2b, 0xd2, 0xe1, 0x45, 0x41, 0x04, 0x79, 0x01, 0xe2, 0x2a, 0x94, 0x0a, 0x9c, 0x7e, 0x3a, 0xc0, 0xa8, 0x82, 0xd1, 0xa8, 0xaf, 0x6b, 0xfa, 0xea, 0x81, 0x3a, 0x6a, 0x6b, 0xe7, 0x21, 0xf9, 0x26, 0x22, 0x04, 0xaa, 0x9d, 0xa4, 0xe4, 0x77, 0x27, 0xd0},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.command.generateOperationFromMnemonicAndValidator(ctx)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, test.command.signedOperations)
			}
		})
	}
}

func TestGenerateOperationFromSeedAndPath(t *testing.T) {
	ctx := context.Background()

	require.NoError(t, e2types.InitBLS())

	chainInfo := &beacon.ChainInfo{
		Version: 1,
		Validators: []*beacon.ValidatorInfo{
			{
				Index:                 0,
				Pubkey:                phase0.BLSPubKey{0xb3, 0x84, 0xf7, 0x67, 0xd9, 0x64, 0xe1, 0x00, 0xc8, 0xa9, 0xb2, 0x10, 0x18, 0xd0, 0x8c, 0x25, 0xff, 0xeb, 0xae, 0x26, 0x8b, 0x3a, 0xb6, 0xd6, 0x10, 0x35, 0x38, 0x97, 0x54, 0x19, 0x71, 0x72, 0x6d, 0xbf, 0xc3, 0xc7, 0x46, 0x38, 0x84, 0xc6, 0x8a, 0x53, 0x15, 0x15, 0xaa, 0xb9, 0x4c, 0x87},
				WithdrawalCredentials: []byte{0x00, 0x8b, 0xa1, 0xcc, 0x4b, 0x09, 0x1b, 0x91, 0xc1, 0x20, 0x2b, 0xba, 0x3f, 0x50, 0x80, 0x75, 0xd6, 0xff, 0x56, 0x5c, 0x77, 0xe5, 0x59, 0xf0, 0x80, 0x3c, 0x07, 0x92, 0xe0, 0x30, 0x2b, 0xf1},
			},
			{
				Index:                 1,
				Pubkey:                phase0.BLSPubKey{0xb3, 0xd8, 0x9e, 0x2f, 0x29, 0xc7, 0x12, 0xc6, 0xa9, 0xf8, 0xe5, 0xa2, 0x69, 0xb9, 0x76, 0x17, 0xc4, 0xa9, 0x4d, 0xd6, 0xf6, 0x66, 0x2a, 0xb3, 0xb0, 0x7c, 0xe9, 0xe5, 0x43, 0x45, 0x73, 0xf1, 0x5b, 0x5c, 0x98, 0x8c, 0xd1, 0x4b, 0xbd, 0x58, 0x04, 0xf7, 0x71, 0x56, 0xa8, 0xaf, 0x1c, 0xfa},
				WithdrawalCredentials: []byte{0x00, 0x78, 0x6c, 0xb0, 0x2e, 0xd2, 0x8e, 0x5f, 0xbb, 0x1f, 0x7f, 0x9e, 0x93, 0x1a, 0x2b, 0x72, 0x69, 0x29, 0x06, 0xe6, 0xb1, 0x2c, 0xe4, 0x64, 0x39, 0x75, 0xe3, 0x2b, 0x51, 0x76, 0x91, 0xf2},
			},
			{
				Index:                 2,
				Pubkey:                phase0.BLSPubKey{0xaf, 0x9c, 0xe4, 0x4f, 0x50, 0x14, 0x8d, 0xb4, 0x12, 0x19, 0x4a, 0xf0, 0xba, 0xf0, 0xba, 0xb3, 0x6b, 0xd5, 0xc3, 0xe0, 0xc4, 0x93, 0x89, 0x11, 0xa4, 0xe5, 0x02, 0xe3, 0x98, 0xb5, 0x9e, 0x5c, 0xca, 0x7c, 0x78, 0xe3, 0xfe, 0x03, 0x41, 0x95, 0x47, 0x88, 0x79, 0xee, 0xb2, 0x3d, 0xb0, 0xa6},
				WithdrawalCredentials: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x93, 0x1a, 0x2b, 0x72, 0x69, 0x29, 0x06, 0xe6, 0xb1, 0x2c, 0xe4, 0x64, 0x39, 0x75, 0xe3, 0x2b, 0x51, 0x76, 0x91, 0xf2},
			},
			{
				Index:                 3,
				Pubkey:                phase0.BLSPubKey{0x86, 0xd3, 0x30, 0xaf, 0x51, 0xfa, 0x59, 0x3f, 0xa9, 0xf9, 0x3e, 0xdb, 0x9d, 0x16, 0x64, 0x01, 0x86, 0xbe, 0x2e, 0x93, 0xea, 0x94, 0xd2, 0x59, 0x78, 0x1e, 0x1e, 0xb3, 0x4d, 0xeb, 0x84, 0x4c, 0x39, 0x68, 0xd7, 0x5e, 0xa9, 0x1d, 0x19, 0xf1, 0x59, 0xdb, 0xd0, 0x52, 0x3c, 0x6c, 0x5b, 0xa5},
				WithdrawalCredentials: []byte{0x00, 0x81, 0x68, 0x45, 0x6b, 0x6d, 0x9a, 0x32, 0x83, 0x93, 0x1f, 0xea, 0x52, 0x10, 0xda, 0x12, 0x2d, 0x1e, 0x65, 0xe8, 0xed, 0x50, 0xb8, 0xe8, 0xf5, 0x91, 0x11, 0x83, 0xb0, 0x2f, 0xd1, 0x25},
			},
		},
		GenesisValidatorsRoot: phase0.Root{},
		Epoch:                 1,
		CurrentForkVersion:    phase0.Version{},
	}
	validators := make(map[string]*beacon.ValidatorInfo, len(chainInfo.Validators))
	for i := range chainInfo.Validators {
		validators[fmt.Sprintf("%#x", chainInfo.Validators[i].Pubkey)] = chainInfo.Validators[i]
	}

	tests := []struct {
		name      string
		command   *command
		seed      []byte
		path      string
		generated bool
		err       string
		expected  []*capella.SignedBLSToExecutionChange
	}{
		{
			name: "PathInvalid",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				chainInfo:            chainInfo,
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			seed: []byte{0x40, 0x8b, 0x28, 0x5c, 0x12, 0x38, 0x36, 0x00, 0x4f, 0x4b, 0x88, 0x42, 0xc8, 0x93, 0x24, 0xc1, 0xf0, 0x13, 0x82, 0x45, 0x0c, 0x0d, 0x43, 0x9a, 0xf3, 0x45, 0xba, 0x7f, 0xc4, 0x9a, 0xcf, 0x70, 0x54, 0x89, 0xc6, 0xfc, 0x77, 0xdb, 0xd4, 0xe3, 0xdc, 0x1d, 0xd8, 0xcc, 0x6b, 0xc9, 0xf0, 0x43, 0xdb, 0x8a, 0xda, 0x1e, 0x24, 0x3c, 0x4a, 0x0e, 0xaf, 0xb2, 0x90, 0xd3, 0x99, 0x48, 0x08, 0x40},
			path: "invalid",
			err:  "failed to generate validator private key: not master at path component 0",
		},
		{
			name: "ValidatorUnknown",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				chainInfo:            chainInfo,
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			seed: []byte{0x40, 0x8b, 0x28, 0x5c, 0x12, 0x38, 0x36, 0x00, 0x4f, 0x4b, 0x88, 0x42, 0xc8, 0x93, 0x24, 0xc1, 0xf0, 0x13, 0x82, 0x45, 0x0c, 0x0d, 0x43, 0x9a, 0xf3, 0x45, 0xba, 0x7f, 0xc4, 0x9a, 0xcf, 0x70, 0x54, 0x89, 0xc6, 0xfc, 0x77, 0xdb, 0xd4, 0xe3, 0xdc, 0x1d, 0xd8, 0xcc, 0x6b, 0xc9, 0xf0, 0x43, 0xdb, 0x8a, 0xda, 0x1e, 0x24, 0x3c, 0x4a, 0x0e, 0xaf, 0xb2, 0x90, 0xd3, 0x99, 0x48, 0x08, 0x40},
			path: "m/12381/3600/999/0/0",
		},
		{
			name: "ValidatorCredentialsAlreadySet",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				chainInfo:            chainInfo,
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			seed: []byte{0x40, 0x8b, 0x28, 0x5c, 0x12, 0x38, 0x36, 0x00, 0x4f, 0x4b, 0x88, 0x42, 0xc8, 0x93, 0x24, 0xc1, 0xf0, 0x13, 0x82, 0x45, 0x0c, 0x0d, 0x43, 0x9a, 0xf3, 0x45, 0xba, 0x7f, 0xc4, 0x9a, 0xcf, 0x70, 0x54, 0x89, 0xc6, 0xfc, 0x77, 0xdb, 0xd4, 0xe3, 0xdc, 0x1d, 0xd8, 0xcc, 0x6b, 0xc9, 0xf0, 0x43, 0xdb, 0x8a, 0xda, 0x1e, 0x24, 0x3c, 0x4a, 0x0e, 0xaf, 0xb2, 0x90, 0xd3, 0x99, 0x48, 0x08, 0x40},
			path: "m/12381/3600/2/0/0",
		},
		{
			name: "PrivateKeyInvalid",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				chainInfo:            chainInfo,
				privateKey:           "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			seed: []byte{0x40, 0x8b, 0x28, 0x5c, 0x12, 0x38, 0x36, 0x00, 0x4f, 0x4b, 0x88, 0x42, 0xc8, 0x93, 0x24, 0xc1, 0xf0, 0x13, 0x82, 0x45, 0x0c, 0x0d, 0x43, 0x9a, 0xf3, 0x45, 0xba, 0x7f, 0xc4, 0x9a, 0xcf, 0x70, 0x54, 0x89, 0xc6, 0xfc, 0x77, 0xdb, 0xd4, 0xe3, 0xdc, 0x1d, 0xd8, 0xcc, 0x6b, 0xc9, 0xf0, 0x43, 0xdb, 0x8a, 0xda, 0x1e, 0x24, 0x3c, 0x4a, 0x0e, 0xaf, 0xb2, 0x90, 0xd3, 0x99, 0x48, 0x08, 0x40},
			path: "m/12381/3600/0/0/0",
			err:  "failed to create account from private key: invalid private key: err blsSecretKeyDeserialize ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name: "Good",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				chainInfo:            chainInfo,
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
			},
			seed:      []byte{0x40, 0x8b, 0x28, 0x5c, 0x12, 0x38, 0x36, 0x00, 0x4f, 0x4b, 0x88, 0x42, 0xc8, 0x93, 0x24, 0xc1, 0xf0, 0x13, 0x82, 0x45, 0x0c, 0x0d, 0x43, 0x9a, 0xf3, 0x45, 0xba, 0x7f, 0xc4, 0x9a, 0xcf, 0x70, 0x54, 0x89, 0xc6, 0xfc, 0x77, 0xdb, 0xd4, 0xe3, 0xdc, 0x1d, 0xd8, 0xcc, 0x6b, 0xc9, 0xf0, 0x43, 0xdb, 0x8a, 0xda, 0x1e, 0x24, 0x3c, 0x4a, 0x0e, 0xaf, 0xb2, 0x90, 0xd3, 0x99, 0x48, 0x08, 0x40},
			path:      "m/12381/3600/0/0/0",
			generated: true,
			expected: []*capella.SignedBLSToExecutionChange{
				{
					Message: &capella.BLSToExecutionChange{
						ValidatorIndex:     0,
						FromBLSPubkey:      phase0.BLSPubKey{0x99, 0xb1, 0xf1, 0xd8, 0x4d, 0x76, 0x18, 0x54, 0x66, 0xd8, 0x6c, 0x34, 0xbd, 0xe1, 0x10, 0x13, 0x16, 0xaf, 0xdd, 0xae, 0x76, 0x21, 0x7a, 0xa8, 0x6c, 0xd0, 0x66, 0x97, 0x9b, 0x19, 0x85, 0x8c, 0x2c, 0x9d, 0x9e, 0x56, 0xee, 0xbc, 0x1e, 0x06, 0x7a, 0xc5, 0x42, 0x77, 0xa6, 0x17, 0x90, 0xdb},
						ToExecutionAddress: bellatrix.ExecutionAddress{0x8c, 0x1f, 0xf9, 0x78, 0x03, 0x6f, 0x2e, 0x9d, 0x7c, 0xc3, 0x82, 0xef, 0xf7, 0xb4, 0xc8, 0xc5, 0x3c, 0x22, 0xac, 0x15},
					},
					Signature: phase0.BLSSignature{0xb7, 0x8a, 0x05, 0xba, 0xd9, 0x27, 0xfc, 0x89, 0x6f, 0x14, 0x06, 0xb3, 0x2d, 0x64, 0x4a, 0xe1, 0x69, 0xce, 0xcd, 0x89, 0x86, 0xc1, 0xef, 0x8c, 0x0d, 0x03, 0x7d, 0x70, 0x86, 0xf8, 0x5f, 0x13, 0xe1, 0xe1, 0x88, 0xb4, 0x30, 0x96, 0x43, 0xa2, 0xc1, 0x3f, 0xfe, 0xfb, 0x0a, 0xe8, 0x05, 0x11, 0x09, 0x98, 0x53, 0xa0, 0x58, 0x1f, 0x4b, 0x2b, 0xd2, 0xe1, 0x45, 0x41, 0x04, 0x79, 0x01, 0xe2, 0x2a, 0x94, 0x0a, 0x9c, 0x7e, 0x3a, 0xc0, 0xa8, 0x82, 0xd1, 0xa8, 0xaf, 0x6b, 0xfa, 0xea, 0x81, 0x3a, 0x6a, 0x6b, 0xe7, 0x21, 0xf9, 0x26, 0x22, 0x04, 0xaa, 0x9d, 0xa4, 0xe4, 0x77, 0x27, 0xd0},
				},
			},
		},
		{
			name: "GoodPrivateKey",
			command: &command{
				mnemonic:             "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
				chainInfo:            chainInfo,
				withdrawalAddressStr: "0x8c1Ff978036F2e9d7CC382Eff7B4c8c53C22ac15",
				privateKey:           "0x67775f030068b4610d6e1bd04948f547305b2502423fcece4c1091d065b44638",
			},
			seed:      []byte{0x40, 0x8b, 0x28, 0x5c, 0x12, 0x38, 0x36, 0x00, 0x4f, 0x4b, 0x88, 0x42, 0xc8, 0x93, 0x24, 0xc1, 0xf0, 0x13, 0x82, 0x45, 0x0c, 0x0d, 0x43, 0x9a, 0xf3, 0x45, 0xba, 0x7f, 0xc4, 0x9a, 0xcf, 0x70, 0x54, 0x89, 0xc6, 0xfc, 0x77, 0xdb, 0xd4, 0xe3, 0xdc, 0x1d, 0xd8, 0xcc, 0x6b, 0xc9, 0xf0, 0x43, 0xdb, 0x8a, 0xda, 0x1e, 0x24, 0x3c, 0x4a, 0x0e, 0xaf, 0xb2, 0x90, 0xd3, 0x99, 0x48, 0x08, 0x40},
			path:      "m/12381/3600/3/0/0",
			generated: true,
			expected: []*capella.SignedBLSToExecutionChange{
				{
					Message: &capella.BLSToExecutionChange{
						ValidatorIndex:     3,
						FromBLSPubkey:      phase0.BLSPubKey{0x86, 0x71, 0x0a, 0xbb, 0x44, 0xb6, 0xcd, 0xa6, 0x66, 0x57, 0x7b, 0xbb, 0x25, 0x5e, 0x16, 0xd9, 0x8b, 0xf2, 0x52, 0x51, 0x76, 0x22, 0x3f, 0x35, 0x35, 0xc7, 0xdf, 0xf8, 0xe7, 0x0b, 0x3b, 0xc8, 0x92, 0xbb, 0x36, 0x11, 0x33, 0x95, 0x2b, 0x03, 0xd2, 0xb0, 0x78, 0xcd, 0x07, 0x18, 0xca, 0xf3},
						ToExecutionAddress: bellatrix.ExecutionAddress{0x8c, 0x1f, 0xf9, 0x78, 0x03, 0x6f, 0x2e, 0x9d, 0x7c, 0xc3, 0x82, 0xef, 0xf7, 0xb4, 0xc8, 0xc5, 0x3c, 0x22, 0xac, 0x15},
					},
					Signature: phase0.BLSSignature{0x8d, 0x92, 0xb9, 0x1c, 0x5d, 0xfd, 0x98, 0xc7, 0x98, 0xfc, 0x94, 0xe1, 0xe6, 0x69, 0xf3, 0xaa, 0xae, 0x72, 0xb2, 0x36, 0x47, 0xde, 0x88, 0x54, 0xea, 0x16, 0x74, 0x7f, 0xfe, 0xf0, 0x4d, 0x46, 0x5c, 0x07, 0x56, 0x34, 0x03, 0x30, 0x2f, 0xbc, 0x26, 0xa2, 0x6d, 0xec, 0x10, 0x20, 0xe7, 0x67, 0x10, 0xb0, 0x4a, 0x7e, 0x4e, 0x25, 0x89, 0x7e, 0x87, 0x88, 0xda, 0xaf, 0x2b, 0xb5, 0xb7, 0x73, 0x25, 0x64, 0x80, 0xc1, 0xba, 0xf3, 0x1d, 0x33, 0x8f, 0x17, 0xa5, 0x35, 0x74, 0x80, 0xf3, 0x37, 0x0e, 0xea, 0x19, 0x15, 0xd5, 0x69, 0x7e, 0xf6, 0x68, 0xaa, 0x9c, 0x3d, 0x47, 0x19, 0x75, 0xfc},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			generated, err := test.command.generateOperationFromSeedAndPath(ctx, validators, test.seed, test.path)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.generated, generated)
				if generated {
					require.Equal(t, test.expected, test.command.signedOperations)
				}
			}
		})
	}
}
