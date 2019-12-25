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
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/TesraSupernet/tesracrypto/keypair"
	sdkcom "github.com/TesraSupernet/tesrasdk/common"
	"github.com/TesraSupernet/tesrasdk/utils"
	"github.com/TesraSupernet/Tesra/common"
	"github.com/TesraSupernet/Tesra/common/serialization"
	"github.com/TesraSupernet/Tesra/core/types"
	cutils "github.com/TesraSupernet/Tesra/core/utils"
	"github.com/TesraSupernet/Tesra/smartcontract/service/native/global_params"
	"github.com/TesraSupernet/Tesra/smartcontract/service/native/tesra"
)

var (
	TSR_CONTRACT_ADDRESS, _           = utils.AddressFromHexString("0100000000000000000000000000000000000000")
	TSG_CONTRACT_ADDRESS, _           = utils.AddressFromHexString("0200000000000000000000000000000000000000")
	TSR_ID_CONTRACT_ADDRESS, _        = utils.AddressFromHexString("0300000000000000000000000000000000000000")
	GLOABL_PARAMS_CONTRACT_ADDRESS, _ = utils.AddressFromHexString("0400000000000000000000000000000000000000")
	AUTH_CONTRACT_ADDRESS, _          = utils.AddressFromHexString("0600000000000000000000000000000000000000")
	GOVERNANCE_CONTRACT_ADDRESS, _    = utils.AddressFromHexString("0700000000000000000000000000000000000000")
)

var (
	TSR_CONTRACT_VERSION           = byte(0)
	TSG_CONTRACT_VERSION           = byte(0)
	TSR_ID_CONTRACT_VERSION        = byte(0)
	GLOBAL_PARAMS_CONTRACT_VERSION = byte(0)
	AUTH_CONTRACT_VERSION          = byte(0)
	GOVERNANCE_CONTRACT_VERSION    = byte(0)
)

var OPCODE_IN_PAYLOAD = map[byte]bool{0xc6: true, 0x6b: true, 0x6a: true, 0xc8: true, 0x6c: true, 0x68: true, 0x67: true,
	0x7c: true, 0xc1: true}

type NativeContract struct {
	tesraSdk       *TesraSdk
	Tesra          *Tesra
	Tsg          *Tsg
	TsrId        *TsrId
	GlobalParams *GlobalParam
	Auth         *Auth
}

func newNativeContract(tesraSdk *TesraSdk) *NativeContract {
	native := &NativeContract{tesraSdk: tesraSdk}
	Native.Tsr = &Tesra{native: native, tesraSdk: tesraSdk}
	native.Tsg = &Tsg{native: native, tesraSdk: tesraSdk}
	native.TsrId = &TsrId{native: native, tesraSdk: tesraSdk}
	native.GlobalParams = &GlobalParam{native: native, tesraSdk: tesraSdk}
	native.Auth = &Auth{native: native, tesraSdk: tesraSdk}
	return native
}

func (this *NativeContract) NewNativeInvokeTransaction(
	gasPrice,
	gasLimit uint64,
	version byte,
	contractAddress common.Address,
	method string,
	params []interface{},
) (*types.MutableTransaction, error) {
	if params == nil {
		params = make([]interface{}, 0, 1)
	}
	//Params cannot empty, if params is empty, fulfil with empty string
	if len(params) == 0 {
		params = append(params, "")
	}
	invokeCode, err := cutils.BuildNativeInvokeCode(contractAddress, version, method, params)
	if err != nil {
		return nil, fmt.Errorf("BuildNativeInvokeCode error:%s", err)
	}
	return this.tesraSdk.NewInvokeTransaction(gasPrice, gasLimit, invokeCode), nil
}

func (this *NativeContract) InvokeNativeContract(
	gasPrice,
	gasLimit uint64,
	payer,
	singer *Account,
	version byte,
	contractAddress common.Address,
	method string,
	params []interface{},
) (common.Uint256, error) {
	tx, err := this.NewNativeInvokeTransaction(gasPrice, gasLimit, version, contractAddress, method, params)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, singer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *NativeContract) PreExecInvokeNativeContract(
	contractAddress common.Address,
	version byte,
	method string,
	params []interface{},
) (*sdkcom.PreExecResult, error) {
	tx, err := this.NewNativeInvokeTransaction(0, 0, version, contractAddress, method, params)
	if err != nil {
		return nil, err
	}
	return this.tesraSdk.PreExecTransaction(tx)
}

