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

type txnBalanceInfo struct {
	TxnID      string
	TxnDate    time.Time
	LoanID     string
	InsID      string
	WalletID   string
	OpeningBal int64
	TxnType    string
	Amt        int64
	CAmt       int64
	DAmt       int64
	TxnBal     int64
	By         string
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "putTxnInfo" { //Inserting a New Business information
		return putTxnInfo(stub, args)
	} else if function == "getTxnInfo" { // To view a Business information
		return getTxnInfo(stub, args)
	}
	return shim.Success(nil)
}

func putTxnInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 13 {
		return shim.Error("Invalid number of arguments. Needed 13 arguments")
	}

	//TxnDate ->txnDate
	txnDate, err := time.Parse("02/01/06", args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	openBal, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	txnTypeValues := map[string]bool{
		"loan Sanction":          true,
		"disbursement":           true,
		"charges":                true,
		"repayment / collection": true,
		"margin refund":          true,
		"interest refund":        true,
		"tds":                    true,
		"penal charges":          true,
		"cersai carges":          true,
		"factor regn charges":    true,
	}

	txnTypeLower := strings.ToLower(args[7])
	if !txnTypeValues[txnTypeLower] {
		return shim.Error("Invalid Transaction type")
	}

	amt, err := strconv.ParseInt(args[8], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	cAmt, err := strconv.ParseInt(args[9], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	dAmt, err := strconv.ParseInt(args[10], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	txnBal, err := strconv.ParseInt(args[11], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	txnBalance := txnBalanceInfo{args[1], txnDate, args[3], args[4], args[5], openBal, txnTypeLower, amt, cAmt, dAmt, txnBal, args[12]}
	txnBalanceBytes, err := json.Marshal(txnBalance)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(args[0], txnBalanceBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

func getTxnInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Required only one argument")
	}

	txnBalance := txnBalanceInfo{}
	txnBalanceBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get the business information: " + err.Error())
	} else if txnBalanceBytes == nil {
		return shim.Error("No information is avalilable on this businessID " + args[0])
	}

	err = json.Unmarshal(txnBalanceBytes, &txnBalance)
	if err != nil {
		return shim.Error("Unable to parse into the structure " + err.Error())
	}
	jsonString := fmt.Sprintf("%+v", txnBalance)
	return shim.Success([]byte(jsonString))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s\n", err)
	}

}
