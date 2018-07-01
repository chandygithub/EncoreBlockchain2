package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/common/util"
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
	if function == "putLoanBalInfo" {
		return putLoanBalInfo(stub, args)
	} else if function == "getLoanBalInfo" {
		return getLoanBalInfo(stub, args)
	} else if function == "updateLoanBal" {
		return updateLoanBal(stub, args)
	}
	return shim.Error("No function named " + function + " in loanBalance")
}

func putLoanBalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 10 {
		return shim.Error("Invalid number of arguments in putLoanBalInfo. Needed 10 arguments")
	}

	//TxnDate -> transDate
	transDate, err := time.Parse("02/01/2006", args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	txnTypeValues := map[string]bool{
		"disbursement":  true,
		"charges":       true,
		"payment":       true,
		"other changes": true,
	}

	txnTypeLower := strings.ToLower(args[4])
	if !txnTypeValues[txnTypeLower] {
		return shim.Error("Invalid Transaction type " + txnTypeLower)
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

	loanStatusLower := strings.ToLower(args[9])
	if !loanStatusValues[loanStatusLower] {
		return shim.Error("Invalid Loan Status type " + loanStatusLower)
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
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("Required only one argument in getLoanBalInfo, given:" + xLenStr)
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

func updateLoanBal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	/*
			// From Disbursement
		*LoanID  -> args[1]
		*TxnID   -> args[2]
		*TxnDate -> args[3]
		*TxnType -> args[4]
		*DAmt    -> args[5]


		*OpenBal -> LoanBalance from Loan structure
		*CAmt

		*LoanBal -> OpenBal-DAmt+Camt
		*LoanStatus -> depends
	*/

	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 6 {
		return shim.Error("Required 6 arguments in updateLoanBal from Disbursement")
	}

	chaincodeArgs := util.ToChaincodeArgs("getLoanInfo", args[1])
	fmt.Println("calling the other chaincode")
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}
	// spliting the arguments got from loan as response (sanction -> [0] and status -> [1])
	loanArgs := strings.Split(string(response.Payload), ",")
	fmt.Printf("disbursement payload:%s\n", loanArgs)
	fmt.Println("loanArgs[0]", loanArgs[0])
	fmt.Println("[0] type:", reflect.TypeOf(loanArgs[0]))
	fmt.Println("loanArgs[1]", loanArgs[1])
	fmt.Println("[1] type:", reflect.TypeOf(loanArgs[1]))

	openBal, err := strconv.ParseInt(loanArgs[0], 10, 64)
	if err != nil {
		return shim.Error("Error in parsing the openbalance in LoanBalance: " + err.Error())
	}
	CAmt := int64(0)
	DAmt, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return shim.Error("Error in parsing the DAmt in LoanBalance: " + err.Error())
	}

	var status string
	status = loanArgs[1] // status of the current loan
	loanBal := openBal - DAmt + CAmt
	loanBalString := strconv.FormatInt(loanBal, 10)
	if status == "sanctioned" || status == "partly disbursed" {

		if openBal-loanBal == 0 {
			status = "disbursed"
		} else {
			status = "partly disbursed"
		}
	}

	//Updating loanBalance ledger

	loanBalance := loanBalanceInfo{}
	loanBalanceBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get the loan balance information: " + err.Error())
	} else if loanBalanceBytes == nil {
		return shim.Error("No information is avalilable on this loan balance " + args[0])
	}

	err = json.Unmarshal(loanBalanceBytes, &loanBalance)
	if err != nil {
		return shim.Error("Unable to parse loan balance into the structure " + err.Error())
	}

	loanBalance.LoanStatus = status
	loanBalance.CAmt = CAmt
	loanBalance.DAmt = DAmt
	loanBalance.LoanBal = loanBal

	loanBalanceBytes, _ = json.Marshal(loanBalance)
	stub.PutState(args[0], loanBalanceBytes)
	fmt.Println("written into loan balance ledger")

	fmt.Printf("Status:%s\n", status)
	chaincodeArgs = util.ToChaincodeArgs("updateLoanInfo", args[1], status, loanBalString)
	fmt.Println("calling the other chaincode")
	response = stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s\n", err)
	}

}
