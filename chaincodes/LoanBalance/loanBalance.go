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

type loanBalanceInfo struct {
	LoanID     string
	TxnID      string
	TxnDate    time.Time
	TxnType    string
	OpenBal    int64
	CAmt       int64
	DAmt       int64
	LoanBal    int64
	LoanStatus string
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "putTxnInfo" { //Inserting a New Business information
		return putLoanBalInfo(stub, args)
	} else if function == "getTxnInfo" { // To view a Business information
		return getLoanBalInfo(stub, args)
	}
	return shim.Success(nil)
}

func putLoanBalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 10 {
		return shim.Error("Invalid number of arguments. Needed 10 arguments")
	}

	//TxnDate -> transDate
	transDate, err := time.Parse("02/01/06", args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	txnTypeValues := map[string]bool{
		"disbursement":  true,
		"chargse":       true,
		"payment":       true,
		"other changes": true,
	}

	txnTypeLower := strings.ToLower(args[4])
	if !txnTypeValues[txnTypeLower] {
		return shim.Error("Invalid Transaction type")
	}

	openBal, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	cAmt, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	dAmt, err := strconv.ParseInt(args[7], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	loanBal, err := strconv.ParseInt(args[8], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	loanStatusValues := map[string]bool{
		"open":           true,
		"sanctioned":     true,
		"part disbursed": true,
		"disbursed":      true,
		"part collected": true,
		"collected":      true,
		"overdue":        true,
	}
	loanStatusLower := strings.ToLower(args[4])
	if !loanStatusValues[loanStatusLower] {
		return shim.Error("Invalid Loan Status type")
	}

	loanBalance := loanBalanceInfo{args[1], args[2], transDate, txnTypeLower, openBal, cAmt, dAmt, loanBal, loanStatusLower}
	loanBalanceBytes, err := json.Marshal(loanBalance)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(args[0], loanBalanceBytes)

	return shim.Success(nil)

}

func getLoanBalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Required only one argument")
	}

	loanBalance := loanBalanceInfo{}
	loanBalanceBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get the business information: " + err.Error())
	} else if loanBalanceBytes == nil {
		return shim.Error("No information is avalilable on this businessID " + args[0])
	}

	err = json.Unmarshal(loanBalanceBytes, &loanBalance)
	if err != nil {
		return shim.Error("Unable to parse into the structure " + err.Error())
	}
	jsonString := fmt.Sprintf("%+v", loanBalance)
	return shim.Success([]byte(jsonString))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s\n", err)
	}

}