type Tesra struct {
	tesraSdk *TesraSdk
	native *NativeContract
}

func (this *Tesra) NewTransferTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &tesra.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.NewMultiTransferTransaction(gasPrice, gasLimit, []*tesra.State{state})
}

func (this *Tesra) Transfer(gasPrice, gasLimit uint64, payer *Account, from *Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferTransaction(gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tesra) NewMultiTransferTransaction(gasPrice, gasLimit uint64, states []*tesra.State) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_CONTRACT_VERSION,
		TSR_CONTRACT_ADDRESS,
		tesra.TRANSFER_NAME,
		[]interface{}{states})
}

func (this *Tesra) MultiTransfer(gasPrice, gasLimit uint64, payer *Account, states []*tesra.State, signer *Account) (common.Uint256, error) {
	tx, err := this.NewMultiTransferTransaction(gasPrice, gasLimit, states)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tesra) NewTransferFromTransaction(gasPrice, gasLimit uint64, sender, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &tesra.TransferFrom{
		Sender: sender,
		From:   from,
		To:     to,
		Value:  amount,
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_CONTRACT_VERSION,
		TSR_CONTRACT_ADDRESS,
		tesra.TRANSFERFROM_NAME,
		[]interface{}{state},
	)
}

func (this *Tesra) TransferFrom(gasPrice, gasLimit uint64, payer *Account, sender *Account, from, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferFromTransaction(gasPrice, gasLimit, sender.Address, from, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, sender)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tesra) NewApproveTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &tesra.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_CONTRACT_VERSION,
		TSR_CONTRACT_ADDRESS,
		tesra.APPROVE_NAME,
		[]interface{}{state},
	)
}

func (this *Tesra) Approve(gasPrice, gasLimit uint64, payer *Account, from *Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewApproveTransaction(gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tesra) Allowance(from, to common.Address) (uint64, error) {
	type allowanceStruct struct {
		From common.Address
		To   common.Address
	}
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_CONTRACT_ADDRESS,
		TSR_CONTRACT_VERSION,
		tesra.ALLOWANCE_NAME,
		[]interface{}{&allowanceStruct{From: from, To: to}},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Tesra) Symbol() (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_CONTRACT_ADDRESS,
		TSR_CONTRACT_VERSION,
		tesra.SYMBOL_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Tesra) BalanceOf(address common.Address) (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_CONTRACT_ADDRESS,
		TSR_CONTRACT_VERSION,
		tesra.BALANCEOF_NAME,
		[]interface{}{address[:]},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Tesra) Name() (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_CONTRACT_ADDRESS,
		TSR_CONTRACT_VERSION,
		tesra.NAME_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Tesra) Decimals() (byte, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_CONTRACT_ADDRESS,
		TSR_CONTRACT_VERSION,
		tesra.DECIMALS_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	decimals, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return byte(decimals.Uint64()), nil
}

func (this *Tesra) TotalSupply() (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_CONTRACT_ADDRESS,
		TSR_CONTRACT_VERSION,
		tesra.TOTAL_SUPPLY_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

type Tsg struct {
	tesraSdk *TesraSdk
	native *NativeContract
}

func (this *Tsg) NewTransferTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &tesra.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.NewMultiTransferTransaction(gasPrice, gasLimit, []*tesra.State{state})
}

func (this *Tsg) Transfer(gasPrice, gasLimit uint64, payer *Account, from *Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferTransaction(gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tsg) NewMultiTransferTransaction(gasPrice, gasLimit uint64, states []*tesra.State) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSG_CONTRACT_VERSION,
		TSG_CONTRACT_ADDRESS,
		tesra.TRANSFER_NAME,
		[]interface{}{states})
}

func (this *Tsg) MultiTransfer(gasPrice, gasLimit uint64, states []*tesra.State, signer *Account) (common.Uint256, error) {
	tx, err := this.NewMultiTransferTransaction(gasPrice, gasLimit, states)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tsg) NewTransferFromTransaction(gasPrice, gasLimit uint64, sender, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &tesra.TransferFrom{
		Sender: sender,
		From:   from,
		To:     to,
		Value:  amount,
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSG_CONTRACT_VERSION,
		TSG_CONTRACT_ADDRESS,
		tesra.TRANSFERFROM_NAME,
		[]interface{}{state},
	)
}

