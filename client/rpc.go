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

//RPC client for Tesra
package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/TesraSupernet/tesrasdk/utils"
	"github.com/TesraSupernet/Tesra/common"
	"github.com/TesraSupernet/Tesra/core/types"
	"io/ioutil"
	"net/http"
	"time"
)

//RpcClient for TesraSupernet rpc api
type RpcClient struct {
	addr       string
	httpClient *http.Client
}

//NewRpcClient return RpcClient instance
func NewRpcClient() *RpcClient {
	return &RpcClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   5,
				DisableKeepAlives:     false, //enable keepalive
				IdleConnTimeout:       time.Second * 300,
				ResponseHeaderTimeout: time.Second * 300,
			},
			Timeout: time.Second * 300, //timeout for http response
		},
	}
}

//SetAddress set rpc server address. Simple http://localhost:20336
func (this *RpcClient) SetAddress(addr string) *RpcClient {
	this.addr = addr
	return this
}

//SetHttpClient set http client to RpcClient. In most cases SetHttpClient is not necessary
func (this *RpcClient) SetHttpClient(httpClient *http.Client) *RpcClient {
	this.httpClient = httpClient
	return this
}

//GetVersion return the version of Tesra
func (this *RpcClient) getVersion(qid string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_VERSION, []interface{}{})
}

func (this *RpcClient) getNetworkId(qid string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_NETWORK_ID, []interface{}{})
}

//GetBlockByHash return block with specified block hash in hex string code
func (this *RpcClient) getBlockByHash(qid, hash string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_BLOCK, []interface{}{hash})
}

//GetBlockByHeight return block by specified block height
func (this *RpcClient) getBlockByHeight(qid string, height uint32) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_BLOCK, []interface{}{height})
}

func (this *RpcClient) getBlockInfoByHeight(qid string, height uint32) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_BLOCK, []interface{}{height, 1})
}

//GetBlockCount return the total block count of Tesra
func (this *RpcClient) getBlockCount(qid string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_BLOCK_COUNT, []interface{}{})
}

func (this *RpcClient) getCurrentBlockHeight(qid string) ([]byte, error) {
	data, err := this.getBlockCount(qid)
	if err != nil {
		return nil, err
	}
	count, err := utils.GetUint32(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(count - 1)
}

//GetCurrentBlockHash return the current block hash of Tesra
func (this *RpcClient) getCurrentBlockHash(qid string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_CURRENT_BLOCK_HASH, []interface{}{})
}

//GetBlockHash return block hash by block height
func (this *RpcClient) getBlockHash(qid string, height uint32) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_BLOCK_HASH, []interface{}{height})
}

//GetStorage return smart contract storage item.
//addr is smart contact address
//key is the key of value in smart contract
func (this *RpcClient) getStorage(qid, contractAddress string, key []byte) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_STORAGE, []interface{}{contractAddress, hex.EncodeToString(key)})
}

//GetSmartContractEvent return smart contract event execute by invoke transaction by hex string code
func (this *RpcClient) getSmartContractEvent(qid, txHash string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_SMART_CONTRACT_EVENT, []interface{}{txHash})
}

func (this *RpcClient) getSmartContractEventByBlock(qid string, blockHeight uint32) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_SMART_CONTRACT_EVENT, []interface{}{blockHeight})
}

//GetRawTransaction return transaction by transaction hash
func (this *RpcClient) getRawTransaction(qid, txHash string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_TRANSACTION, []interface{}{txHash})
}

//GetSmartContract return smart contract deployed in TesraSupernet by specified smart contract address
func (this *RpcClient) getSmartContract(qid, contractAddress string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_SMART_CONTRACT, []interface{}{contractAddress})
}

//GetMerkleProof return the merkle proof whether tx is exist in ledger. Param txHash is in hex string code
func (this *RpcClient) getMerkleProof(qid, txHash string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_MERKLE_PROOF, []interface{}{txHash})
}

func (this *RpcClient) getMemPoolTxState(qid, txHash string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_MEM_POOL_TX_STATE, []interface{}{txHash})
}

func (this *RpcClient) getMemPoolTxCount(qid string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_MEM_POOL_TX_COUNT, []interface{}{})
}

func (this *RpcClient) getBlockHeightByTxHash(qid, txHash string) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_BLOCK_HEIGHT_BY_TX_HASH, []interface{}{txHash})
}

func (this *RpcClient) getBlockTxHashesByHeight(qid string, height uint32) ([]byte, error) {
	return this.sendRpcRequest(qid, RPC_GET_BLOCK_TX_HASH_BY_HEIGHT, []interface{}{height})
}

func (this *RpcClient) sendRawTransaction(qid string, tx *types.Transaction, isPreExec bool) ([]byte, error) {
	txData := hex.EncodeToString(common.SerializeToBytes(tx))
	params := []interface{}{txData}
	if isPreExec {
		params = append(params, 1)
	}
	return this.sendRpcRequest(qid, RPC_SEND_TRANSACTION, params)
}

//sendRpcRequest send Rpc request to Tesra
func (this *RpcClient) sendRpcRequest(qid, method string, params []interface{}) ([]byte, error) {
	rpcReq := &JsonRpcRequest{
		Version: JSON_RPC_VERSION,
		Id:      qid,
		Method:  method,
		Params:  params,
	}
	data, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("JsonRpcRequest json.Marsha error:%s", err)
	}
	resp, err := this.httpClient.Post(this.addr, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("http post request:%s error:%s", data, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rpc response body error:%s", err)
	}
	rpcRsp := &JsonRpcResponse{}
	err = json.Unmarshal(body, rpcRsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal JsonRpcResponse:%s error:%s", body, err)
	}
	if rpcRsp.Error != 0 {
		return nil, fmt.Errorf("JsonRpcResponse error code:%d desc:%s result:%s", rpcRsp.Error, rpcRsp.Desc, rpcRsp.Result)
	}
	return rpcRsp.Result, nil
}
