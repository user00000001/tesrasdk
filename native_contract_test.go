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
	"github.com/TesraSupernet/tesracrypto/keypair"
	"testing"
	"time"
)

func TestTsrId_RegIDWithPublicKey(t *testing.T) {
	return
	Init()
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_RegIDWithPublicKey NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_RegIDWithPublicKey GetControllerByIndex error:%s", err)
		return
	}
	txHash, err := testTesraSdk.Native.TsrId.RegIDWithPublicKey(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_RegIDWithPublicKey RegIDWithPublicKey error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	event, err := testTesraSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		t.Errorf("TestTsrId_RegIDWithPublicKey GetSmartContractEvent error:%s", err)
		return
	}
	fmt.Printf("TestTsrId_RegIDWithPublicKey Event: %+v\n", event)
	fmt.Println("testIdentity.ID:", testIdentity.ID)
	ddo, err := testTesraSdk.Native.TsrId.GetDDO(testIdentity.ID)
	if err != nil {
		t.Errorf("TestTsrId_RegIDWithPublicKey GetDDO error:%s", err)
		return
	}
	fmt.Printf("TestTsrId_RegIDWithPublicKey DDO:%+v\n", ddo)
}

func TestTsrId_RegIDWithAttributes(t *testing.T) {
	return
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_RegIDWithPublicKey NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_RegIDWithPublicKey GetControllerByIndex error:%s", err)
		return
	}
	attributes := make([]*DDOAttribute, 0)
	attr1 := &DDOAttribute{
		Key:       []byte("Hello"),
		Value:     []byte("World"),
		ValueType: []byte("string"),
	}
	attributes = append(attributes, attr1)
	attr2 := &DDOAttribute{
		Key:       []byte("Foo"),
		Value:     []byte("Bar"),
		ValueType: []byte("string"),
	}
	attributes = append(attributes, attr2)
	_, err = testTesraSdk.Native.TsrId.RegIDWithAttributes(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, testDefController, attributes)
	if err != nil {
		t.Errorf("TestTsrId_RegIDWithPublicKey RegIDWithAttributes error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)

	ddo, err := testTesraSdk.Native.TsrId.GetDDO(testIdentity.ID)
	if err != nil {
		t.Errorf("GetDDO error:%s", err)
		return
	}

	owners := ddo.Owners
	if owners[0].Value != hex.EncodeToString(keypair.SerializePublicKey(testDefController.GetPublicKey())) {
		t.Errorf("TestTsrId_RegIDWithPublicKey pubkey %s != %s", owners[0].Value, hex.EncodeToString(keypair.SerializePublicKey(testDefController.GetPublicKey())))
		return
	}
	attrs := ddo.Attributes
	if len(attributes) != len(attrs) {
		t.Errorf("TestTsrId_RegIDWithPublicKey attribute size %d != %d", len(attrs), len(attributes))
		return
	}
	fmt.Printf("Owner:%+v\n", owners[0])
	if string(attr1.Key) != string(attrs[0].Key) ||
		string(attr1.Value) != string(attrs[0].Value) ||
		string(attr1.ValueType) != string(attrs[0].ValueType) {
		t.Errorf("TestTsrId_RegIDWithPublicKey Attribute %s != %s", attrs[0], attr1)
	}
	if string(attr2.Key) != string(attrs[1].Key) ||
		string(attr2.Value) != string(attrs[1].Value) ||
		string(attr2.ValueType) != string(attrs[1].ValueType) {
		t.Errorf("TestTsrId_RegIDWithPublicKey Attribute %s != %s", attrs[1], attr2)
	}
}

func TestTsrId_Key(t *testing.T) {
	return
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_Key NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_Key GetControllerByIndex error:%s", err)
		return
	}
	_, err = testTesraSdk.Native.TsrId.RegIDWithPublicKey(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_Key RegIDWithPublicKey error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)

	controller1, err := testIdentity.NewDefaultSettingController("2", testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_Key NewDefaultSettingController error:%s", err)
		return
	}

	_, err = testTesraSdk.Native.TsrId.AddKey(testGasPrice, testGasLimit, nil, testIdentity.ID, testDefAcc, controller1.PublicKey, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_Key AddKey error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)

	owners, err := testTesraSdk.Native.TsrId.GetPublicKeys(testIdentity.ID)
	if err != nil {
		t.Errorf("TestTsrId_Key GetPublicKeys error:%s", err)
		return
	}

	if len(owners) != 2 {
		t.Errorf("TestTsrId_Key owner size:%d != 2", len(owners))
		return
	}

	if owners[0].Value != hex.EncodeToString(keypair.SerializePublicKey(testDefController.PublicKey)) {
		t.Errorf("TestTsrId_Key owner index:%d pubkey:%s != %s", owners[0].pubKeyIndex, owners[0].Value, hex.EncodeToString(keypair.SerializePublicKey(testDefController.PublicKey)))
		return
	}

	if owners[1].Value != hex.EncodeToString(keypair.SerializePublicKey(controller1.PublicKey)) {
		t.Errorf("TestTsrId_Key owner index:%d pubkey:%s != %s", owners[1].pubKeyIndex, owners[1].Value, hex.EncodeToString(keypair.SerializePublicKey(controller1.PublicKey)))
		return
	}

	_, err = testTesraSdk.Native.TsrId.RevokeKey(testGasPrice, testGasLimit, nil, testIdentity.ID, testDefAcc, testDefController.PublicKey, controller1)
	if err != nil {
		t.Errorf("TestTsrId_Key RevokeKey error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)

	owners, err = testTesraSdk.Native.TsrId.GetPublicKeys(testIdentity.ID)
	if err != nil {
		t.Errorf("TestTsrId_Key GetPublicKeys error:%s", err)
		return
	}

	if len(owners) != 1 {
		t.Errorf("TestTsrId_Key owner size:%d != 1 after remove", len(owners))
		return
	}

	state, err := testTesraSdk.Native.TsrId.GetKeyState(testIdentity.ID, 1)
	if err != nil {
		t.Errorf("TestTsrId_Key GetKeyState error:%s", err)
		return
	}

	if state != KEY_STATUS_REVOKE {
		t.Errorf("TestTsrId_Key remove key state != %s", KEY_STATUS_REVOKE)
		return
	}

	state, err = testTesraSdk.Native.TsrId.GetKeyState(testIdentity.ID, 2)
	if err != nil {
		t.Errorf("TestTsrId_Key GetKeyState error:%s", err)
		return
	}
	if state != KEY_STSTUS_IN_USE {
		t.Errorf("TestTsrId_Key GetKeyState state != %s", KEY_STSTUS_IN_USE)
		return
	}
}

func TestTsrId_Attribute(t *testing.T) {
	return
	Init()
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_Attribute NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_Attribute GetControllerByIndex error:%s", err)
		return
	}
	_, err = testTesraSdk.Native.TsrId.RegIDWithPublicKey(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_Attribute RegIDWithPublicKey error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)

	attributes := make([]*DDOAttribute, 0)
	attr1 := &DDOAttribute{
		Key:       []byte("Foo"),
		Value:     []byte("Bar"),
		ValueType: []byte("string"),
	}
	attributes = append(attributes, attr1)
	attr2 := &DDOAttribute{
		Key:       []byte("Hello"),
		Value:     []byte("World"),
		ValueType: []byte("string"),
	}
	attributes = append(attributes, attr2)
	_, err = testTesraSdk.Native.TsrId.AddAttributes(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, attributes, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_Attribute AddAttributes error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	attrs, err := testTesraSdk.Native.TsrId.GetAttributes(testIdentity.ID)
	if len(attributes) != len(attrs) {
		t.Errorf("TestTsrId_Attribute GetAttributes len:%d != %d", len(attrs), len(attributes))
		return
	}
	fmt.Println("attrs:", attrs)
	if string(attr1.Key) != string(attrs[0].Key) || string(attr1.Value) != string(attrs[0].Value) || string(attr1.ValueType) != string(attrs[0].ValueType) {
		t.Errorf("TestTsrId_Attribute attribute:%s != %s", attrs[0], attr1)
		return
	}

	_, err = testTesraSdk.Native.TsrId.RemoveAttribute(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, attr1.Key, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_Attribute RemoveAttribute error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	attrs, err = testTesraSdk.Native.TsrId.GetAttributes(testIdentity.ID)
	if len(attrs) != 1 {
		t.Errorf("TestTsrId_Attribute GetAttributes len:%d != 1", len(attrs))
		return
	}
	if string(attr2.Key) != string(attrs[0].Key) || string(attr2.Value) != string(attrs[0].Value) || string(attr2.ValueType) != string(attrs[0].ValueType) {
		t.Errorf("TestTsrId_Attribute attribute:%s != %s", attrs[0], attr2)
		return
	}
}

func TestTsrId_Recovery(t *testing.T) {
	return
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_Recovery NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_Recovery GetControllerByIndex error:%s", err)
		return
	}
	_, err = testTesraSdk.Native.TsrId.RegIDWithPublicKey(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_Recovery RegIDWithPublicKey error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	_, err = testTesraSdk.Native.TsrId.SetRecovery(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, testDefAcc.Address, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_Recovery SetRecovery error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	ddo, err := testTesraSdk.Native.TsrId.GetDDO(testIdentity.ID)
	if err != nil {
		t.Errorf("TestTsrId_Recovery GetDDO error:%s", err)
		return
	}
	if ddo.Recovery != testDefAcc.Address.ToBase58() {
		t.Errorf("TestTsrId_Recovery recovery address:%s != %s", ddo.Recovery, testDefAcc.Address.ToBase58())
		return
	}

	acc1, err := testWallet.NewDefaultSettingAccount(testPasswd)
	if err != nil {
		t.Errorf("TestTsrId_Recovery NewDefaultSettingAccount error:%s", err)
		return
	}

	txHash, err := testTesraSdk.Native.TsrId.SetRecovery(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, acc1.Address, testDefController)

	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	evt, err := testTesraSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		t.Errorf("TestTsrId_Recovery GetSmartContractEvent:%s error:%s", txHash.ToHexString(), err)
		return
	}

	if evt != nil && evt.State == 1 {
		t.Errorf("TestTsrId_Recovery duplicate add recovery should failed")
		return
	}
	_, err = testTesraSdk.Native.TsrId.ChangeRecovery(testGasPrice, testGasLimit, nil, testDefAcc, testIdentity.ID, acc1.Address, testDefAcc.Address, testDefController)
	if err != nil {
		t.Errorf("TestTsrId_Recovery ChangeRecovery error:%s", err)
		return
	}
	testTesraSdk.WaitForGenerateBlock(30*time.Second, 1)
	ddo, err = testTesraSdk.Native.TsrId.GetDDO(testIdentity.ID)
	if err != nil {
		t.Errorf("TestTsrId_Recovery GetDDO error:%s", err)
		return
	}
	if ddo.Recovery != acc1.Address.ToBase58() {
		t.Errorf("TestTsrId_Recovery recovery address:%s != %s", ddo.Recovery, acc1.Address.ToBase58())
		return
	}
}
