/*
 * Copyright (C) 2019 The TesraSupernet Authors
 * This file is part of The TesraSupernet library.
 *
 * The TesraSupernet is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The TesraSupernet is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The TesraSupernet.  If not, see <http://www.gnu.org/licenses/>.
 */

package tesra_go_sdk

import (
	"encoding/hex"
	"fmt"
	"github.com/TesraSupernet/tesracrypto/signature"
	common2 "github.com/TesraSupernet/tesrasdk/common"
	"github.com/TesraSupernet/Tesra/common"
	"github.com/TesraSupernet/Tesra/core/payload"
	"github.com/TesraSupernet/Tesra/core/utils"
	"github.com/TesraSupernet/Tesra/core/validation"
	"github.com/TesraSupernet/Tesra/smartcontract/event"
	"github.com/TesraSupernet/Tesra/smartcontract/service/native/tsr"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var (
	testTesraSdk   *TesraSdk
	testWallet   *Wallet
	testPasswd   = []byte("123456")
	testDefAcc   *Account
	testGasPrice = uint64(500)
	testGasLimit = uint64(20000)
	testNetUrl   = "http://polaris1.tsr.io:20336"
)

func init() {
	var err error
	testWallet, err = testTesraSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("OpenWallet err: %s\n", err)
		return
	}
	testTesraSdk = NewTesraSdk()
	testTesraSdk.NewRpcClient().SetAddress(testNetUrl)
	testDefAcc, err = testWallet.GetDefaultAccount(testPasswd)
	if err != nil {
		fmt.Printf("GetDefaultAccount err: %s\n", err)
		return
	}
}
func TestTsrId_NewRegIDWithAttributesTransaction(t *testing.T) {
	testTesraSdk = NewTesraSdk()
}
func TestParseNativeTxPayload(t *testing.T) {
	testTesraSdk = NewTesraSdk()
	pri, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb7cf")
	assert.Nil(t, err)
	acc, err := NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	state := &tsr.State{
		From:  acc.Address,
		To:    acc.Address,
		Value: uint64(100),
	}
	transfers := make([]*tsr.State, 0)
	for i := 0; i < 1; i++ {
		transfers = append(transfers, state)
	}
	_, err = testTesraSdk.Native.Tsr.NewMultiTransferTransaction(500, 20000, transfers)
	assert.Nil(t, err)
	_, err = testTesraSdk.Native.Tsr.NewTransferFromTransaction(500, 20000, acc.Address, acc.Address, acc.Address, 20)
	assert.Nil(t, err)
}

func TestParsePayload(t *testing.T) {
	testTesraSdk = NewTesraSdk()
	//transferMulti
	payloadHex := "00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c0114c1087472616e736665721400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"
	//one transfer
	payloadHex = "00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0164c86c51c1087472616e736665721400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	//one transferFrom
	payloadHex = "00c66b6a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a14d2c124dd088190f709b684e0bc676d70c41b3776c86a0114c86c0c7472616e7366657246726f6d1400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	payloadBytes, err := common.HexToBytes(payloadHex)
	assert.Nil(t, err)
	_, err = ParsePayload(payloadBytes)
	assert.Nil(t, err)
}

