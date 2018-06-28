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

type loanInfo struct {
	InstNum            string //Instrument Number
	ExposureBusinessID string
	ProgramID          string
	SanctionAmt        int64
	SanctionDate       time.Time
	SanctionAuthority  string
	ROI                float64
	DueDate            time.Time
	ValueDate          time.Time
	LoanStatus         string
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "newLoanInfo" {
		return newLoanInfo(stub, args)
	} else if function == "getLoanInfo" {
		return getLoanInfo(stub, args)
	}
	return shim.Success(nil)
}

func newLoanInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 10 {
		return shim.Error("Invalid number of arguments")
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

	//SanctionAmt -> sAmt
	sAmt, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing
	sDateStr := args[5][:8]
	sTime := args[5][9:]
	sStr := sDateStr + "T" + sTime

	//SanctionDate ->sDate
	sDate, err := time.Parse("02/01/06T15:04:05", sStr)
	if err != nil {
		return shim.Error(err.Error())
	}

	roi, err := strconv.ParseFloat(args[7], 32)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Parsing into date for storage but hh:mm:ss will also be stored as
	//00:00:00 .000Z with the date
	//DueDate -> dDate
	dDate, err := time.Parse("02/01/2006", args[8])
	if err != nil {
		return shim.Error(err.Error())
	}

	//Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing
	vDateStr := args[5][:8]
	vTime := args[5][9:]
	vStr := vDateStr + "T" + vTime

	//ValueDate ->vDate
	vDate, err := time.Parse("02/01/06T15:04:05", vStr)
	if err != nil {
		return shim.Error(err.Error())
	}

	loanStatusValues := map[string]bool{
		"open":              true,
		"sanctioned":        true,
		"part disbursed":    true,
		"disbursed":         true,
		"part collected":    true,
		"collected/settled": true,
		"overdue":           true,
	}

	loanStatusValuesLower := strings.ToLower(args[6])
	if !loanStatusValues[loanStatusValuesLower] {
		return shim.Error("Invalid Instrument Status " + args[10])
	}

	loan := loanInfo{args[1], args[2], args[3], sAmt, sDate, args[6], roi, dDate, vDate, loanStatusValuesLower}
	loanBytes, err := json.Marshal(loan)
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(args[0], loanBytes)
	return shim.Success(nil)
}

func getLoanInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid number of arguments")
	}

	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if loanBytes == nil {
		return shim.Error("No data exists on this loanID: " + args[0])
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error(err.Error())
	}
	loanString := fmt.Sprintf("%+v", loan)
	return shim.Success([]byte(loanString))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start the chaincode")
	}
}
