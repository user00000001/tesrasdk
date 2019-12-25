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
	"fmt"
	"github.com/TesraSupernet/tesracrypto/keypair"
	tesra_go_sdk "github.com/TesraSupernet/tesrasdk"
	"github.com/TesraSupernet/tesrasdk/utils"
	"github.com/TesraSupernet/Tesra/common"
	"github.com/TesraSupernet/Tesra/core/types"
	"math/big"
	"testing"
	"time"
)

const scriptHash = "5e0aebb3dcc7af619e019a8f2195151d4d59644d"

func TestTep1(t *testing.T) {
	contractAddr, err := utils.AddressFromHexString(scriptHash)
	if err != nil {
		t.Fatal(err)
	}
	tesraSdk := tesra_go_sdk.NewTesraSdk()
	tesraSdk.NewRpcClient().SetAddress("http://polaris1.tsr.io:20336")
	tep1 := NewTep1(contractAddr, tesraSdk)
	name, err := tep1.Name()
	if err != nil {
		t.Fatal(err)
	}
	symbol, err := tep1.Symbol()
	if err != nil {
		t.Fatal(err)
	}
	decimals, err := tep1.Decimals()
	if err != nil {
		t.Fatal(err)
	}
	totalSupply, err := tep1.TotalSupply()
	if err != nil {
		t.Fatal(err)
	}

	wallet, err := tesraSdk.OpenWallet("../../wallet.json")
	if err != nil {
		fmt.Println("OpenWallet error:", err)
		return
		t.Fatal(err)
	}
	if wallet.GetAccountCount() < 2 {
		t.Fatal("account not enough")
	}
	acc, err := wallet.GetDefaultAccount([]byte("passwordtest"))
	if err != nil {
		t.Fatal(err)
	}
	balance, err := tep1.BalanceOf(acc.Address)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("name %s, symbol %s, decimals %d, totalSupply %d, balanceOf %s is %d",
		name, symbol, decimals, totalSupply, acc.Address.ToBase58(), balance)

	anotherAccount, err := wallet.GetAccountByIndex(2, []byte("passwordtest"))
	if err != nil {
		t.Fatal(err)
	}
	m := 2
	multiSignAddr, err := types.AddressFromMultiPubKeys([]keypair.PublicKey{acc.PublicKey, anotherAccount.PublicKey}, m)
	if err != nil {
		t.Fatal(err)
	}
	amount := big.NewInt(1000)
	gasPrice := uint64(500)
	gasLimit := uint64(500000)
	transferTx, err := tep1.Transfer(acc, multiSignAddr, amount, nil, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("transferTx %s: from %s to multi-sign addr %s, amount %d", transferTx.ToHexString(),
		acc.Address.ToBase58(), multiSignAddr.ToBase58(), amount)
	accounts := []*tesra_go_sdk.Account{acc, anotherAccount}
	transferMultiSignTx, err := tep1.MultiSignTransfer(accounts, m, acc.Address, amount, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("transferMultiSignTx %s: from %s to multi-sign addr %s, amount %d", transferMultiSignTx.ToHexString(),
		multiSignAddr.ToBase58(), acc.Address.ToBase58(), amount)
	approveTx, err := tep1.Approve(acc, multiSignAddr, amount, nil, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("approveTx %s: owner %s approve to multi-sign spender addr %s, amount %d", approveTx.ToHexString(),
		acc.Address.ToBase58(), multiSignAddr.ToBase58(), amount)
	multiSignTransferFromTx, err := tep1.MultiSignTransferFrom(accounts, m, acc.Address, multiSignAddr, amount,
		gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("multiSignTransferFromTx %s: owner %s, multi-sign spender addr %s, to %s, amount %d",
		multiSignTransferFromTx.ToHexString(), acc.Address.ToBase58(), multiSignAddr.ToBase58(), multiSignAddr.ToBase58(),
		amount)
	multiSignApproveTx, err := tep1.MultiSignApprove(accounts, m, acc.Address, amount, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("multiSignApproveTx %s: multi-sign owner %s approve to spender addr %s, amount %d",
		multiSignApproveTx.ToHexString(), multiSignAddr.ToBase58(), acc.Address.ToBase58(), amount)
	transferFromTx, err := tep1.TransferFrom(acc, multiSignAddr, acc.Address, amount, nil, gasPrice, gasLimit)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("transferFromTx %s: multi-sign owner %s, spender addr %s, to %s, amount %d",
		transferFromTx.ToHexString(), multiSignAddr.ToBase58(), acc.Address.ToBase58(), acc.Address.ToBase58(), amount)
	_, _ = tesraSdk.WaitForGenerateBlock(30 * time.Second)

	eventsFromTx, err := tep1.FetchTxTransferEvent(transferTx.ToHexString())
	if err != nil {
		t.Fatal(err)
	}
	for _, evt := range eventsFromTx {
		t.Logf("tx %s transfer event: %s", transferTx.ToHexString(), evt.String())
	}

	height := uint32(1791727)
	eventsFromBlock, err := tep1.FetchBlockTransferEvent(height)
	if err != nil {
		t.Fatal(err)
	}
	for _, evt := range eventsFromBlock {
		t.Logf("block %d transfer event: %s", height, evt.String())
	}
}

func TestTep1_FetchTxTransferEvent(t *testing.T) {
	contractAddr, err := utils.AddressFromHexString(scriptHash)
	if err != nil {
		t.Fatal(err)
	}
	//from address
	bs, _ := common.HexToBytes("83c12e967885ba0a1285a0c628acbfb1185af8bc")
	addr, _ := common.AddressParseFromBytes(bs)
	fmt.Println(addr.ToBase58())
	tesraSdk := tesra_go_sdk.NewTesraSdk()
	tesraSdk.NewRpcClient().SetAddress("http://polaris1.tsr.io:20336")
	tep1 := NewTep1(contractAddr, tesraSdk)
	res, _ := tep1.FetchTxTransferEvent("8074fabad95400c6705478593f2b2fce865aa356c166e63214d8a9af036ee739")
	fmt.Println(res)
}
