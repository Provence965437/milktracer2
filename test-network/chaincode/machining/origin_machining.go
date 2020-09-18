// origin_machining.go
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type machining struct{}
type machininginfo struct {
	ID     string
	Age    int
	Action string
}
type milkinfo struct {
	ID     string
	MachID string
	CowID  string
	MFD    string
	IsSold bool
}

func (t *machining) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("success invok and Not opter !!!!!! "))
}

func (t *machining) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if len(args) < 1 {
		return shim.Error("need least a arg")
	}

	var opttype = fn

	if opttype == "newmachining" {
		err := NewMachining(stub, args)
		if err != nil {
			return shim.Error("newmachining failed!")
		}
		return shim.Success([]byte("newmachining success "))

	} else if opttype == "upinfo" {
		err := UpdateMachInfo(stub, args)
		if err != nil {
			return shim.Error("upinfo failed!")
		}
		return shim.Success([]byte("upinfo success "))

	} else if opttype == "getstate" {
		keyvalue, err := GetMachState(stub, args)
		if err != nil {
			return shim.Error("getstate failed!")
		}
		MachInfo := machininginfo{}
		json.Unmarshal(keyvalue, &MachInfo)
		strMilkhInfo, _ := json.MarshalIndent(MachInfo, "", "")
		return shim.Success(strMilkhInfo)
	} else if opttype == "getmilkinfo" {
		keyvalue, err := GetMilkInfo(stub, args)
		if err != nil {
			return shim.Error("getmilkinfo failed!")
		}
		MilkInfo := milkinfo{}
		json.Unmarshal(keyvalue, &MilkInfo)
		strMilkhInfo, _ := json.MarshalIndent(MilkInfo, "", "")
		return shim.Success(strMilkhInfo)
	} else if opttype == "newmilk" {
		err := NewMilk(stub, args)
		if err != nil {
			return shim.Error("newmilk failed!")
		}
		return shim.Success([]byte("newmilk success "))

	} else if opttype == "getmachininghistory" {
		jsonKey, err := GetMachHistory(stub, args)
		if err != nil {
			return shim.Error("gethistory failed!")
		}
		return shim.Success([]byte(jsonKey))
	} else {
		return shim.Success([]byte("success invoke but func is invalid! "))
	}
}

//NewMachining ...
func NewMachining(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("putvalue need 2 args")
	}
	strID := args[0]
	nAge, _ := strconv.Atoi(args[1])
	strCreate := "create"

	info := &machininginfo{
		ID:     strID,
		Age:    nAge,
		Action: strCreate}
	MachAsBytes, _ := json.Marshal(info)
	err := stub.PutState(strID, MachAsBytes)
	if err != nil {
		return fmt.Errorf("fail to putvalue:%s", strID)
	}
	return nil
}

func NewMilk(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("putvalue need 2 args")
	}
	strID := args[0]
	strMachID := args[1]
	strCowID := args[2]
	bIsSold := false
	tm := time.Unix(int64(time.Now().Unix()), 0)
	strDate := tm.Format("2006-01-02 03:04:05 PM")

	info := &milkinfo{
		ID:     strID,
		MachID: strMachID,
		CowID:  strCowID,
		MFD:    strDate,
		IsSold: bIsSold}

	MilkAsBytes, _ := json.Marshal(info)
	err := stub.PutState(strID, MilkAsBytes)
	if err != nil {
		return fmt.Errorf("fail to putvalue:%s", strID)
	}
	return nil
}

func UpdateMachInfo(stub shim.ChaincodeStubInterface, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("putvalue need 3 args")
	}
	strID := args[0]
	strType := args[1]
	MachInfo := machininginfo{}
	machAsBytes, err := stub.GetState(strID)
	if err != nil {
		return fmt.Errorf("fail to find %s data", strID)
	}

	json.Unmarshal(machAsBytes, &MachInfo)

	switch strType {
	case "action":
		MachInfo.Action = args[2]
	case "age":
		MachInfo.Age, _ = strconv.Atoi(args[2])
	default:
		return fmt.Errorf("no this infotype:%s", args[1])
	}
	machAsBytes, _ = json.Marshal(MachInfo)

	err = stub.PutState(args[0], machAsBytes)
	if err != nil {
		return fmt.Errorf("fail to putvalue:%s", args[0])
	}
	return nil
}

func GetMachState(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("getstate need 1 args")
	}
	MachState, err := stub.GetState(args[0])
	if err != nil {
		return nil, fmt.Errorf("fail getstate%s", args[0])
	}
	return MachState, nil
}

func GetMilkInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("getstate need 1 args")
	}
	MilkInfo, err := stub.GetState(args[0])
	if err != nil {
		return nil, fmt.Errorf("fail getstate%s", args[0])
	}
	return MilkInfo, nil
}

func GetMachHistory(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("GetMachHistory need 1 args")
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
func main() {
	err := shim.Start(new(machining))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}