func TestParsePayloadRandom(t *testing.T) {
	testTesraSdk = NewTesraSdk()
	pri, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb7cf")
	assert.Nil(t, err)
	acc, err := NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	assert.Nil(t, err)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000000; i++ {
		amount := rand.Intn(1000000)
		state := &tsr.State{
			From:  acc.Address,
			To:    acc.Address,
			Value: uint64(amount),
		}
		param := []*tsr.State{state}
		invokeCode, err := utils.BuildNativeInvokeCode(TSR_CONTRACT_ADDRESS, 0, "transfer", []interface{}{param})
		res, err := ParsePayload(invokeCode)
		assert.Nil(t, err)
		if res["param"] == nil {
			fmt.Println("amount:", amount)
			fmt.Println(res["param"])
			return
		} else {
			stateInfos := res["param"].([]common2.StateInfo)
			assert.Equal(t, uint64(amount), stateInfos[0].Value)
		}
		tr := tsr.TransferFrom{
			Sender: acc.Address,
			From:   acc.Address,
			To:     acc.Address,
			Value:  uint64(amount),
		}
		invokeCode, err = utils.BuildNativeInvokeCode(TSR_CONTRACT_ADDRESS, 0, "transferFrom", []interface{}{tr})
		res, err = ParsePayload(invokeCode)
		assert.Nil(t, err)
		if res["param"] == nil {
			fmt.Println("amount:", amount)
			fmt.Println(res["param"])
			return
		} else {
			stateInfos := res["param"].(common2.TransferFromInfo)
			assert.Equal(t, uint64(amount), stateInfos.Value)
		}
	}
}
func TestParsePayloadRandomMulti(t *testing.T) {
	testTesraSdk = NewTesraSdk()
	pri, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb7cf")
	assert.Nil(t, err)
	acc, err := NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	assert.Nil(t, err)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100000; i++ {
		amount := rand.Intn(10000000)
		state := &tsr.State{
			From:  acc.Address,
			To:    acc.Address,
			Value: uint64(amount),
		}
		paramLen := rand.Intn(20)
		if paramLen == 0 {
			paramLen += 1
		}
		params := make([]*tsr.State, 0)
		for i := 0; i < paramLen; i++ {
			params = append(params, state)
		}
		invokeCode, err := utils.BuildNativeInvokeCode(TSR_CONTRACT_ADDRESS, 0, "transfer", []interface{}{params})
		res, err := ParsePayload(invokeCode)
		assert.Nil(t, err)
		if res["param"] == nil {
			fmt.Println(res["param"])
			fmt.Println(amount)
			fmt.Println("invokeCode:", common.ToHexString(invokeCode))
			return
		} else {
			stateInfos := res["param"].([]common2.StateInfo)
			for i := 0; i < paramLen; i++ {
				assert.Equal(t, uint64(amount), stateInfos[i].Value)
			}
		}
	}
}

func TestTesraSdk_TrabsferFrom(t *testing.T) {
	testTesraSdk = NewTesraSdk()
	payloadHex := "00c66b1421ab6ece5c9e44fa5e35261ef42cc6bc31d98e9c6a7cc814c1d2d106f9d2276b383958973b9fca8e4f48cc966a7cc80400e1f5056a7cc86c51c1087472616e736665721400000000000000000000000000000000000000020068164f6e746f6c6f67792e4e61746976652e496e766f6b65"
	payloadBytes, err := common.HexToBytes(payloadHex)
	assert.Nil(t, err)
	res, err := ParsePayload(payloadBytes)
	assert.Nil(t, err)
	fmt.Println("res:", res)

	//java sdk,  transferFrom
	//amount =100
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc8149018fbdfe16d5b1054165ab892b0e040919bd1ca6a7cc8143e7c40c2a2a98e3f95adace19b12ef4a1d7a35066a7cc801646a7cc86c0c7472616e7366657246726f6d1400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"
	//amount =10
	//payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc8149018fbdfe16d5b1054165ab892b0e040919bd1ca6a7cc8143e7c40c2a2a98e3f95adace19b12ef4a1d7a35066a7cc85a6a7cc86c0c7472616e7366657246726f6d1400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	//amount = 1000000000
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc8149018fbdfe16d5b1054165ab892b0e040919bd1ca6a7cc8143e7c40c2a2a98e3f95adace19b12ef4a1d7a35066a7cc80400ca9a3b6a7cc86c0c7472616e7366657246726f6d1400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	//java sdk, transfer
	//amount = 100
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc814d2c124dd088190f709b684e0bc676d70c41b37766a7cc801646a7cc86c51c1087472616e736665721400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	//amount = 10
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc814d2c124dd088190f709b684e0bc676d70c41b37766a7cc85a6a7cc86c51c1087472616e736665721400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"
	//amount = 1000000000
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc814d2c124dd088190f709b684e0bc676d70c41b37766a7cc80400ca9a3b6a7cc86c51c1087472616e736665721400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	payloadBytes, err = common.HexToBytes(payloadHex)
	assert.Nil(t, err)
	res, err = ParsePayload(payloadBytes)
	assert.Nil(t, err)
	fmt.Println("res:", res)
}

