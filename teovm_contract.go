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
	"fmt"
	sdkcom "github.com/TesraSupernet/tesrasdk/common"
	"github.com/TesraSupernet/Tesra/cmd/utils"
	"github.com/TesraSupernet/Tesra/common"
	"github.com/TesraSupernet/Tesra/core/payload"
	"github.com/TesraSupernet/Tesra/core/types"
	httpcom "github.com/TesraSupernet/Tesra/http/base/common"
)

type TeoVMContract struct {
	tesraSdk *TesraSdk
}

func newTeoVMContract(tesraSdk *TesraSdk) *TeoVMContract {
	return &TeoVMContract{
		tesraSdk: tesraSdk,
	}
}

func (this *TeoVMContract) NewDeployTeoVMCodeTransaction(gasPrice, gasLimit uint64, contract payload.DeployCode) (*types.MutableTransaction, error) {

	return utils.NewDeployCodeTransaction(gasPrice, gasLimit, contract.GetRawCode(), payload.TEOVM_TYPE, contract.Name,
		contract.Version, contract.Author, contract.Email, contract.Description)
}

//DeploySmartContract Deploy smart contract to Tesra
func (this *TeoVMContract) DeployTeoVMSmartContract(
	gasPrice,
	gasLimit uint64,
	singer *Account,
	needStorage bool,
	code,
	name,
	version,
	author,
	email,
	desc string) (common.Uint256, error) {
	codeBs, err := common.HexToBytes(code)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	tx, err := utils.NewDeployCodeTransaction(gasPrice, gasLimit, codeBs, payload.TEOVM_TYPE, name, version, author, email, desc)
	err = this.tesraSdk.SignToTransaction(tx, singer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.tesraSdk.SendTransaction(tx)
}

func (this *TeoVMContract) NewTeoVMInvokeTransaction(
	gasPrice,
	gasLimit uint64,
	contractAddress common.Address,
	params []interface{},
) (*types.MutableTransaction, error) {
	invokeCode, err := httpcom.BuildTeoVMInvokeCode(contractAddress, params)
	if err != nil {
		return nil, err
	}
	return this.tesraSdk.NewInvokeTransaction(gasPrice, gasLimit, invokeCode), nil
}

func (this *TeoVMContract) InvokeTeoVMContract(
	gasPrice,
	gasLimit uint64,
	payer,
	signer *Account,
	contractAddress common.Address,
	params []interface{}) (common.Uint256, error) {
	tx, err := this.NewTeoVMInvokeTransaction(gasPrice, gasLimit, contractAddress, params)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("NewTeoVMInvokeTransaction error:%s", err)
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

func (this *TeoVMContract) PreExecInvokeTeoVMContract(
	contractAddress common.Address,
	params []interface{}) (*sdkcom.PreExecResult, error) {
	tx, err := this.NewTeoVMInvokeTransaction(0, 0, contractAddress, params)
	if err != nil {
		return nil, err
	}
	return this.tesraSdk.PreExecTransaction(tx)
}
