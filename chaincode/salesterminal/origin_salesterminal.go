// origin_salesterminal
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type order struct {
	ID     string
	MilkID string
}
type milkinfo struct {
	ID     string
	MachID string
	CowID  string
	MFD    string
	IsSold bool
}
type salesterminal struct{}

func (t *salesterminal) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("success invok and Not opter !!!!!! "))
}

func (t *salesterminal) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if len(args) < 1 {
		return shim.Error("need least a arg")
	}

	var opttype = fn

	if opttype == "neworder" {
		err := NewOrder(stub, args)
		if err != nil {
			return shim.Error("neworder failed!")
		}
		return shim.Success([]byte("neworder success "))
	} else if opttype == "getorderstate" {
		var keyvalue []byte
		var err error
		keyvalue, err = GetOrderState(stub, args)

		if err != nil {
			return shim.Error("find error!")
		}
		return shim.Success(keyvalue)
	} else if opttype == "getorderhistory" {
		jsonKey, err := GetOrderHistory(stub, args)
		if err != nil {
			return shim.Error("gethistory failed!")
		}
		return shim.Success([]byte(jsonKey))
	} else if opttype == "getmilkhistory" {

		jsonKey, err := GetMilkHistory(stub, args)
		if err != nil {
			return shim.Error("getmilkhistory failed")
		}
		return shim.Success([]byte(jsonKey))
	} else {
		return shim.Success([]byte("success invok and No operation !!!!!!!!"))
	}
}

func NewOrder(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("putvalue need 2 args")
	}
	strID := args[0]
	strMilkID := args[1]
	info := &order{
		ID:     strID,
		MilkID: strMilkID}
	OrderAsBytes, _ := json.Marshal(info)
	err := stub.PutState(strID, OrderAsBytes)
	if err != nil {
		return fmt.Errorf("fail to neworder:%s", strID)
	}
	return nil
}
func GetOrderState(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("getstate need 1 args")
	}
	MilkState, err := stub.GetState(args[0])
	if err != nil {
		return nil, fmt.Errorf("fail getstate%s", args[0])
	}
	return MilkState, nil
}
func GetOrderHistory(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("GetOrderHistory need 1 args")
	}
	strID := args[0]
	keysIter, err := stub.GetHistoryForKey(strID)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("GetHistoryForKey failed.Error accessing state %s", err))
	}
	defer keysIter.Close()
	var keys []string
	for keysIter.HasNext() {
		response, iterErr := keysIter.Next()
		if iterErr != nil {
			return "", fmt.Errorf(fmt.Sprintf("GetHistoryForKey operation failed.Error accessing state %s ", iterErr))
		}

		//交易编号
		txid := response.TxId
		//交易的值
		txvalue := response.Value
		//当前交易的状态
		txstatus := response.IsDelete
		//交易发生的时间戳
		txtimestamp := response.Timestamp
		tm := time.Unix(txtimestamp.Seconds, 0)
		datestr := tm.Format("2006-01-02 03:04:05 PM")
		fmt.Printf("Tx info - txid:%s value: %s if delete %t datetime:%s\n", txid, string(txvalue), txstatus, datestr)
		keys = append(keys, string(txvalue)+":"+datestr)
	}

	jsonKeys, err := json.Marshal(keys)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("query operation failed.Error marshaling JSON :%s", err))
	}

	return string(jsonKeys), nil
}
func GetMilkHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, fmt.Errorf("GetMilkHistory need 1 args")
	}
	var keys []string
	MilkID := args[0]

	milk_query_parm := []string{"getmilkinfo", MilkID}
	queryArgs2 := make([][]byte, len(milk_query_parm))
	for i, arg := range milk_query_parm {
		queryArgs2[i] = []byte(arg)
	}

	milk_response := stub.InvokeChaincode("machining", queryArgs2, "mychannel")

	if milk_response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", milk_response.Payload)
		fmt.Printf(errStr)
		return nil, fmt.Errorf(errStr)
	}

	milk_query_parm1 := []string{"getmachininghistory", MilkID}
	queryArgs3 := make([][]byte, len(milk_query_parm1))
	for i, arg := range milk_query_parm1 {
		queryArgs3[i] = []byte(arg)
	}

	milk_response2 := stub.InvokeChaincode("machining", queryArgs3, "mychannel")

	if milk_response2.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", milk_response2.Payload)
		fmt.Printf(errStr)
		return nil, fmt.Errorf(errStr)
	}

	result_milkhistory := string(milk_response2.Payload)
	var strmilkhis []string
	if err := json.Unmarshal([]byte(result_milkhistory), &strmilkhis); err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("query operation failed.Error marshaling JSON:%s", err))
	}

	for _, v := range strmilkhis {
		keys = append(keys, v)
	}

	Milkinfo := &milkinfo{}
	json.Unmarshal(milk_response.Payload, Milkinfo)
	//获取工厂编号
	machID := Milkinfo.MachID
	cowID := Milkinfo.CowID

	//调用加工厂的chaincode获取加工厂的溯源信息
	machining_history_parm := []string{"getmachininghistory", machID}
	queryArgs := make([][]byte, len(machining_history_parm))
	for i, arg := range machining_history_parm {
		queryArgs[i] = []byte(arg)
	}

	response := stub.InvokeChaincode("machining", queryArgs, "mychannel")

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		fmt.Printf(errStr)
		return nil, fmt.Errorf(errStr)
	}

	//获取加工的信息
	result := string(response.Payload)

	fmt.Printf("machining info -  result : %s  \n ", result)

	var strmachinfo []string
	if err := json.Unmarshal([]byte(result), &strmachinfo); err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("query operation failed.Error marshaling JSON:%s", err))
	}

	for _, v := range strmachinfo {
		keys = append(keys, v)
	}

	//通过牛奶的编号获取溯源信息
	cow_parms := []string{"gethistory", cowID}
	queryArgs1 := make([][]byte, len(cow_parms))
	for i, arg := range cow_parms {
		queryArgs1[i] = []byte(arg)
	}

	cow_response := stub.InvokeChaincode("dairyfarm", queryArgs1, "mychannel")

	if cow_response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", cow_response.Payload)
		fmt.Printf(errStr)
		return nil, fmt.Errorf(errStr)
	}

	cow_result := string(cow_response.Payload)

	fmt.Printf("cow info - result :%s \n", cow_result)

	var cowhistorys []string
	if err := json.Unmarshal([]byte(cow_result), &cowhistorys); err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("query operation failed.Error marshaling JSON:%s", err))
	}

	for _, v1 := range cowhistorys {
		keys = append(keys, v1)
	}

	jsonKeys, err := json.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("query operation failed.Error marshaling JSON:%s", err))
	}

	return jsonKeys, nil

}
func main() {
	err := shim.Start(new(salesterminal))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}