func (this *Tsg) TransferFrom(gasPrice, gasLimit uint64, payer *Account, sender *Account, from, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferFromTransaction(gasPrice, gasLimit, sender.Address, from, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, sender)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tsg) NewWithdrawTSGTransaction(gasPrice, gasLimit uint64, address common.Address, amount uint64) (*types.MutableTransaction, error) {
	return this.NewTransferFromTransaction(gasPrice, gasLimit, address, TSR_CONTRACT_ADDRESS, address, amount)
}

func (this *Tsg) WithdrawTSG(gasPrice, gasLimit uint64, payer *Account, address *Account, amount uint64) (common.Uint256, error) {
	tx, err := this.NewWithdrawTSGTransaction(gasPrice, gasLimit, address.Address, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, address)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tsg) NewApproveTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &tesra.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSG_CONTRACT_VERSION,
		TSG_CONTRACT_ADDRESS,
		tesra.APPROVE_NAME,
		[]interface{}{state},
	)
}

func (this *Tsg) Approve(gasPrice, gasLimit uint64, payer *Account, from *Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewApproveTransaction(gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Tsg) Allowance(from, to common.Address) (uint64, error) {
	type allowanceStruct struct {
		From common.Address
		To   common.Address
	}
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSG_CONTRACT_ADDRESS,
		TSG_CONTRACT_VERSION,
		tesra.ALLOWANCE_NAME,
		[]interface{}{&allowanceStruct{From: from, To: to}},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Tsg) UnboundTSG(address common.Address) (uint64, error) {
	return this.Allowance(TSR_CONTRACT_ADDRESS, address)
}

func (this *Tsg) Symbol() (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSG_CONTRACT_ADDRESS,
		TSG_CONTRACT_VERSION,
		tesra.SYMBOL_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Tsg) BalanceOf(address common.Address) (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSG_CONTRACT_ADDRESS,
		TSG_CONTRACT_VERSION,
		tesra.BALANCEOF_NAME,
		[]interface{}{address[:]},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Tsg) Name() (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSG_CONTRACT_ADDRESS,
		TSG_CONTRACT_VERSION,
		tesra.NAME_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Tsg) Decimals() (byte, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSG_CONTRACT_ADDRESS,
		TSG_CONTRACT_VERSION,
		tesra.DECIMALS_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	decimals, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return byte(decimals.Uint64()), nil
}

func (this *Tsg) TotalSupply() (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSG_CONTRACT_ADDRESS,
		TSG_CONTRACT_VERSION,
		tesra.TOTAL_SUPPLY_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

type TsrId struct {
	tesraSdk *TesraSdk
	native *NativeContract
}

func (this *TsrId) NewRegIDWithPublicKeyTransaction(gasPrice, gasLimit uint64, tsrId string, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type regIDWithPublicKey struct {
		TsrId  string
		PubKey []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"regIDWithPublicKey",
		[]interface{}{
			&regIDWithPublicKey{
				TsrId:  tsrId,
				PubKey: keypair.SerializePublicKey(pubKey),
			},
		},
	)
}

func (this *TsrId) RegIDWithPublicKey(gasPrice, gasLimit uint64, payer *Account, signer *Account, tsrId string, controller *Controller) (common.Uint256, error) {
	tx, err := this.NewRegIDWithPublicKeyTransaction(gasPrice, gasLimit, tsrId, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *TsrId) NewRegIDWithAttributesTransaction(gasPrice, gasLimit uint64, tsrId string, pubKey keypair.PublicKey, attributes []*DDOAttribute) (*types.MutableTransaction, error) {
	type regIDWithAttribute struct {
		TsrId      string
		PubKey     []byte
		Attributes []*DDOAttribute
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"regIDWithAttributes",
		[]interface{}{
			&regIDWithAttribute{
				TsrId:      tsrId,
				PubKey:     keypair.SerializePublicKey(pubKey),
				Attributes: attributes,
			},
		},
	)
}

func (this *TsrId) RegIDWithAttributes(gasPrice, gasLimit uint64, payer, signer *Account, tsrId string, controller *Controller, attributes []*DDOAttribute) (common.Uint256, error) {
	tx, err := this.NewRegIDWithAttributesTransaction(gasPrice, gasLimit, tsrId, controller.PublicKey, attributes)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *TsrId) GetDDO(tsrId string) (*DDO, error) {
	result, err := this.native.PreExecInvokeNativeContract(
		TSR_ID_CONTRACT_ADDRESS,
		TSR_ID_CONTRACT_VERSION,
		"getDDO",
		[]interface{}{tsrId},
	)
	if err != nil {
		return nil, err
	}
	data, err := result.Result.ToByteArray()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	keyData, err := serialization.ReadVarBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("key ReadVarBytes error:%s", err)
	}
	owners, err := this.getPublicKeys(tsrId, keyData)
	if err != nil {
		return nil, fmt.Errorf("getPublicKeys error:%s", err)
	}
	attrData, err := serialization.ReadVarBytes(buf)
	attrs, err := this.getAttributes(tsrId, attrData)
	if err != nil {
		return nil, fmt.Errorf("getAttributes error:%s", err)
	}
	recoveryData, err := serialization.ReadVarBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("recovery ReadVarBytes error:%s", err)
	}
	var addr string
	if len(recoveryData) != 0 {
		address, err := common.AddressParseFromBytes(recoveryData)
		if err != nil {
			return nil, fmt.Errorf("AddressParseFromBytes error:%s", err)
		}
		addr = address.ToBase58()
	}

	ddo := &DDO{
		TsrId:      tsrId,
		Owners:     owners,
		Attributes: attrs,
		Recovery:   addr,
	}
	return ddo, nil
}

func (this *TsrId) NewAddKeyTransaction(gasPrice, gasLimit uint64, tsrId string, newPubKey, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type addKey struct {
		TsrId     string
		NewPubKey []byte
		PubKey    []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"addKey",
		[]interface{}{
			&addKey{
				TsrId:     tsrId,
				NewPubKey: keypair.SerializePublicKey(newPubKey),
				PubKey:    keypair.SerializePublicKey(pubKey),
			},
		})
}

func (this *TsrId) AddKey(gasPrice, gasLimit uint64, payer *Account, tsrId string, signer *Account, newPubKey keypair.PublicKey, controller *Controller) (common.Uint256, error) {
	tx, err := this.NewAddKeyTransaction(gasPrice, gasLimit, tsrId, newPubKey, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *TsrId) NewRevokeKeyTransaction(gasPrice, gasLimit uint64, tsrId string, removedPubKey, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type removeKey struct {
		TsrId      string
		RemovedKey []byte
		PubKey     []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"removeKey",
		[]interface{}{
			&removeKey{
				TsrId:      tsrId,
				RemovedKey: keypair.SerializePublicKey(removedPubKey),
				PubKey:     keypair.SerializePublicKey(pubKey),
			},
		},
	)
}

func (this *TsrId) RevokeKey(gasPrice, gasLimit uint64, payer *Account, tsrId string, signer *Account, removedPubKey keypair.PublicKey, controller *Controller) (common.Uint256, error) {
	tx, err := this.NewRevokeKeyTransaction(gasPrice, gasLimit, tsrId, removedPubKey, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *TsrId) NewSetRecoveryTransaction(gasPrice, gasLimit uint64, tsrId string, recovery common.Address, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type addRecovery struct {
		TsrId    string
		Recovery common.Address
		Pubkey   []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"addRecovery",
		[]interface{}{
			&addRecovery{
				TsrId:    tsrId,
				Recovery: recovery,
				Pubkey:   keypair.SerializePublicKey(pubKey),
			},
		})
}

func (this *TsrId) SetRecovery(gasPrice, gasLimit uint64, payer, signer *Account, tsrId string, recovery common.Address, controller *Controller) (common.Uint256, error) {
	tx, err := this.NewSetRecoveryTransaction(gasPrice, gasLimit, tsrId, recovery, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *TsrId) NewChangeRecoveryTransaction(gasPrice, gasLimit uint64, tsrId string, newRecovery, oldRecovery common.Address) (*types.MutableTransaction, error) {
	type changeRecovery struct {
		TsrId       string
		NewRecovery common.Address
		OldRecovery common.Address
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"changeRecovery",
		[]interface{}{
			&changeRecovery{
				TsrId:       tsrId,
				NewRecovery: newRecovery,
				OldRecovery: oldRecovery,
			},
		})
}

func (this *TsrId) ChangeRecovery(gasPrice, gasLimit uint64, payer, signer *Account, tsrId string, newRecovery, oldRecovery common.Address, controller *Controller) (common.Uint256, error) {
	tx, err := this.NewChangeRecoveryTransaction(gasPrice, gasLimit, tsrId, newRecovery, oldRecovery)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *TsrId) NewAddAttributesTransaction(gasPrice, gasLimit uint64, tsrId string, attributes []*DDOAttribute, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type addAttributes struct {
		TsrId      string
		Attributes []*DDOAttribute
		PubKey     []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"addAttributes",
		[]interface{}{
			&addAttributes{
				TsrId:      tsrId,
				Attributes: attributes,
				PubKey:     keypair.SerializePublicKey(pubKey),
			},
		})
}

func (this *TsrId) AddAttributes(gasPrice, gasLimit uint64, payer, signer *Account, tsrId string, attributes []*DDOAttribute, controller *Controller) (common.Uint256, error) {
	tx, err := this.NewAddAttributesTransaction(gasPrice, gasLimit, tsrId, attributes, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}

	return this.tesraSdk.SendTransaction(tx)
}

func (this *TsrId) NewRemoveAttributeTransaction(gasPrice, gasLimit uint64, tsrId string, key []byte, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type removeAttribute struct {
		TsrId  string
		Key    []byte
		PubKey []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"removeAttribute",
		[]interface{}{
			&removeAttribute{
				TsrId:  tsrId,
				Key:    key,
				PubKey: keypair.SerializePublicKey(pubKey),
			},
		})
}

func (this *TsrId) RemoveAttribute(gasPrice, gasLimit uint64, payer, signer *Account, tsrId string, removeKey []byte, controller *Controller) (common.Uint256, error) {
	tx, err := this.NewRemoveAttributeTransaction(gasPrice, gasLimit, tsrId, removeKey, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}

	return this.tesraSdk.SendTransaction(tx)
}

func (this *TsrId) GetAttributes(tsrId string) ([]*DDOAttribute, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_ID_CONTRACT_ADDRESS,
		TSR_ID_CONTRACT_VERSION,
		"getAttributes",
		[]interface{}{tsrId})
	if err != nil {
		return nil, err
	}
	data, err := preResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("ToByteArray error:%s", err)
	}
	return this.getAttributes(tsrId, data)
}

func (this *TsrId) getAttributes(tsrId string, data []byte) ([]*DDOAttribute, error) {
	buf := bytes.NewBuffer(data)
	attributes := make([]*DDOAttribute, 0)
	for {
		if buf.Len() == 0 {
			break
		}
		key, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("key ReadVarBytes error:%s", err)
		}
		valueType, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("value type ReadVarBytes error:%s", err)
		}
		value, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("value ReadVarBytes error:%s", err)
		}
		attributes = append(attributes, &DDOAttribute{
			Key:       key,
			Value:     value,
			ValueType: valueType,
		})
	}
	//reverse
	for i, j := 0, len(attributes)-1; i < j; i, j = i+1, j-1 {
		attributes[i], attributes[j] = attributes[j], attributes[i]
	}
	return attributes, nil
}

