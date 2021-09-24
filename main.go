// +build js
// +build wasm

package main

import (
	"git.xx.network/elixxir/mainnet-commitments/client"
	"gitlab.com/xx_network/primitives/utils"
	"syscall/js"
)

var done chan bool

func main() {
	f := js.FuncOf(SignAndTransmit)
	js.Global().Set("signAndTransmit", f)
	done = make(chan bool)
	<-done
}

// SignAndTransmit signs & transmits info to the commitments server
func SignAndTransmit(this js.Value, inputs []js.Value) interface{} {
	commitmentsCertPath := inputs[0].String()
	certPath := inputs[1].String()
	keyPath := inputs[2].String()
	idfPath := inputs[3].String()
	wallet := inputs[4].String()
	address := inputs[5].String()

	var cert, key, idf, commitmentCert []byte
	var err error
	var ep string
	// Read certificate file
	if ep, err = utils.ExpandPath(certPath); err == nil {
		cert, err = utils.ReadFile(ep)
		if err != nil {
			return map[string]interface{}{"Error": err.Error()}
		}
	} else {
		return map[string]interface{}{"Error": err.Error()}
	}

	// Read key file
	if ep, err = utils.ExpandPath(keyPath); err == nil {
		key, err = utils.ReadFile(ep)
		if err != nil {
			return map[string]interface{}{"Error": err.Error()}
		}
	} else {
		return map[string]interface{}{"Error": err.Error()}
	}

	// Read id file
	if ep, err = utils.ExpandPath(idfPath); err == nil {
		idf, err = utils.ReadFile(ep)
		if err != nil {
			return map[string]interface{}{"Error": err.Error()}
		}
	} else {
		return map[string]interface{}{"Error": err.Error()}
	}

	if ep, err = utils.ExpandPath(commitmentsCertPath); err == nil {
		commitmentCert, err = utils.ReadFile(ep)
		if err != nil {
			return map[string]interface{}{"Error": err.Error()}
		}
	} else {
		return map[string]interface{}{"Error": err.Error()}
	}

	// Sign & transmit information
	err = client.SignAndTransmit(key, cert, idf, commitmentCert, wallet, address)
	if err != nil {
		return map[string]interface{}{"Error": err.Error()}
	}
	return map[string]interface{}{}
}