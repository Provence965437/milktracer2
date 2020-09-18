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
	ID         string
	Age        int
	Weight     int
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
	} else if opttype == "upcowinfo" { //设值
		err := UpdateCowInfo(stub, args)
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
		return shim.Success([]byte(fmt.Sprintf("invoke sucess but func is wrong:%s", opttype)))
	}
}

//newcow
func newcow(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 4 {
		return fmt.Errorf("putvalue need 4 args")
	}
	strID := args[0]
	intAGE, _ := strconv.Atoi(args[1])
	floatWEIGHT, _ := strconv.Atoi(args[2])
	var boolHEAL bool
	switch args[3] {
	case "1", "t", "T", "true", "TRUE", "True":
		boolHEAL = true
	case "0", "f", "F", "false", "FALSE", "False":
		boolHEAL = false
	}
	info := &cowinfo{
		ID:         strID,
		Age:        intAGE,
		Weight:     floatWEIGHT,
		Heal:       boolHEAL,
		Lastaction: "create"}
	cowAsBytes, _ := json.Marshal(info)
	err := stub.PutState(args[0], cowAsBytes)
	if err != nil {
		return fmt.Errorf("fail to putvalue:%s", args[0])
	}
	return nil
}

//UpdateCowInfo define
func UpdateCowInfo(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("putvalue need 2 args")
	}
	cowid := args[0]
	infotype := args[1]
	Cowinfo := cowinfo{}
	cowAsBytes, err := stub.GetState(cowid)
	if err != nil {
		return fmt.Errorf("fail to find %s data", cowid)
	}
	json.Unmarshal(cowAsBytes, &Cowinfo)
	switch infotype {
	case "action":
		Cowinfo.Lastaction = args[2]
	case "age":
		Cowinfo.Age, _ = strconv.Atoi(args[2])
	case "weight":
		Cowinfo.Weight, _ = strconv.Atoi(args[2])
	case "heal":
		if args[2] == "true" {
			Cowinfo.Heal = true
		} else {
			Cowinfo.Heal = false
		}

	default:
		return fmt.Errorf("no this infotype:%s", args[1])
	}
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
