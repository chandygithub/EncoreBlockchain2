package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
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

	if function == "newPICinfo" {
		return newPICinfo(stub, args)
	}
	return shim.Error("penalinterestcc: " + "no function named " + function + " found in Interest Refund")
}

func newPICinfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("penalCharges.cc: " + "Invalid number of arguments in newPICinfo(Interest Refund) (required:10) given:" + xLenStr)
	}

	/*
	 *TxnType string    //args[1]
	 *TxnDate time.Time //args[2]
	 *LoanID  string    //args[3] loan
	 *InsID   string    //args[4]
	 *Amt     int64     //args[5]
	 *FromID  string    //args[6]  Bank
	 *ToID    string    //args[7]  Seller
	 *By      string    //args[8]
	 *PprID   string    //args[9]
	 */

	amt, _ := strconv.ParseInt(args[5], 10, 64)
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 				UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// The transaction object has been created and written into the ledger
	// The JSON file is 'transaction'function
	// Now to create a TXN_Bal_Update obj for 4 times
	// Calling TXN_Balance CC based on TXN_Type
	/*
			    a. Debiting (Decreasing) Business Wallet
		        b. Crediting (Incresing) Bank Wallet
		        c. Debiting (Decreasing) Business Charges O/s Wallet
		        d. Debiting (Decreasing) Loan Charges Wallet
	*/

	//Validations

	// Must be Existing Loan with Status as Overdue
	chaincodeArgs := toChaincodeArgs("getLoanStatus", args[3])
	response := stub.InvokeChaincode("loancc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("penalChargescc: can't get loanStatus" + response.Message)
	}
	status := string(response.Payload)
	if status != "overdue" {
		return shim.Error("penalCharges.cc: " + "loan status for loanID " + args[3] + " is not overdue")
	}

	//TXN Amt must be > Zero
	if amt <= 0 {
		return shim.Error("penalCharges.cc: " + "Transaction Amount in Penal Interest Collectionis less than or equal to zero")
	}

	//####################################################################################################################

	//#####################################################################################################################
	//Calling for updating Business Main_Wallet
	//####################################################################################################################

	cAmtString := "0"
	dAmtString := args[5]

	walletID, err := getWalletID(stub, "businesscc", args[7], "main")
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest CollectionBusiness Main WalletID " + err.Error())
	}

	openBalance, err := getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest CollectionBusiness Main WalletValue " + err.Error())
	}
	openBalString := strconv.FormatInt(openBalance, 10)
	bal := openBalance - amt

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	txnBalString := strconv.FormatInt(bal, 10)
	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger

	argsList := []string{"1PIC", args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	txnResponse := putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Bank Main_Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, err = getWalletID(stub, "bankcc", args[6], "main")
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest CollectionBank Main WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Bank Main WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance + amt
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	argsList = []string{"2PIC", args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Business Charges O/s Wallet
	//####################################################################################################################

	cAmtString = "0"
	dAmtString = args[5]

	walletID, err = getWalletID(stub, "businesscc", args[7], "chargesOut")
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Business Charges O/s WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Business Charges O/s WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance - amt
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	argsList = []string{"3PIC", args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + txnResponse.Message)
	}

	//####################################################################################################################
	//Calling for updating Loan Charges Wallet
	//####################################################################################################################

	cAmtString = "0"
	dAmtString = args[5]

	walletID, err = getWalletID(stub, "loancc", args[3], "charges")
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Loan Charges Wallet WalletID " + err.Error())
	}

	openBalance, err = getWalletValue(stub, walletID)
	if err != nil {
		return shim.Error("penalCharges.cc: " + "Penal Interest Collection Loan Charges WalletValue " + err.Error())
	}
	openBalString = strconv.FormatInt(openBalance, 10)

	//amt, _ = strconv.ParseInt(args[5], 10, 64)

	bal = openBalance - amt
	txnBalString = strconv.FormatInt(bal, 10)

	response = walletUpdation(stub, walletID, bal)
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	argsList = []string{"4PIC", args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	txnResponse = putInTxnBal(stub, argsListStr)
	if txnResponse.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + txnResponse.Message)
	}

	//####################################################################################################################

	return shim.Success(nil)
}

func putInTxnBal(stub shim.ChaincodeStubInterface, argsListStr string) pb.Response {

	chaincodeArgs := toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("calling the txnbalcc chaincode from Interest Refund")
	response := stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("penalCharges.cc: " + response.Message)
	}
	fmt.Println(string(response.Payload))
	return shim.Success(nil)
}

func getWalletID(stub shim.ChaincodeStubInterface, ccName string, id string, walletType string) (string, error) {

	// STEP-1
	// using FromID, get a walletID from bank structure

	chaincodeArgs := toChaincodeArgs("getWalletID", id, walletType)
	response := stub.InvokeChaincode(ccName, chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return "0", errors.New(response.Message)
	}
	walletID := string(response.GetPayload())
	return walletID, nil

}

func getWalletValue(stub shim.ChaincodeStubInterface, walletID string) (int64, error) {

	walletArgs := toChaincodeArgs("getWallet", walletID)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return 0, errors.New(walletResponse.Message)
	}
	balString := string(walletResponse.Payload)
	balance, _ := strconv.ParseInt(balString, 10, 64)
	return balance, nil
}

func walletUpdation(stub shim.ChaincodeStubInterface, walletID string, amt int64) pb.Response {

	txnBalString := strconv.FormatInt(amt, 10)
	walletArgs := toChaincodeArgs("updateWallet", walletID, txnBalString)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return shim.Error(walletResponse.Message)
	}
	return shim.Success(nil)

}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start Penal Interest Collectionchaincode:", err)
	}
}
