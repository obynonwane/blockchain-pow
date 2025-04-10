package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Tx struct {
	FromID string `json:"from"`
	ToID   string `json:"to"`
	Value  uint64 `json:"value"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {

	privateKey, err := crypto.LoadECDSA("zblock/accounts/kennedy.ecdsa")
	if err != nil {
		return fmt.Errorf("unable to laod private key for node: %w", err)
	}

	tx := Tx{
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		ToID:   "Aaron",
		Value:  1000,
	}

	// get a slice of byte easily of tx
	data, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	// salt the signature such that I am able to determine which, chain It is meant for
	stamp := []byte(fmt.Sprintf("\x19Ardan Signed Message:\n%d", len(data)))

	// get a 32 bytes hash of the byte slice of data above and also concatenat the stamp
	// keccak256 is a variadic function
	v := crypto.Keccak256(stamp, data)

	sig, err := crypto.Sign(v, privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign: %w", err)
	}

	// printout the signature
	fmt.Println("SIG:", hexutil.Encode(sig))
	// ==================================================================================
	// OVER THE WIRE

	// returns the public key from the signature - using the ECDSA
	publicKey, err := crypto.SigToPub(v, sig)
	if err != nil {
		return fmt.Errorf("unable to pub: %w", err)
	}

	// return the ethereum common address := maybe by returning the first twenty bytes
	fmt.Println("PUB", crypto.PubkeyToAddress(*publicKey).String())

	// ==================================================================================

	tx = Tx{
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		ToID:   "Frank",
		Value:  250,
	}

	// get a slice of byte easily of tx
	data, err = json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	// get a 32 bytes hash of the byte slice of data above
	v2 := crypto.Keccak256(data)

	sig2, err := crypto.Sign(v2, privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign: %w", err)
	}

	// printout the signature
	fmt.Println("SIG:", hexutil.Encode(sig2))
	// ==================================================================================
	// OVER THE WIRE

	tx2 := Tx{
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		ToID:   "Frank",
		Value:  250,
	}

	// get a slice of byte easily of tx
	data, err = json.Marshal(tx2)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	// get a 32 bytes hash of the data slice
	v2 = crypto.Keccak256(data)

	// returns the public key from the signature - using the ECDSA
	publicKey, err = crypto.SigToPub(v2, sig2)
	if err != nil {
		return fmt.Errorf("unable to pub: %w", err)
	}

	// return the ethereum common address := maybe by returning the first twenty bytes
	fmt.Println("PUB", crypto.PubkeyToAddress(*publicKey).String())

	return nil
}

// Get the data
// get the byte slice
// get the 32byte hash of the byte slice (using keccak or sha256)
// then sign the 32 byte hash with your private key
// then forward it to the node for processing

// On the node side
// 1. You can check if the signature is properly structured - since you dont have the private key

//===================Structure of the signature================================
// there is nothing one can do with a signature other than to verify if it is structured properly
// It comes in three (3) parts
// 1. R value (second 32 byte is the second point on the elliptic curve(called secp256k1) )
// 2. S value (first 32 byte is fist point on the elliptic curve(called secp256k1) )
// 3. V value
// The ECDSA algorithm need only one of the R or S value to work, ethreum added the third Value V which is either 0 or 1
// so the V tells whether to use the first value of or the second value

// if the signature ends with zer (0), it is basically saying when using the signature use the first value
// if the signature ends with one (1) it is basically saying when using the signature use the second value
