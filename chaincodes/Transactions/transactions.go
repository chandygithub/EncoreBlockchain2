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

type transactionInfo struct {
	TxnType string
	TxnDate time.Time
	LoanID  string
	InsID   string
	Amt     int64
	FromID  string
	ToID    string
	By      string
	PprID   string
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "newTxnInfo" {
		return newTxnInfo(stub, args)
	} else if function == "getTxnInfo" {
		return getTxnInfo(stub, args)
	}
	return shim.Success(nil)
}

func newTxnInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 9 {
		return shim.Error("Invalid number of arguments")
	}

	tTypeValues := map[string]bool{
		"disbursement": true,
		"collection":   true,
		"refund":       true,
	}

	//Converting into lower case for comparison
	tTypeLower := strings.ToLower(args[1])
	if !tTypeValues[tTypeLower] {
		return shim.Error("Invalid transaction type " + args[1])
	}

	//TxnDate -> tDate
	tDate, err := time.Parse("02/01/2006", args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	amt, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	transaction := transactionInfo{tTypeLower, tDate, args[3], args[4], amt, args[6], args[7], args[8], args[9]}
	txnBytes, err := json.Marshal(transaction)
	err = stub.PutState(args[0], txnBytes)
	return shim.Success(nil)
}

func getTxnInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid number of arguments")
	}

	txnBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if txnBytes == nil {
		return shim.Error("No data exists on this loanID: " + args[0])
	}

	transaction := transactionInfo{}
	err = json.Unmarshal(txnBytes, transaction)
	if err != nil {
		return shim.Error(err.Error())
	}

	tString := fmt.Sprintf("%+v", transaction)
	return shim.Success([]byte(tString))

}
func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start the chaincode")
	}
}
