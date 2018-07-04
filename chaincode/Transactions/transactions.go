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
	TxnType string    `json:"TxnType"`      //args[1]
	TxnDate time.Time `json:"TxnDate"`      //args[2]
	LoanID  string    `json:"LoanID"`       //args[3]
	InsID   string    `json:"InstrumentID"` //args[4]
	Amt     int64     `json:"TxnAmount"`    //args[5]
	FromID  string    `json:"From"`         //args[6]
	ToID    string    `json:"To"`           //args[7]
	By      string    `json:"By"`           //args[8]
	PprID   string    `json:"PPR_ID"`       //args[9]
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "newTxnInfo" {
		return c.newTxnInfo(stub, args)
	} else if function == "getTxnInfo" {
		return c.getTxnInfo(stub, args)
	}
	return shim.Success(nil)
}

func (c *chainCode) newTxnInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 10 {
		return shim.Error("Invalid number of arguments for transaction")
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
		return shim.Error("Date error, txncc: newTxnInfo, " + err.Error())
	}

	amt, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return shim.Error("Error in converting amt to string in transactions:newTxnInfo(), " + err.Error())
	}

	switch tTypeLower {

	case "disbursement":
		argsStr := strings.Join(args, ",")
		chaincodeArgs := toChaincodeArgs("newTxnInfo", argsStr)
		fmt.Println("calling the disbursement chaincode")
		response := stub.InvokeChaincode("disbursementcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error(response.Message)
		}
		//TODO: put it at last for redability
		transaction := transactionInfo{tTypeLower, tDate, args[3], args[4], amt, args[6], args[7], args[8], args[9]}
		txnBytes, err := json.Marshal(transaction)
		err = stub.PutState(args[0], txnBytes)
		if err != nil {
			return shim.Error("Cannot write into ledger the transactino details")
		} else {
			fmt.Println("Successfully inserted the transaction " + args[0] + " into the ledger")
		}

	case "charges":
		argsStr := strings.Join([]string{args[2], args[3], args[4], args[6], args[7], args[5], args[8], args[1]}, ",")
		fmt.Println("the Charges arguments: " + argsStr)
		chaincodeArgs := toChaincodeArgs("newTxnInfo", argsStr)
		fmt.Println("calling the charges chaincode")
		response := stub.InvokeChaincode("chargescc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error(response.Message)
		}

	default:
		fmt.Println("incorrect txnType")
		return shim.Error("incorrect txnType from txncc")
	}

	return shim.Success(nil)
}

func (c *chainCode) getTxnInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Need only one argument for getting txn info")
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
	// marshal and return????
	return shim.Success([]byte(tString))

}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start the chaincode")
	}
}
