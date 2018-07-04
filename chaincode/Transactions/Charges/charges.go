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
	if function == "putTxnBalInfo" { //Inserting a New Business information
		return c.putTxnBalInfo(stub, args)
	} else {
		return shim.Error("incorrect function")
	}
}

func (c *chainCode) putTxnBalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	/*
	 * arg[0]	:	Key
	 * arg[1]	:	txnID	(0)
	 * arg[2]	:	date					//given	args[0]
	 * arg[3]	:	LoanID					//given	args[1]
	 * arg[4]	:	insID					//given args[2]
	 * arg[5]	:	bankID	=>	walletID	//given args[3]
	 *			:	bissID	=>	walletID	//given args[4]
	 * arg[6]	:	openBal
	 * arg[7]	:	txnType : charges		//given args[7]
	 * arg[8]	:	amt						//given args[5]
	 * arg[9]	:	cAmt (calc)
	 * arg[10]	:	dAmt (calc)
	 * arg[11]	:	txnBal
	 * arg[12]	:	by						//given args[6]
	 */
	if len(args) == 1 {
		args = strings.Split(args[0], ",")
	}
	if len(args) != 8 {
		return shim.Error("incorrcect number of arguments in chargescc")
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 													UPDATING WALLETS																///
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	/*
	 *	business load wallet increased
	 * 	bank revenue wallet incresed
	 *	bank asset wallet increased
	 */

	//####################################################################################################################
	//Calling for updating Business_loan_wallet
	//####################################################################################################################

	cAmtString := args[5]
	dAmtString := "0"

	walletID, openBalString, txnBalString, err := c.getWalletInfo(stub, args[4], "loan", "businesscc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error("could'nt get wallet:" + args[4] + " info: " + err.Error())
	}

	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList := []string{"1txnbal", args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr := strings.Join(argsList, ",")
	chaincodeArgs := toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("from charges calling the other chaincode: txnBalcc")
	response := stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}
	//fmt.Println(response.GetPayload())
	//successfully updated Bank's main wallet and written the txn thing to the ledger
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	//####################################################################################################################
	//Calling for updating Bank Charges_Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = c.getWalletInfo(stub, args[3], "charges", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error(err.Error())
	}

	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{"2txnbal", args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("from charges calling the other chaincode: txnBalcc")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}
	//fmt.Println(response.GetPayload())
	//successfully updated Bank's main wallet and written the txn thing to the ledger
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	//####################################################################################################################
	//Calling for updating Bank Asset_Wallet
	//####################################################################################################################

	cAmtString = args[5]
	dAmtString = "0"

	walletID, openBalString, txnBalString, err = c.getWalletInfo(stub, args[3], "asset", "bankcc", cAmtString, dAmtString)
	if err != nil {
		return shim.Error(err.Error())
	}

	// STEP-4 generate txn_balance_object and write it to the Txn_Bal_Ledger
	argsList = []string{"3txnbal", args[0], args[2], args[3], args[4], walletID, openBalString, args[1], args[5], cAmtString, dAmtString, txnBalString, args[8]}
	argsListStr = strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("putTxnBalInfo", argsListStr)
	fmt.Println("from charges calling the other chaincode: txnBalcc")
	response = stub.InvokeChaincode("txnbalcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}
	//fmt.Println(response.GetPayload())
	//successfully updated Bank's main wallet and written the txn thing to the ledger
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	/*
		//####################################################################################################################
		//Calling for Loan Balance Update
		//####################################################################################################################

		cAmtString = "0"
		dAmtString = args[5]
		argStrings := []string{"1loanbal", args[3], args[0], args[2], args[1], cAmtString, dAmtString} // 6 variables for updateLoanBalance
		argStr := strings.Join(argStrings, ",")
		chaincodeArgs = toChaincodeArgs("updateLoanBal", argStr)
		//sending to loanBalUp chaincode not loanBalance Chaincode
		response = stub.InvokeChaincode("loanbalcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("Error in updating Loan Balance: " + response.Message)
		}
		//walletID := string(response.GetPayload())
	*/
	return shim.Success(nil)
}

func (c *chainCode) getWalletInfo(stub shim.ChaincodeStubInterface, participantID string, walletType string, ccName string, cAmtStr string, dAmtStr string) (string, string, string, error) {

	// STEP-1
	// using FromID, get a walletID from bank structure
	// bankID = bankID

	chaincodeArgs := toChaincodeArgs("getWalletID", participantID, walletType)
	response := stub.InvokeChaincode(ccName, chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return "", "", "", errors.New(response.Message)
	}
	walletID := string(response.GetPayload())

	// STEP-2
	// getting Balance from walletID
	// walletFcn := "getWallet"
	walletArgs := toChaincodeArgs("getWallet", walletID)
	walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return "", "", "", errors.New(walletResponse.Message)
	}
	openBalString := string(walletResponse.Payload)

	openBal, err := strconv.ParseInt(openBalString, 10, 64)
	if err != nil {
		return "", "", "", errors.New("Error in converting the openBalance")
	}
	cAmt, err := strconv.ParseInt(cAmtStr, 10, 64)
	if err != nil {
		return "", "", "", errors.New("Error in converting the cAmt")
	}
	dAmt, err := strconv.ParseInt(dAmtStr, 10, 64)
	if err != nil {
		return "", "", "", errors.New("Error in converting the dAmt")
	}

	txnBal := openBal - dAmt + cAmt
	txnBalString := strconv.FormatInt(txnBal, 10)

	// STEP-3
	// update wallet of ID walletID here, and write it to the wallet_ledger
	// walletFcn := "updateWallet"
	walletArgs = toChaincodeArgs("updateWallet", walletID, txnBalString)
	walletResponse = stub.InvokeChaincode("walletcc", walletArgs, "myc")
	if walletResponse.Status != shim.OK {
		return "", "", "", errors.New("Not shim ok: " + walletResponse.Message)
	}

	return walletID, openBalString, txnBalString, nil
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start the chaincode")
	}
}