//transferFrom
func TestTesraSdk_ParseNativeTxPayload2(t *testing.T) {
	var err error
	assert.Nil(t, err)
	pri, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb7cf")
	acc, err := NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)

	pri2, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb8cf")
	assert.Nil(t, err)

	pri3, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb9cf")
	assert.Nil(t, err)
	acc, err = NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)

	acc2, err := NewAccountFromPrivateKey(pri2, signature.SHA256withECDSA)

	acc3, err := NewAccountFromPrivateKey(pri3, signature.SHA256withECDSA)
	amount := 1000000000
	txFrom, err := testTesraSdk.Native.Tsr.NewTransferFromTransaction(500, 20000, acc.Address, acc2.Address, acc3.Address, uint64(amount))
	assert.Nil(t, err)
	tx, err := txFrom.IntoImmutable()
	assert.Nil(t, err)
	invokeCode, ok := tx.Payload.(*payload.InvokeCode)
	assert.True(t, ok)
	code := invokeCode.Code
	res, err := ParsePayload(code)
	assert.Nil(t, err)
	rp := res["param"].(common2.TransferFromInfo)
	assert.Equal(t, acc.Address.ToBase58(), rp.Sender)
	assert.Equal(t, acc2.Address.ToBase58(), rp.From)
	assert.Equal(t, uint64(amount), rp.Value)
	assert.Equal(t, "transferFrom", res["functionName"].(string))
	fmt.Println("res:", res)
}
func TestTesraSdk_ParseNativeTxPayload(t *testing.T) {
	testTesraSdk = NewTesraSdk()
	var err error
	assert.Nil(t, err)
	pri, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb7cf")
	acc, err := NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)

	pri2, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb8cf")
	assert.Nil(t, err)

	pri3, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb9cf")
	assert.Nil(t, err)
	acc, err = NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)

	acc2, err := NewAccountFromPrivateKey(pri2, signature.SHA256withECDSA)

	acc3, err := NewAccountFromPrivateKey(pri3, signature.SHA256withECDSA)
	y, _ := common.HexToBytes(acc.Address.ToHexString())

	fmt.Println("acc:", common.ToHexString(common.ToArrayReverse(y)))
	assert.Nil(t, err)

	amount := uint64(1000000000)
	tx, err := testTesraSdk.Native.Tsr.NewTransferTransaction(500, 20000, acc.Address, acc2.Address, amount)
	assert.Nil(t, err)

	tx2, err := tx.IntoImmutable()
	assert.Nil(t, err)
	res, err := ParseNativeTxPayload(tx2.ToArray())
	assert.Nil(t, err)
	fmt.Println("res:", res)
	states := res["param"].([]common2.StateInfo)
	assert.Equal(t, acc.Address.ToBase58(), states[0].From)
	assert.Equal(t, acc2.Address.ToBase58(), states[0].To)
	assert.Equal(t, amount, states[0].Value)
	assert.Equal(t, "transfer", res["functionName"].(string))

	transferFrom, err := testTesraSdk.Native.Tsr.NewTransferFromTransaction(500, 20000, acc.Address, acc2.Address, acc3.Address, 10)
	transferFrom2, err := transferFrom.IntoImmutable()
	r, err := ParseNativeTxPayload(transferFrom2.ToArray())
	assert.Nil(t, err)
	fmt.Println("res:", r)
	rp := r["param"].(common2.TransferFromInfo)
	assert.Equal(t, acc.Address.ToBase58(), rp.Sender)
	assert.Equal(t, acc2.Address.ToBase58(), rp.From)
	assert.Equal(t, uint64(10), rp.Value)

	tsgTransfer, err := testTesraSdk.Native.Tsg.NewTransferTransaction(uint64(500), uint64(20000), acc.Address, acc2.Address, 100000000)
	assert.Nil(t, err)
	tsgTx, err := tsgTransfer.IntoImmutable()
	assert.Nil(t, err)
	res, err = ParseNativeTxPayload(tsgTx.ToArray())
	assert.Nil(t, err)
	fmt.Println("res:", res)
}