func (this *TsrId) VerifySignature(tsrId string, keyIndex int, controller *Controller) (bool, error) {
	tx, err := this.native.NewNativeInvokeTransaction(
		0, 0,
		TSR_ID_CONTRACT_VERSION,
		TSR_ID_CONTRACT_ADDRESS,
		"verifySignature",
		[]interface{}{tsrId, keyIndex})
	if err != nil {
		return false, err
	}
	err = this.tesraSdk.SignToTransaction(tx, controller)
	if err != nil {
		return false, err
	}
	preResult, err := this.tesraSdk.PreExecTransaction(tx)
	if err != nil {
		return false, err
	}
	return preResult.Result.ToBool()
}

func (this *TsrId) GetPublicKeys(tsrId string) ([]*DDOOwner, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_ID_CONTRACT_ADDRESS,
		TSR_ID_CONTRACT_VERSION,
		"getPublicKeys",
		[]interface{}{
			tsrId,
		})
	if err != nil {
		return nil, err
	}
	data, err := preResult.Result.ToByteArray()
	if err != nil {
		return nil, err
	}
	return this.getPublicKeys(tsrId, data)
}

func (this *TsrId) getPublicKeys(tsrId string, data []byte) ([]*DDOOwner, error) {
	buf := bytes.NewBuffer(data)
	owners := make([]*DDOOwner, 0)
	for {
		if buf.Len() == 0 {
			break
		}
		index, err := serialization.ReadUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("index ReadUint32 error:%s", err)
		}
		pubKeyId := fmt.Sprintf("%s#keys-%d", tsrId, index)
		pkData, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("PubKey Idenx:%d ReadVarBytes error:%s", index, err)
		}
		pubKey, err := keypair.DeserializePublicKey(pkData)
		if err != nil {
			return nil, fmt.Errorf("DeserializePublicKey Index:%d error:%s", index, err)
		}
		keyType := keypair.GetKeyType(pubKey)
		owner := &DDOOwner{
			pubKeyIndex: index,
			PubKeyId:    pubKeyId,
			Type:        GetKeyTypeString(keyType),
			Curve:       GetCurveName(pkData),
			Value:       hex.EncodeToString(pkData),
		}
		owners = append(owners, owner)
	}
	return owners, nil
}

