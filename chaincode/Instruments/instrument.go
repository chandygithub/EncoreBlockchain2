package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type instrumentInfo struct {
	InstrumentRefNo string
	InstrumenDate   time.Time
	SellBusinessID  string
	BuyBusinsessID  string
	InsAmount       int64
	InsStatus       string
	InsDueDate      time.Time
	ProgramID       string
	UploadBatchNo   string
	ValueDate       time.Time
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "enterInstrument" {
		return enterInstrument(stub, args)
	} else if function == "getInstrument" {
		return getInstrument(stub, args)
	}

	return shim.Success(nil)
}

func enterInstrument(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 11 {
		return shim.Error("Invalid number of arguments")
	}

	//InstrumentDate -> instDate
	instDate, err := time.Parse("02/01/2006", args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	insAmt, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	insStatusValues := map[string]bool{
		"open":              true,
		"sanctioned":        true,
		"part disbursed":    true,
		"disbursed":         true,
		"part collected":    true,
		"collected/settled": true,
		"overdue":           true,
	}

	insStatusValuesLower := strings.ToLower(args[6])
	if !insStatusValues[insStatusValuesLower] {
		return shim.Error("Invalid Instrument Status " + args[6])
	}

	// UNCOMMENT THIS WHILE ALL THE CHAINCODES ARE LINKED
	// SO THAT CHECKING FOR A PROGRAM ID CAN WORK PROPERLY
	/*
		//Checking if the programID exist or not
		chk, err := stub.GetState(args[2])
		if err == nil {
			return shim.Error("This program does not exist")
		} else if chk == nil {
			return shim.Error("There is no information on this program")
		}
	*/

	//InsDueDate -> insDate
	insDate, err := time.Parse("02/01/2006", args[7])
	if err != nil {
		return shim.Error(err.Error())
	}

	//Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing

	vString := args[10][:8] + "T" + args[10][9:] //removing the ":" part from the string

	//ValueDate -> vDate
	vDate, err := time.Parse("02/01/2006T15:04:05", vString)
	if err != nil {
		return shim.Error(err.Error())
	}

	inst := instrumentInfo{args[1], instDate, args[3], args[4], insAmt, insStatusValuesLower, insDate, args[8], args[9], vDate}
	instBytes, err := json.Marshal(inst)
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(args[0], instBytes)
	return shim.Success(nil)
}

func getInstrument(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid number of arguments")
	}

	insBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if insBytes == nil {
		return shim.Error("No data exists on this InstrumentID: " + args[0])
	}

	ins := instrumentInfo{}
	err = json.Unmarshal(insBytes, &ins)
	insString := fmt.Sprintf("%+v", ins)
	return shim.Success([]byte(insString))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start the chaincode")
	}
}