func TestTesraSdk_GenerateMnemonicCodesStr2(t *testing.T) {
	mnemonic := make(map[string]bool)
	testTesraSdk := NewTesraSdk()
	for i := 0; i < 100000; i++ {
		mnemonicStr, err := testTesraSdk.GenerateMnemonicCodesStr()
		assert.Nil(t, err)
		if mnemonic[mnemonicStr] == true {
			panic("there is the same mnemonicStr ")
		} else {
			mnemonic[mnemonicStr] = true
		}
	}
}

func TestTesraSdk_GenerateMnemonicCodesStr(t *testing.T) {
	testTesraSdk := NewTesraSdk()
	for i := 0; i < 1000; i++ {
		mnemonic, err := testTesraSdk.GenerateMnemonicCodesStr()
		assert.Nil(t, err)
		private, err := testTesraSdk.GetPrivateKeyFromMnemonicCodesStrBip44(mnemonic, 0)
		assert.Nil(t, err)
		acc, err := NewAccountFromPrivateKey(private, signature.SHA256withECDSA)
		assert.Nil(t, err)
		si, err := signature.Sign(acc.SigScheme, acc.PrivateKey, []byte("test"), nil)
		boo := signature.Verify(acc.PublicKey, []byte("test"), si)
		assert.True(t, boo)

		tx, err := testTesraSdk.Native.Tsr.NewTransferTransaction(0, 0, acc.Address, acc.Address, 10)
		assert.Nil(t, err)
		testTesraSdk.SignToTransaction(tx, acc)
		tx2, err := tx.IntoImmutable()
		assert.Nil(t, err)
		res := validation.VerifyTransaction(tx2)
		assert.Equal(t, "not an error", res.Error())
	}
}

func TestGenerateMemory(t *testing.T) {
	expectedPrivateKey := []string{"915f5df65c75afe3293ed613970a1661b0b28d0cb711f21c489d8785977df0cd", "dbf1090889ba8b19aa01fa31c8b1ce29828bd2fa664afd95cc62e6055b74e112",
		"1487a8e53e4f4e2e1991781bcd14b3d334d3b2965cb48c976b234da29d7cf242", "79f85da015f079469c6e04aa0fc23523187d0f72c29450073d858ddeed272617"}
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	mnemonic = "ecology cricket napkin scrap board purpose picnic toe bean heart coast retire"
	testTesraSdk := NewTesraSdk()
	for i := 0; i < len(expectedPrivateKey); i++ {
		privk, err := testTesraSdk.GetPrivateKeyFromMnemonicCodesStrBip44(mnemonic, uint32(i))
		assert.Nil(t, err)
		assert.Equal(t, expectedPrivateKey[i], common.ToHexString(privk))
	}
}

func TestTesraSdk_CreateWallet(t *testing.T) {
	return
	wal, err := testTesraSdk.CreateWallet("./wallet2.dat")
	assert.Nil(t, err)
	if err != nil {
		return
	}
	_, err = wal.NewDefaultSettingAccount(testPasswd)
	assert.Nil(t, err)
	wal.Save()
}

