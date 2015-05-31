package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/tendermint/tendermint/account"
	rtypes "github.com/tendermint/tendermint/rpc/core/types"
	cclient "github.com/tendermint/tendermint/rpc/core_client"
	"github.com/tendermint/tendermint/types"
)

//------------------------------------------------------------------------------------
// core functions with string args, pertaining to a subcommand

func coreOutput(addr, amtS string) ([]byte, error) {
	if amtS == "" {
		return nil, fmt.Errorf("output must specify an amount with the --amt flag")
	}

	if addr == "" {
		return nil, fmt.Errorf("output must specify an addr with the --addr flag")
	}

	addrBytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, fmt.Errorf("addr is bad hex: %v", err)
	}

	amt, err := strconv.ParseUint(amtS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("amt is misformatted: %v", err)
	}

	txOutput := types.TxOutput{
		Address: addrBytes,
		Amount:  amt,
	}

	n, errPtr := new(int64), new(error)
	buf := new(bytes.Buffer)
	txOutput.WriteSignBytes(buf, n, errPtr)
	if *errPtr != nil {
		return nil, *errPtr
	}
	return buf.Bytes(), nil

}

func coreInput(pubkey, amtS, nonceS, addr string) ([]byte, error) {
	pub, addrBytes, amt, nonce, err := checkCommon(pubkey, addr, amtS, nonceS)
	if err != nil {
		return nil, err
	}

	txInput := types.TxInput{
		Address:  addrBytes,
		Amount:   amt,
		Sequence: uint(nonce),
		PubKey:   pub,
	}

	n, errPtr := new(int64), new(error)
	buf := new(bytes.Buffer)
	txInput.WriteSignBytes(buf, n, errPtr)
	if *errPtr != nil {
		return nil, *errPtr
	}
	return buf.Bytes(), nil
}

func coreSend(chainID, pubkey, addr, toAddr, amtS, nonceS string) (*types.SendTx, error) {
	pub, addrBytes, amt, nonce, err := checkCommon(pubkey, addr, amtS, nonceS)
	if err != nil {
		return nil, err
	}

	if toAddr == "" {
		return nil, fmt.Errorf("destination address must be given with --to flag")
	}

	toAddrBytes, err := hex.DecodeString(toAddr)
	if err != nil {
		return nil, fmt.Errorf("toAddr is bad hex: %v", err)
	}

	tx := types.NewSendTx()
	_ = addrBytes
	tx.AddInputWithNonce(pub, amt, uint(nonce))
	tx.AddOutput(toAddrBytes, amt)

	return tx, nil
}

func coreCall(chainID, pubkey, addr, toAddr, amtS, nonceS, gasS, feeS, data string) (*types.CallTx, error) {
	pub, _, amt, nonce, err := checkCommon(pubkey, addr, amtS, nonceS)
	if err != nil {
		return nil, err
	}

	toAddrBytes, err := hex.DecodeString(toAddr)
	if err != nil {
		return nil, fmt.Errorf("toAddr is bad hex: %v", err)
	}

	fee, err := strconv.ParseUint(feeS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("fee is misformatted: %v", err)
	}

	gas, err := strconv.ParseUint(gasS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("gas is misformatted: %v", err)
	}

	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("data is bad hex: %v", err)
	}

	tx := types.NewCallTxWithNonce(pub, toAddrBytes, dataBytes, amt, gas, fee, uint(nonce))
	return tx, nil
}

func coreName(chainID, pubkey, addr, amtS, nonceS, feeS, name, data string) (*types.NameTx, error) {
	pub, _, amt, nonce, err := checkCommon(pubkey, addr, amtS, nonceS)
	if err != nil {
		return nil, err
	}

	fee, err := strconv.ParseUint(feeS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("fee is misformatted: %v", err)
	}

	tx := types.NewNameTxWithNonce(pub, name, data, amt, fee, uint(nonce))
	return tx, nil
}

//------------------------------------------------------------------------------------

type HTTPResponse struct {
	Response string
	Error    string
}

func requestResponse(req *http.Request) (string, string, error) {
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode >= 400 {
		return "", "", fmt.Errorf(resp.Status)
	}
	return unpackResponse(resp)
}

func unpackResponse(resp *http.Response) (string, string, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	r := new(HTTPResponse)
	if err := json.Unmarshal(b, r); err != nil {
		return "", "", err
	}
	return r.Response, r.Error, nil
}

func coreSign(signBytes, signAddr, signRPC string) ([]byte, error) {
	req, _ := http.NewRequest("GET", signRPC+"/sign", nil)
	req.Header.Add("hash", signBytes)
	req.Header.Add("addr", signAddr)
	sig, errS, err := requestResponse(req)
	if err != nil {
		return nil, err
	}
	if errS != "" {
		return nil, fmt.Errorf("%s", errS)
	}
	sigBytes, err := hex.DecodeString(sig)
	return sigBytes, err
}

func coreBroadcast(tx types.Tx, broadcastRPC string) (*rtypes.Receipt, error) {
	client := cclient.NewClient(broadcastRPC, "JSONRPC")
	rb, err := client.BroadcastTx(tx)
	if err != nil {
		return nil, err
	}
	return &(rb.Receipt), nil
}

func checkCommon(pubkey, addr, amtS, nonceS string) (pub account.PubKey, addrBytes []byte, amt uint64, nonce uint64, err error) {
	if amtS == "" {
		err = fmt.Errorf("input must specify an amount with the --amt flag")
		return
	}

	if nonceS == "" {
		err = fmt.Errorf("input must specify a nonce with the --nonce flag")
		return
	}

	if pubkey == "" && addr == "" {
		err = fmt.Errorf("at least one of --pubkey or --addr must be given")
		return
	}

	fmt.Printf("PUBKEY %s\n", pubkey)
	pubKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		err = fmt.Errorf("pubkey is bad hex: %v", err)
		return
	}
	fmt.Printf("PUBKEY %X\n", pubKeyBytes)

	addrBytes, err = hex.DecodeString(addr)
	if err != nil {
		err = fmt.Errorf("addr is bad hex: %v", err)
		return
	}

	amt, err = strconv.ParseUint(amtS, 10, 64)
	if err != nil {
		err = fmt.Errorf("amt is misformatted: %v", err)
	}

	nonce, err = strconv.ParseUint(nonceS, 10, 64)
	if err != nil {
		err = fmt.Errorf("nonce is misformatted: %v", err)
		return
	}

	if len(pubKeyBytes) > 0 {
		pub = account.PubKeyEd25519(pubKeyBytes)
		addrBytes = pub.Address()
	}

	return
}
