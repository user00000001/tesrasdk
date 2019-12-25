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
package main

import (
	"fmt"
	sdk "github.com/TesraSupernet/tesrasdk"

	"encoding/binary"
	"github.com/TesraSupernet/tesrasdk/utils"
	"github.com/TesraSupernet/Tesra/common"
	"io/ioutil"
	"time"
)

func main() {
	fmt.Println("==========================start============================")
	//testUrl := "http://127.0.0.1:20336"
	testUrl := "http://polaris1.tsr.io:20336"
	//initialize tesraSdk
	tesraSdk := sdk.NewTesraSdk()
	//suppose you already start up a local wasm TesraSupernet node
	tesraSdk.NewRpcClient().SetAddress(testUrl)
	//your wallet file
	wallet, err := tesraSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("error in OpenWallet:%s\n", err)
		return
	}

	//modify me
	walletpassword := "123456"

	//we get the first account of the wallet by your password
	signer, err := wallet.GetAccountByAddress("AVBzcUtgdgS94SpBmw4rDMhYA4KDq1YTzy", []byte(walletpassword))
	if err != nil {
		fmt.Printf("error in GetDefaultAccount:%s\n", err)
		return
	}
	fmt.Printf("===signer address is %s\n", signer.Address.ToBase58())
	//get a compiled wasm file from tsr_cpp
	wasmfile := "./TEP1.wasm"

	//set timeout
	timeoutSec := 30 * time.Second
	address1 := "AX8opZCQBpEpYsFPKpZHNguWz2s3xpT7Wk"

	// read wasm file and get the Hex fmt string
	code, err := ioutil.ReadFile(wasmfile)
	if err != nil {
		fmt.Printf("error in ReadFile:%s\n", err)

		return
	}

	codeHash := common.ToHexString(code)

	//===========================================
	gasprice := uint64(500)
	invokegaslimit := uint64(200000)
	deploygaslimit := uint64(200000000)
	// deploy the wasm contract
	fmt.Println("======DeployWasmVMSmartContract ==========")
	txHash, err := tesraSdk.WasmVM.DeployWasmVMSmartContract(
		gasprice,
		deploygaslimit,
		signer,
		codeHash,
		"TEP1 wasm",
		"1.0",
		"author",
		"email",
		"desc",
	)
	if err != nil {
		fmt.Printf("error in DeployWasmVMSmartContract:%s\n", err)

		return
	}
	_, err = tesraSdk.WaitForGenerateBlock(timeoutSec)
	if err != nil {
		fmt.Printf("error in WaitForGenerateBlock:%s\n", err)

		return
	}
	fmt.Printf("the deploy contract txhash is %s\n", txHash.ToHexString())

	//calculate the contract address from code
	contractAddr, err := utils.GetContractAddress(codeHash)
	if err != nil {
		fmt.Printf("error in GetContractAddress:%s\n", err)

		return
	}
	fmt.Printf("the contractAddr is %s\n", contractAddr.ToBase58())

	fmt.Println("======InvokeWasmVMSmartContract init==========")

	//============================================
	//invoke wasm method
	//we invoke "init" method first
	txHash, err = tesraSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, signer, contractAddr, "init", []interface{}{})
	if err != nil {
		fmt.Printf("error in InvokeWasmVMSmartContract:%s\n", err)
		return
	}
	_, err = tesraSdk.WaitForGenerateBlock(timeoutSec)
	if err != nil {
		fmt.Printf("error in WaitForGenerateBlock:%s\n", err)
		return
	}
	fmt.Printf("init txhash is :%s\n", txHash.ToHexString())
	//get smartcontract event by txhash
	events, err := tesraSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		fmt.Printf("error in GetSmartContractEvent:%s\n", err)

		return
	}
	fmt.Printf("event is %v\n", events)
	//State = 0 means transaction failed
	if events.State == 0 {
		fmt.Printf("error in events.State is 0 failed.\n")

		return
	}
	fmt.Printf("events.Notify:%v", events.Notify)
	for _, notify := range events.Notify {
		fmt.Printf("%+v\n", notify)
	}

	//next we test transfer method
	//1.  we get another address from wallet,suppose you have created in the wallet
	account2, err := wallet.GetAccountByAddress(address1, []byte(walletpassword))
	if err != nil {
		fmt.Printf("error in GetAccountByAddress:%s\n", err)
		return
	}
	fmt.Println("======InvokeWasmVMSmartContract transfer==========")

	//2. we construct a tx transfer 500 token from signer account to account2
	txHash, err = tesraSdk.WasmVM.InvokeWasmVMSmartContract(
		gasprice, invokegaslimit, nil, signer, contractAddr, "transfer", []interface{}{signer.Address, account2.Address, uint64(500)})
	if err != nil {
		fmt.Printf("error in InvokeWasmVMSmartContract:%s\n", err)
		return
	}
	_, err = tesraSdk.WaitForGenerateBlock(timeoutSec)
	if err != nil {
		fmt.Printf("error in WaitForGenerateBlock:%s\n", err)

		return
	}
	//get smartcontract event by txhash
	events, err = tesraSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		fmt.Printf("error in GetSmartContractEvent:%s\n", err)
		return
	}
	fmt.Printf("event is %v\n", events)
	//State = 0 means transaction failed
	if events.State == 0 {
		fmt.Printf("error in events.State is 0 failed.\n")

		return
	}
	fmt.Printf("events.Notify:%v", events.Notify)
	for _, notify := range events.Notify {
		//you check the notify here
		fmt.Printf("%+v\n", notify)
	}

	//we will query the balance using pre-execuse method
	res, err := tesraSdk.WasmVM.PreExecInvokeWasmVMContract(contractAddr, "balanceOf", []interface{}{signer.Address})
	if err != nil {
		fmt.Printf("error in PreExecInvokeWasmVMContract:%s\n", err)

		return
	}
	bs, err := res.Result.ToByteArray()
	if err != nil {
		fmt.Printf("error in ToByteArray:%s\n", err)

		return
	}
	fmt.Printf("balance of %s is %d\n", signer.Address.ToBase58(), binary.LittleEndian.Uint64(bs))

	fmt.Println("==============================end ==========================")
}