func (this *TsrId) GetKeyState(tsrId string, keyIndex int) (string, error) {
	type keyState struct {
		TsrId    string
		KeyIndex int
	}
	preResult, err := this.native.PreExecInvokeNativeContract(
		TSR_ID_CONTRACT_ADDRESS,
		TSR_ID_CONTRACT_VERSION,
		"getKeyState",
		[]interface{}{
			&keyState{
				TsrId:    tsrId,
				KeyIndex: keyIndex,
			},
		})
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

type GlobalParam struct {
	tesraSdk *TesraSdk
	native *NativeContract
}

func (this *GlobalParam) GetGlobalParams(params []string) (map[string]string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		global_params.GET_GLOBAL_PARAM_NAME,
		[]interface{}{params})
	if err != nil {
		return nil, err
	}
	results, err := preResult.Result.ToByteArray()
	if err != nil {
		return nil, err
	}
	queryParams := new(global_params.Params)
	err = queryParams.Deserialization(common.NewZeroCopySource(results))
	if err != nil {
		return nil, err
	}
	globalParams := make(map[string]string, len(params))
	for _, param := range params {
		index, values := queryParams.GetParam(param)
		if index < 0 {
			continue
		}
		globalParams[param] = values.Value
	}
	return globalParams, nil
}

func (this *GlobalParam) NewSetGlobalParamsTransaction(gasPrice, gasLimit uint64, params map[string]string) (*types.MutableTransaction, error) {
	var globalParams global_params.Params
	for k, v := range params {
		globalParams.SetParam(global_params.Param{Key: k, Value: v})
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.SET_GLOBAL_PARAM_NAME,
		[]interface{}{globalParams})
}

