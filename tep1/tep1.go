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
package tep1

import (
	"encoding/hex"
	"fmt"
	"github.com/TesraSupernet/tesracrypto/keypair"
	tesra_go_sdk "github.com/TesraSupernet/tesrasdk"
	scomm "github.com/TesraSupernet/tesrasdk/common"
	"github.com/TesraSupernet/tesrasdk/utils"
	"github.com/TesraSupernet/Tesra/common"
	"github.com/TesraSupernet/Tesra/core/types"
	"math/big"
)

type Tep1 struct {
	ContractAddress common.Address
	sdk             *tesra_go_sdk.TesraSdk
}

func NewTep1(address common.Address, sdk *tesra_go_sdk.TesraSdk) *Tep1 {
	return &Tep1{
		ContractAddress: address,
		sdk:             sdk,
	}
}

func (this *Tep1) Name() (string, error) {
	preResult, err := this.sdk.TeoVM.PreExecInvokeTeoVMContract(this.ContractAddress,
		[]interface{}{"name", []interface{}{}})
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Tep1) Symbol() (string, error) {
	preResult, err := this.sdk.TeoVM.PreExecInvokeTeoVMContract(this.ContractAddress,
		[]interface{}{"symbol", []interface{}{}})
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Tep1) Decimals() (*big.Int, error) {
	preResult, err := this.sdk.TeoVM.PreExecInvokeTeoVMContract(this.ContractAddress,
		[]interface{}{"decimals", []interface{}{}})
	if err != nil {
		return nil, err
	}
	return preResult.Result.ToInteger()
}

func (this *Tep1) TotalSupply() (*big.Int, error) {
	preResult, err := this.sdk.TeoVM.PreExecInvokeTeoVMContract(this.ContractAddress,
		[]interface{}{"totalSupply", []interface{}{}})
	if err != nil {
		return nil, err
	}
	return preResult.Result.ToInteger()
}

func (this *Tep1) BalanceOf(account common.Address) (*big.Int, error) {
	preResult, err := this.sdk.TeoVM.PreExecInvokeTeoVMContract(this.ContractAddress,
		[]interface{}{"balanceOf", []interface{}{account}})
	if err != nil {
		return nil, err
	}
	return preResult.Result.ToInteger()
}

func (this *Tep1) Transfer(from *tesra_go_sdk.Account, to common.Address, amount *big.Int, payer *tesra_go_sdk.Account, gasPrice,
	gasLimit uint64) (common.Uint256, error) {
	return this.sdk.TeoVM.InvokeTeoVMContract(gasPrice, gasLimit, payer, from, this.ContractAddress,
		[]interface{}{"transfer", []interface{}{from.Address, to, amount}})
}

func (this *Tep1) MultiSignTransfer(fromAccounts []*tesra_go_sdk.Account, m int, to common.Address, amount *big.Int,
	gasPrice, gasLimit uint64) (common.Uint256, error) {
	pubKeys := make([]keypair.PublicKey, 0)
	for _, acc := range fromAccounts {
		pubKeys = append(pubKeys, acc.PublicKey)
	}
	fromAddr, err := types.AddressFromMultiPubKeys(pubKeys, m)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("generate multi-sign address failed, err: %s", err)
	}
	mutableTx, err := this.sdk.TeoVM.NewTeoVMInvokeTransaction(gasPrice, gasLimit, this.ContractAddress,
		[]interface{}{"transfer", []interface{}{fromAddr, to, amount}})
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("construct tx failed, err: %s", err)
	}
	for _, signer := range fromAccounts {
		err = this.sdk.MultiSignToTransaction(mutableTx, uint16(m), pubKeys, signer)
		if err != nil {
			return common.UINT256_EMPTY, fmt.Errorf("multi sign failed, err: %s", err)
		}
	}
	return this.sdk.SendTransaction(mutableTx)
}

// there are no plan to support multi sign of TransferMulti
func (this *Tep1) TransferMulti(fromAccounts []*tesra_go_sdk.Account, to []common.Address, amount []*big.Int,
	gasPrice, gasLimit uint64) (common.Uint256, error) {
	if len(fromAccounts) != len(to) || len(fromAccounts) != len(amount) || len(to) != len(amount) {
		return common.UINT256_EMPTY, fmt.Errorf("param invalid")
	}
	args := make([]*State, 0)
	for i, from := range fromAccounts {
		args = append(args, &State{
			From:   from.Address,
			To:     to[i],
			Amount: amount[i],
		})
	}
	mutableTx, err := this.sdk.TeoVM.NewTeoVMInvokeTransaction(gasPrice, gasLimit, this.ContractAddress,
		[]interface{}{"transferMulti", []interface{}{args}})
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("construct tx failed, err: %s", err)
	}
	for _, signer := range fromAccounts {
		err = this.sdk.SignToTransaction(mutableTx, signer)
		if err != nil {
			return common.UINT256_EMPTY, fmt.Errorf("sign tx failed, err: %s", err)
		}
	}
	return this.sdk.SendTransaction(mutableTx)
}

func (this *Tep1) Approve(owner *tesra_go_sdk.Account, spender common.Address, amount *big.Int, payer *tesra_go_sdk.Account, gasPrice,
	gasLimit uint64) (common.Uint256, error) {
	return this.sdk.TeoVM.InvokeTeoVMContract(gasPrice, gasLimit, payer, owner, this.ContractAddress,
		[]interface{}{"approve", []interface{}{owner.Address, spender, amount}})
}

