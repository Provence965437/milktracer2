// origin_dairyfarm.go
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type dairyfarm struct{}
type cowinfo struct {
	ID         int
	Age        int
	Eight      int
	Heal       bool
	Lastaction string
}

//init
func (t *dairyfarm) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("suceess invok and Not opter!!!!!!"))
}

//invoke
func (t *dairyfarm) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters() //get function and args

	if len(args) < 1 {
		return shim.Error("need least a arg")
	}

	var opttype = fn

	if opttype == "newcow" { //new cow up chain
		err := newcow(stub, args)
		if err != nil {
			return shim.Error("newcow failed!")
		}
		return shim.Success([]byte("newcow success "))
	} else if opttype == "putvalue" { //设值
		err := putvalue(stub, args)
		if err != nil {
			return shim.Error("putvalue failed!")
		}
		return shim.Success([]byte("putvalue success "))
	} else if opttype == "getlastvalue" { //取值
		keyvalue, err := getlastvalue(stub, args)
		if err != nil {
			return shim.Error("getlastvalue failed!")
		}
		Cowinfo := cowinfo{}
		json.Unmarshal(keyvalue, &Cowinfo)
		strCowinfo, _ := json.MarshalIndent(Cowinfo, "", "")
		return shim.Success(strCowinfo)
	} else if opttype == "gethistory" {
		jsonKey, err := gethistory(stub, args)
		if err != nil {
			return shim.Error("gethistory failed!")
		}
		return shim.Success([]byte(jsonKey))
	} else {
		return shim.Success([]byte("success invoke but func is invalid! "))
	}
}

//newcow
func newcow(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 4 {
		return fmt.Errorf("putvalue need 4 args")
	}
	intID, _ := strconv.Atoi(args[0])
	intAGE, _ := strconv.Atoi(args[1])
	floatEIGHT, _ := strconv.Atoi(args[2])
	var boolHEAL bool
	switch args[3] {
	case "1", "t", "T", "true", "TRUE", "True":
		boolHEAL = true
	case "0", "f", "F", "false", "FALSE", "False":
		boolHEAL = false
	}
	info := &cowinfo{
		ID:         intID,
		Age:        intAGE,
		Eight:      floatEIGHT,
		Heal:       boolHEAL,
		Lastaction: "create"}
	cowAsBytes, _ := json.Marshal(info)
	err := stub.PutState(args[0], cowAsBytes)
	if err != nil {
		return fmt.Errorf("fail to putvalue:%s", args[0])
	}
	return nil
}

//putvalue
func putvalue(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("putvalue need 2 args")
	}
	Cowinfo := cowinfo{}
	cowAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return fmt.Errorf("fail to find %s data", args[0])
	}
	json.Unmarshal(cowAsBytes, &Cowinfo)
	Cowinfo.Lastaction = args[1]
	cowAsBytes, _ = json.Marshal(Cowinfo)

	err = stub.PutState(args[0], cowAsBytes)
	if err != nil {
		return fmt.Errorf("fail to putvalue:%s", args[0])
	}
	return nil
}

// getlastvalue
func getlastvalue(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("putvalue need 1 args")
	}
	cowAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, fmt.Errorf("fail getlastvalue%s", args[0])
	}
	return cowAsBytes, nil
}

//gethistory
func gethistory(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("assetname need 1 args")
	}
	assetname := args[0]
	keysIter, err := stub.GetHistoryForKey(assetname)

	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("GetHistoryForKey failed.Error accessing state %s", err))
	}
	defer keysIter.Close()
	var keys []string
	for keysIter.HasNext() {
		response, iterErr := keysIter.Next()
		if iterErr != nil {
			return "", fmt.Errorf(fmt.Sprintf("GetHistoryForKey operation failed.Error accessing state %s", err))
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

//main
func main() {
	err := shim.Start(new(dairyfarm))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode:%s", err)
	}
}