func (this *GlobalParam) SetGlobalParams(gasPrice, gasLimit uint64, payer, signer *Account, params map[string]string) (common.Uint256, error) {
	tx, err := this.NewSetGlobalParamsTransaction(gasPrice, gasLimit, params)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *GlobalParam) NewTransferAdminTransaction(gasPrice, gasLimit uint64, newAdmin common.Address) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.TRANSFER_ADMIN_NAME,
		[]interface{}{newAdmin})
}

func (this *GlobalParam) TransferAdmin(gasPrice, gasLimit uint64, payer, signer *Account, newAdmin common.Address) (common.Uint256, error) {
	tx, err := this.NewTransferAdminTransaction(gasPrice, gasLimit, newAdmin)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *GlobalParam) NewAcceptAdminTransaction(gasPrice, gasLimit uint64, admin common.Address) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.ACCEPT_ADMIN_NAME,
		[]interface{}{admin})
}

func (this *GlobalParam) AcceptAdmin(gasPrice, gasLimit uint64, payer, signer *Account) (common.Uint256, error) {
	tx, err := this.NewAcceptAdminTransaction(gasPrice, gasLimit, signer.Address)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *GlobalParam) NewSetOperatorTransaction(gasPrice, gasLimit uint64, operator common.Address) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.SET_OPERATOR,
		[]interface{}{operator},
	)
}