func TestNewTesraSdk(t *testing.T) {
	testTesraSdk = NewTesraSdk()
	testWallet, _ = testTesraSdk.OpenWallet("./wallet.dat")
	event := &event.NotifyEventInfo{
		ContractAddress: common.ADDRESS_EMPTY,
		States:          []interface{}{"transfer", "Abc3UVbyL1kxd9sK6N9hzAT2u91ftbpoXT", "AFmseVrdL9f9oyCzZefL9tG6UbviEH9ugK", uint64(10000000)},
	}
	e, err := testTesraSdk.ParseNaitveTransferEvent(event)
	assert.Nil(t, err)
	fmt.Println(e)
}

func TestTesraSdk_GetTxData(t *testing.T) {
	testTesraSdk = NewTesraSdk()
	testWallet, _ = testTesraSdk.OpenWallet("./wallet.dat")
	acc, _ := testWallet.GetAccountByAddress("AXdmdzbyf3WZKQzRtrNQwAR91ZxMUfhXkt", testPasswd)
	tx, _ := testTesraSdk.Native.Tsr.NewTransferTransaction(500, 10000, acc.Address, acc.Address, 100)
	testTesraSdk.SignToTransaction(tx, acc)
	tx2, _ := tx.IntoImmutable()
	sink := common.NewZeroCopySink(nil)
	tx2.Serialization(sink)
	txData := hex.EncodeToString(sink.Bytes())
	tx3, _ := testTesraSdk.GetMutableTx(txData)
	assert.Equal(t, tx, tx3)
}

func Init() {
	testTesraSdk = NewTesraSdk()
	testTesraSdk.NewRpcClient().SetAddress(testNetUrl)

	var err error
	var wallet *Wallet
	if !common.FileExisted("./wallet.dat") {
		wallet, err = testTesraSdk.CreateWallet("./wallet.dat")
		if err != nil {
			fmt.Println("[CreateWallet] error:", err)
			return
		}
	} else {
		wallet, err = testTesraSdk.OpenWallet("./wallet.dat")
		if err != nil {
			fmt.Println("[CreateWallet] error:", err)
			return
		}
	}
	_, err = wallet.NewDefaultSettingAccount(testPasswd)
	if err != nil {
		fmt.Println("")
		return
	}
	wallet.Save()
	testWallet, err = testTesraSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("account.Open error:%s\n", err)
		return
	}
	testDefAcc, err = testWallet.GetDefaultAccount(testPasswd)
	if err != nil {
		fmt.Printf("GetDefaultAccount error:%s\n", err)
		return
	}

	return
	ws := testTesraSdk.NewWebSocketClient()
	err = ws.Connect("ws://localhost:20335")
	if err != nil {
		fmt.Printf("Connect ws error:%s", err)
		return
	}
}

