///////////////////////////////////////////////////////////////////////////////
// Copyright © 2021 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package client

import (
	"crypto"
	"crypto/tls"
	"encoding/base64"
	"git.xx.network/elixxir/mainnet-commitments/messages"
	"git.xx.network/elixxir/mainnet-commitments/utils"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/signature/rsa"
	utils2 "gitlab.com/xx_network/primitives/utils"
)

func SignAndTransmit(keyPath, idfPath, nominatorWallet, validatorWallet, serverAddress, serverCert, contract, email string, selectedStake int) error {
	var key, idfBytes []byte
	var err error
	var ep string

	// Read key file
	if ep, err = utils2.ExpandPath(keyPath); err == nil {
		key, err = utils2.ReadFile(ep)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// Read id file
	if ep, err = utils2.ExpandPath(idfPath); err == nil {
		idfBytes, err = utils2.ReadFile(ep)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	h := crypto.BLAKE2b_512.New()
	_, err = h.Write([]byte(contract))
	if err != nil {
		return errors.WithMessage(err, "Failed to write contract to hash")
	}

	return signAndTransmit(key, idfBytes, h.Sum(nil), nominatorWallet, validatorWallet, serverAddress, serverCert, email, selectedStake)
}

// SignAndTransmit creates a Client object & transmits commitment info to the server
func signAndTransmit(pk, idfBytes, contractBytes []byte, nominatorWallet, validatorWallet, serverAddress, serverCert, email string, selectedStake int) error {
	// Create new resty client
	cl := resty.New()
	cl.SetRootCertificateFromString(serverCert)
	cl.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// Hash & sign node info
	key, err := rsa.LoadPrivateKeyFromPem(pk)
	if err != nil {
		return errors.WithMessage(err, "Failed to load private key")
	}
	hashed, hash, err := utils.HashNodeInfo(nominatorWallet, validatorWallet, idfBytes, contractBytes)
	if err != nil {
		return errors.WithMessage(err, "Failed to hash node info")
	}
	sig, err := rsa.Sign(csprng.NewSystemRNG(), key, hash, hashed, nil)
	if err != nil {
		return errors.WithMessage(err, "Failed to sign node info")
	}

	// Build message body & post to server
	body := messages.Commitment{
		IDF:             base64.URLEncoding.EncodeToString(idfBytes),
		Contract:        base64.URLEncoding.EncodeToString(contractBytes),
		NominatorWallet: nominatorWallet,
		ValidatorWallet: validatorWallet,
		Signature:       base64.URLEncoding.EncodeToString(sig),
		Email:           email,
		SelectedStake:   selectedStake,
	}
	resp, err := cl.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(messages.Commitment{}).
		Post(serverAddress + "/commitment")

	if err != nil {
		return errors.WithMessagef(err, "Failed to register commitment, received response: %+v", resp)
	} else if resp.IsError() {
		return errors.Errorf("Failed to register commitment, received response: %+v", resp)
	}

	return nil
}

func GetInfo(nid, serverCert, serverAddress string) ([]byte, error) {
	cl := resty.New()
	cl.SetRootCertificateFromString(serverCert)
	cl.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	resp, err := cl.R().SetQueryParam("id", nid).Get(serverAddress + "/info")
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to get commitment info, received response: %+v", resp)
	} else if resp.IsError() {
		return nil, errors.Errorf("Failed process request for commitment info, received response: %+v", resp)
	}

	return resp.Body(), nil
}