func (this *GlobalParam) SetOperator(gasPrice, gasLimit uint64, payer, signer *Account, operator common.Address) (common.Uint256, error) {
	tx, err := this.NewSetOperatorTransaction(gasPrice, gasLimit, operator)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *GlobalParam) NewCreateSnapshotTransaction(gasPrice, gasLimit uint64) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.CREATE_SNAPSHOT_NAME,
		[]interface{}{},
	)
}

func (this *GlobalParam) CreateSnapshot(gasPrice, gasLimit uint64, payer, signer *Account) (common.Uint256, error) {
	tx, err := this.NewCreateSnapshotTransaction(gasPrice, gasLimit)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

type Auth struct {
	tesraSdk *TesraSdk
	native *NativeContract
}

func (this *Auth) NewAssignFuncsToRoleTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, adminId, role []byte, funcNames []string, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"assignFuncsToRole",
		[]interface{}{
			contractAddress,
			adminId,
			role,
			funcNames,
			keyIndex,
		})
}

func (this *Auth) AssignFuncsToRole(gasPrice, gasLimit uint64, contractAddress common.Address, payer, signer *Account, adminId, role []byte, funcNames []string, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewAssignFuncsToRoleTransaction(gasPrice, gasLimit, contractAddress, adminId, role, funcNames, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Auth) NewDelegateTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, from, to, role []byte, period, level, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"delegate",
		[]interface{}{
			contractAddress,
			from,
			to,
			role,
			period,
			level,
			keyIndex,
		})
}

func (this *Auth) Delegate(gasPrice, gasLimit uint64, payer, signer *Account, contractAddress common.Address, from, to, role []byte, period, level, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewDelegateTransaction(gasPrice, gasLimit, contractAddress, from, to, role, period, level, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Auth) NewWithdrawTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, initiator, delegate, role []byte, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"withdraw",
		[]interface{}{
			contractAddress,
			initiator,
			delegate,
			role,
			keyIndex,
		})
}

func (this *Auth) Withdraw(gasPrice, gasLimit uint64, payer, signer *Account, contractAddress common.Address, initiator, delegate, role []byte, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewWithdrawTransaction(gasPrice, gasLimit, contractAddress, initiator, delegate, role, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Auth) NewAssignTesraIDsToRoleTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, admtsrId, role []byte, persons [][]byte, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"assignTesraIDsToRole",
		[]interface{}{
			contractAddress,
			admtsrId,
			role,
			persons,
			keyIndex,
		})
}

func (this *Auth) AssignTesraIDsToRole(gasPrice, gasLimit uint64, payer, signer *Account, contractAddress common.Address, admtsrId, role []byte, persons [][]byte, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewAssignTesraIDsToRoleTransaction(gasPrice, gasLimit, contractAddress, admtsrId, role, persons, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Auth) NewTransferTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, newAdminId []byte, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"transfer",
		[]interface{}{
			contractAddress,
			newAdminId,
			keyIndex,
		})
}

func (this *Auth) Transfer(gasPrice, gasLimit uint64, payer, signer *Account, contractAddress common.Address, newAdminId []byte, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewTransferTransaction(gasPrice, gasLimit, contractAddress, newAdminId, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *Auth) NewVerifyTokenTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, caller []byte, funcName string, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"verifyToken",
		[]interface{}{
			contractAddress,
			caller,
			funcName,
			keyIndex,
		})
}

func (this *Auth) VerifyToken(gasPrice, gasLimit uint64, payer, signer *Account, contractAddress common.Address, caller []byte, funcName string, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewVerifyTokenTransaction(gasPrice, gasLimit, contractAddress, caller, funcName, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	if payer != nil {
		this.tesraSdk.SetPayer(tx, payer.Address)
		err = this.tesraSdk.SignToTransaction(tx, payer)
		if err != nil {
			return common.UINT256_EMPTY, err
		}
	}
	err = this.tesraSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}