func TestTesra_Transfer(t *testing.T) {
	return
	Init()
	testTesraSdk = NewTesraSdk()
	testTesraSdk.NewRpcClient().SetAddress(testNetUrl)
	testWallet, _ = testTesraSdk.OpenWallet("./wallet.dat")
	txHash, err := testTesraSdk.Native.Tsr.Transfer(testGasPrice, testGasLimit, nil, testDefAcc, testDefAcc.Address, 1)
	if err != nil {
		t.Errorf("NewTransferTransaction error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	evts, err := testTesraSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		t.Errorf("GetSmartContractEvent error:%s", err)
		return
	}
	fmt.Printf("TxHash:%s\n", txHash.ToHexString())
	fmt.Printf("State:%d\n", evts.State)
	fmt.Printf("GasConsume:%d\n", evts.GasConsumed)
	for _, notify := range evts.Notify {
		fmt.Printf("ContractAddress:%s\n", notify.ContractAddress)
		fmt.Printf("States:%+v\n", notify.States)
	}
}

func TestTsg_WithDrawTSG(t *testing.T) {
	Init()
	unboundTSG, err := testTesraSdk.Native.Tsg.UnboundTSG(testDefAcc.Address)
	if err != nil {
		t.Errorf("UnboundTSG error:%s", err)
		return
	}
	fmt.Printf("Address:%s UnboundTSG:%d\n", testDefAcc.Address.ToBase58(), unboundTSG)
	_, err = testTesraSdk.Native.Tsg.WithdrawTSG(500, 20000, nil, testDefAcc, unboundTSG)
	if err != nil {
		t.Errorf("WithDrawTSG error:%s", err)
		return
	}
	fmt.Printf("Address:%s WithDrawTSG amount:%d success\n", testDefAcc.Address.ToBase58(), unboundTSG)
}

func TestGlobalParam_GetGlobalParams(t *testing.T) {
	Init()
	gasPrice := "gasPrice"
	params := []string{gasPrice}
	results, err := testTesraSdk.Native.GlobalParams.GetGlobalParams(params)
	if err != nil {
		t.Errorf("GetGlobalParams:%+v error:%s", params, err)
		return
	}
	fmt.Printf("Params:%s Value:%v\n", gasPrice, results[gasPrice])
}

func TestGlobalParam_SetGlobalParams(t *testing.T) {
	return
	Init()
	gasPrice := "gasPrice"
	globalParams, err := testTesraSdk.Native.GlobalParams.GetGlobalParams([]string{gasPrice})
	if err != nil {
		t.Errorf("GetGlobalParams error:%s", err)
		return
	}
	gasPriceValue, err := strconv.Atoi(globalParams[gasPrice])
	if err != nil {
		t.Errorf("Get prama value error:%s", err)
		return
	}
	_, err = testTesraSdk.Native.GlobalParams.SetGlobalParams(testGasPrice, testGasLimit, nil, testDefAcc, map[string]string{gasPrice: strconv.Itoa(gasPriceValue + 1)})
	if err != nil {
		t.Errorf("SetGlobalParams error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	globalParams, err = testTesraSdk.Native.GlobalParams.GetGlobalParams([]string{gasPrice})
	if err != nil {
		t.Errorf("GetGlobalParams error:%s", err)
		return
	}
	gasPriceValueAfter, err := strconv.Atoi(globalParams[gasPrice])
	if err != nil {
		t.Errorf("Get prama value error:%s", err)
		return
	}
	fmt.Printf("After set params gasPrice:%d\n", gasPriceValueAfter)
}

func TestWsScribeEvent(t *testing.T) {
	return
	Init()
	wsClient := testTesraSdk.ClientMgr.GetWebSocketClient()
	err := wsClient.SubscribeEvent()
	if err != nil {
		t.Errorf("SubscribeTxHash error:%s", err)
		return
	}
	defer wsClient.UnsubscribeTxHash()

	actionCh := wsClient.GetActionCh()
	timer := time.NewTimer(time.Minute * 3)
	for {
		select {
		case <-timer.C:
			return
		case action := <-actionCh:
			fmt.Printf("Action:%s\n", action.Action)
			fmt.Printf("Result:%s\n", action.Result)
		}
	}
}

func TestWsTransfer(t *testing.T) {
	return
	Init()
	wsClient := testTesraSdk.ClientMgr.GetWebSocketClient()
	testTesraSdk.ClientMgr.SetDefaultClient(wsClient)
	txHash, err := testTesraSdk.Native.Tsr.Transfer(testGasPrice, testGasLimit, nil, testDefAcc, testDefAcc.Address, 1)
	if err != nil {
		t.Errorf("NewTransferTransaction error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	evts, err := testTesraSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		t.Errorf("GetSmartContractEvent error:%s", err)
		return
	}
	fmt.Printf("TxHash:%s\n", txHash.ToHexString())
	fmt.Printf("State:%d\n", evts.State)
	fmt.Printf("GasConsume:%d\n", evts.GasConsumed)
	for _, notify := range evts.Notify {
		fmt.Printf("ContractAddress:%s\n", notify.ContractAddress)
		fmt.Printf("States:%+v\n", notify.States)
	}
}
