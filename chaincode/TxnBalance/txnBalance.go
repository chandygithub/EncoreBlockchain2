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
	TxnID      string    `json:"TxnID"`
	TxnDate    time.Time `json:"TxnDate"`
	LoanID     string    `json:"LoanID"`
	InsID      string    `json:"InsID"`
	WalletID   string    `json:"WalletID"`
	OpeningBal int64     `json:"OpeningBalance"`
	TxnType    string    `json:"TxnType"`
	Amt        int64     `json:"Amount"`
	CAmt       int64     `json:"CreditAmount"`
	DAmt       int64     `json:"DebitAmount"`
	TxnBal     int64     `json:"TxnBalance"`
	By         string    `json:"By"`
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "putTxnBalInfo" { //Inserting a New Business information
		return c.putTxnBalInfo(stub, args)
	} else if function == "getTxnBalInfo" { // To view a Business information
		return c.getTxnBalInfo(stub, args)
	}
	return shim.Error("txnbalcc: " + "Inside txnBalcc:Invoke(), Function does not exit" + function)
}

func (c *chainCode) putTxnBalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 13 {
		return shim.Error("txnbalcc: " + "Invalid number of arguments for txnBal. Needed 13 arguments")
	}

	ifExists, err := stub.GetState(args[0])
	if ifExists != nil {
		return shim.Error("txnbalcc: " + "TxnBalanceId " + args[0] + " exits. Cannot create new ID")
	}

	//TxnDate ->txnDate
	txnDate, err := time.Parse("02/01/2006", args[2])
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in txndate " + err.Error())
	}

	openBal, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in openbal " + err.Error())
	}

	txnTypeValues := map[string]bool{
		"disbursement":              true,
		"repayment":                 true,
		"margin_refund":             true,
		"interest_refund":           true,
		"penal_interest_collection": true,
		"loan_sanction":             true,
		"charges":                   true,
		"interest_in_advance":       true,
		"accrual":                   true,
		"interest_accrued_charges":  true,
		"penal_charges":             true,
		"TDS":                       true,
	}

	txnTypeLower := strings.ToLower(args[7])
	if !txnTypeValues[txnTypeLower] {
		return shim.Error("txnbalcc: " + "txnbal Invalid Transaction type")
	}

	amt, err := strconv.ParseInt(args[8], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in amt" + err.Error())
	}

	cAmt, err := strconv.ParseInt(args[9], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in camt" + err.Error())
	}

	dAmt, err := strconv.ParseInt(args[10], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in damt " + err.Error())
	}

	txnBal, err := strconv.ParseInt(args[11], 10, 64)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal err in txnbal " + err.Error())
	}

	txnBalance := txnBalanceInfo{args[1], txnDate, args[3], args[4], args[5], openBal, txnTypeLower, amt, cAmt, dAmt, txnBal, args[12]}
	txnBalanceBytes, err := json.Marshal(txnBalance)
	if err != nil {
		return shim.Error("txnbalcc: " + err.Error())
	}
	err = stub.PutState(args[0], txnBalanceBytes)
	if err != nil {
		return shim.Error("txnbalcc: " + "txnbal cannot write to ledger: " + err.Error())
	}
	fmt.Println("Succefully wrote txnID " + args[0] + " into the ledger")
	return shim.Success(nil)

}

func (c *chainCode) getTxnBalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("txnbalcc: " + "Required only one argument")
	}

	txnBalance := txnBalanceInfo{}
	txnBalanceBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("txnbalcc: " + "Failed to get the business information: " + err.Error())
	} else if txnBalanceBytes == nil {
		return shim.Error("txnbalcc: " + "No information is avalilable on this businessID " + args[0])
	}

	err = json.Unmarshal(txnBalanceBytes, &txnBalance)
	if err != nil {
		return shim.Error("txnbalcc: " + "Unable to parse into the structure " + err.Error())
	}
	jsonString := fmt.Sprintf("%+v", txnBalance)
	fmt.Printf("Transaction info %s : %s", args[0], jsonString)
	return shim.Success([]byte(jsonString))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("txnbalcc: "+"Error starting Simple chaincode: %s\n", err)
	}
}