func (this *Tep1) MultiSignApprove(ownerAccounts []*tesra_go_sdk.Account, m int, spender common.Address,
	amount *big.Int, gasPrice, gasLimit uint64) (common.Uint256, error) {
	pubKeys := make([]keypair.PublicKey, 0)
	for _, acc := range ownerAccounts {
		pubKeys = append(pubKeys, acc.PublicKey)
	}
	owner, err := types.AddressFromMultiPubKeys(pubKeys, m)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("generate multi-sign address failed, err: %s", err)
	}
	mutableTx, err := this.sdk.TeoVM.NewTeoVMInvokeTransaction(gasPrice, gasLimit, this.ContractAddress,
		[]interface{}{"approve", []interface{}{owner, spender, amount}})
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("construct tx failed, err: %s", err)
	}
	for _, signer := range ownerAccounts {
		err = this.sdk.MultiSignToTransaction(mutableTx, uint16(m), pubKeys, signer)
		if err != nil {
			return common.UINT256_EMPTY, fmt.Errorf("multi sign failed, err: %s", err)
		}
	}
	return this.sdk.SendTransaction(mutableTx)
}

func (this *Tep1) TransferFrom(spender *tesra_go_sdk.Account, from, to common.Address, amount *big.Int, payer *tesra_go_sdk.Account, gasPrice,
	gasLimit uint64) (common.Uint256, error) {
	return this.sdk.TeoVM.InvokeTeoVMContract(gasPrice, gasLimit, payer, spender, this.ContractAddress,
		[]interface{}{"transferFrom", []interface{}{spender.Address, from, to, amount}})
}

func (this *Tep1) MultiSignTransferFrom(spenders []*tesra_go_sdk.Account, m int, from, to common.Address,
	amount *big.Int, gasPrice, gasLimit uint64) (common.Uint256, error) {
	pubKeys := make([]keypair.PublicKey, 0)
	for _, acc := range spenders {
		pubKeys = append(pubKeys, acc.PublicKey)
	}
	spender, err := types.AddressFromMultiPubKeys(pubKeys, m)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("generate multi-sign address failed, err: %s", err)
	}
	mutableTx, err := this.sdk.TeoVM.NewTeoVMInvokeTransaction(gasPrice, gasLimit, this.ContractAddress,
		[]interface{}{"approve", []interface{}{spender, from, to, amount}})
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("construct tx failed, err: %s", err)
	}
	for _, signer := range spenders {
		err = this.sdk.MultiSignToTransaction(mutableTx, uint16(m), pubKeys, signer)
		if err != nil {
			return common.UINT256_EMPTY, fmt.Errorf("multi sign failed, err: %s", err)
		}
	}
	return this.sdk.SendTransaction(mutableTx)
}

func (this *Tep1) FetchTxTransferEvent(hash string) ([]*Tep1TransferEvent, error) {
	contractEvt, err := this.sdk.GetSmartContractEvent(hash)
	if err != nil {
		return nil, err
	}
	return this.parseTransferEvent(contractEvt), nil
}

// TODO: fetch approve event

func (this *Tep1) FetchBlockTransferEvent(height uint32) ([]*Tep1TransferEvent, error) {
	contractEvt, err := this.sdk.GetSmartContractEventByBlock(height)
	if err != nil {
		return nil, err
	}
	result := make([]*Tep1TransferEvent, 0)
	for _, evt := range contractEvt {
		result = append(result, this.parseTransferEvent(evt)...)
	}
	return result, nil
}

func (this *Tep1) parseTransferEvent(contractEvt *scomm.SmartContactEvent) []*Tep1TransferEvent {
	result := make([]*Tep1TransferEvent, 0)
	for _, notify := range contractEvt.Notify {
		addr, _ := utils.AddressFromHexString(notify.ContractAddress)
		if addr == this.ContractAddress {
			selfEvt, err := parseTep1TransferEvent(notify)
			if err == nil {
				result = append(result, selfEvt)
			}
		}
	}
	return result
}

func parseTep1TransferEvent(notify *scomm.NotifyEventInfo) (*Tep1TransferEvent, error) {
	state, ok := notify.States.([]interface{})
	if !ok {
		return nil, fmt.Errorf("state.States is not []interface")
	}
	if len(state) != 4 {
		return nil, fmt.Errorf("state length is not 4")
	}
	eventName, ok := state[0].(string)
	if !ok {
		return nil, fmt.Errorf("state.States[0] is not string")
	}
	from, ok := state[1].(string)
	if !ok {
		return nil, fmt.Errorf("state[1] is not string")
	}
	to, ok := state[2].(string)
	if !ok {
		return nil, fmt.Errorf("state[2] is not string")
	}
	amount, ok := state[3].(string)
	if !ok {
		return nil, fmt.Errorf("state[3] is not uint64")
	}
	evt, err := hex.DecodeString(eventName)
	if err != nil {
		return nil, fmt.Errorf("decode event name failed, err: %s", err)
	}
	fr, err := common.HexToBytes(from)
	if err != nil {
		return nil, fmt.Errorf("HexToBytes, err: %s", err)
	}
	fromAddr, err := utils.AddressParseFromBytes(fr)
	if err != nil {
		return nil, fmt.Errorf("decode from failed, err: %s", err)
	}
	toBs, err := common.HexToBytes(to)
	if err != nil {
		return nil, fmt.Errorf("HexToBytes, err: %s", err)
	}
	toAddr, err := utils.AddressParseFromBytes(toBs)
	if err != nil {
		return nil, fmt.Errorf("decode to failed, err: %s", err)
	}
	value, err := hex.DecodeString(amount)
	if err != nil {
		return nil, fmt.Errorf("decode value failed, err: %s", err)
	}
	return &Tep1TransferEvent{
		Name:   string(evt),
		From:   fromAddr,
		To:     toAddr,
		Amount: common.BigIntFromTeoBytes(value),
	}, nil
}
