/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gateway

import (
	"encoding/json"
)

const 	Hsmx509type = "HSM-X.509"

// Hsmx509Identity represents an Hsmx509 identity
type Hsmx509Identity struct {
	IDType      string      `json:"type"`
	Version     int         `json:"version"`
	MspID       string      `json:"mspId"`
	Credentials creds `json:"credentials"`
}

type creds struct {
	Certificate string `json:"certificate"`
}


// Type returns Hsmx509 for this identity type
func (hsmi *Hsmx509Identity) idType() string {
	return Hsmx509type
}

func (hsmi *Hsmx509Identity) mspID() string {
	return hsmi.MspID
}

// Certificate returns the Hsmx509 certificate PEM
func (hsmi *Hsmx509Identity) Certificate() string {
	return hsmi.Credentials.Certificate
}


// NewHsmx509Identity creates an Hsmx509 identity for storage in a wallet
func NewHsmx509Identity(mspid string, cert string) *Hsmx509Identity {
	return &Hsmx509Identity{Hsmx509type, 1, mspid, creds{cert }, }
}


func (hsmi *Hsmx509Identity) toJSON() ([]byte, error) {
	return json.Marshal(hsmi)
}

func (hsmi *Hsmx509Identity) fromJSON(data []byte) (Identity, error) {
	err := json.Unmarshal(data, hsmi)

	if err != nil {
		return nil, err
	}

	return hsmi, nil
}
