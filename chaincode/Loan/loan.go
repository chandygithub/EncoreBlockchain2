package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type loanInfo struct {
	InstNum                     string    `json:"InstrumentNo"`          //[1]//Instrument Number
	ExposureBusinessID          string    `json:"ExposureBusinessID"`    //[2]//buyer for now
	ProgramID                   string    `json:"ProgramID"`             //[3]
	SanctionAmt                 int64     `json:"SanctionAmt"`           //[4]
	SanctionDate                time.Time `json:"SanctionDate"`          //auto generated as created
	SanctionAuthority           string    `json:"SanctionAuthority"`     //[5]
	ROI                         float64   `json:"ROI"`                   //[6]
	DueDate                     time.Time `json:"DueDate"`               //[7]
	ValueDate                   time.Time `json:"ValueDate"`             //[8]//with time
	LoanStatus                  string    `json:"LoanStatus"`            //[]
	LoanDisbursedWalletID       string    `json:"DisbursementWallet"`    //[9]
	LoanChargesWalletID         string    `json:"ChargesWallet"`         //[10]
	LoanAccruedInterestWalletID string    `json:"AccruedInterestWallet"` //[11]
	BuyerBusinessID             string    `json:"BuyerID"`               //[12]
	SellerBusinessID            string    `json:"SellerID"`              //[13]
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

	if function == "newLoanInfo" {
		//Creates a new Loan Data
		return newLoanInfo(stub, args)
	} else if function == "getLoanInfo" {
		//Retrieves the existing data
		return getLoanInfo(stub, args)
	} else if function == "updateLoanInfo" {
		//Updates variables for loan structure
		return updateLoanInfo(stub, args)
	} else if function == "loanIDexists" {
		//Checks the existence of loan ID
		return loanIDexists(stub, args[0])
	} else if function == "getLoanStatus" {
		//Returns the Loan status
		return getLoanStatus(stub, args[0])
	} else if function == "getLoanSancAmt" {
		//Returns the Sanc Amt
		return getLoanSancAmt(stub, args[0])
	} else if function == "getWalletID" {
		//Returns the walletID for the required wallet type
		return getWalletID(stub, args)
	} else if function == "getSellerID" {
		//Returns the Seller Id
		return getSellerID(stub, args[0])
	} else if function == "getBuyerID" {
		return getBuyerID(stub, args[0])
	}
	return shim.Error("loancc: " + "No function named " + function + " in Loanssssssssssss")
}

func newLoanInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 14 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("loancc: " + "Invalid number of arguments in newLoanInfo(loan) (required:14) given: " + xLenStr)
	}

	//Checking existence of loanID
	println("Checking existence of loanID")
	response := loanIDexists(stub, args[0])
	if response.Status != shim.OK {
		return shim.Error("loancc: " + response.Message)
	}

	//Checking existence of ExposureBusinessID
	println("Checking existence of ExposureBusinessID")
	chaincodeArgs := toChaincodeArgs("bisIDexists", args[2])
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("loancc: " + "ExposureBusinessID " + args[2] + " does not exits")
	}

	//Checking if Instrument ID is Instrument Ref. No.
	println("Checking if Instrument ID is Instrument Ref. No.")
	chaincodeArgs = toChaincodeArgs("getInstrument", args[1], args[13])
	response = stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("loancc: " + "Instrument refrence no " + args[1] + " does not exits")
	}

	// getting the sanction amount from the instrument
	chaincodeArgs = toChaincodeArgs("getInstrumentAmt", args[1], args[13])
	response = stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("loancc: " + response.Message)
	}
	instAmtStr := string(response.Payload)
	fmt.Println("instAmtStr: " + instAmtStr)
	instAmt, err := strconv.ParseInt(instAmtStr, 10, 64)
	if err != nil {
		return shim.Error("loancc: " + "Unable to parse instAmt(loan): " + err.Error())
	}

	//Getting the discount percentage
	println("Getting the discount percentage")
	chaincodeArgs = toChaincodeArgs("discountPercentage", args[3], args[2])
	response = stub.InvokeChaincode("pprcc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("loancc: " + "PprId " + args[8] + " does not exits")
	}

	discountPercentStr := string(response.Payload)
	discountPercent, _ := strconv.ParseInt(discountPercentStr, 10, 64)
	amt := instAmt - ((discountPercent * instAmt) / 100)

	//SanctionAmt -> sAmt
	println("SanctionAmt -> sAmt")
	sAmt, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	if sAmt > amt && sAmt == 0 {
		return shim.Error("loancc: " + "Sanction amount exceeds the required value or it is zero : " + args[4])
	}

	//SanctionDate ->sDate
	println("SanctionDate ->sDate")
	sDate := time.Now()

	roi, err := strconv.ParseFloat(args[6], 32)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	//Parsing into date for storage but hh:mm:ss will also be stored as
	println("Parsing into date for storage")
	//00:00:00 .000Z with the date
	//DueDate -> dDate
	dDate, err := time.Parse("02/01/2006", args[7])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}
	if dDate.Weekday().String() == "Sunday" {
		fmt.Println("Since the due date falls on sunday, due date is extended to Monday(loan) : ", dDate.AddDate(0, 0, 1))
	}
	dDate = dDate.AddDate(0, 0, 1)

	//Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing
	println("Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing")
	vDateStr := args[8][:10]
	vTime := args[8][11:]
	vStr := vDateStr + "T" + vTime

	//ValueDate ->vDate
	println("ValueDate ->vDate")
	vDate, err := time.Parse("02/01/2006T15:04:05", vStr)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	hash := sha256.New()
	println("Hashing wallets")
	// Hashing LoanDisbursedWalletID
	LoanDisbursedWalletStr := args[9] + "LoanDisbursedWallet"
	hash.Write([]byte(LoanDisbursedWalletStr))
	md := hash.Sum(nil)
	LoanDisbursedWalletIDsha := hex.EncodeToString(md)
	createWallet(stub, LoanDisbursedWalletIDsha, args[9])

	// Hashing LoanChargesWalletID
	LoanChargesWalletStr := args[10] + "LoanChargesWallet"
	hash.Write([]byte(LoanChargesWalletStr))
	md = hash.Sum(nil)
	LoanChargesWalletIDsha := hex.EncodeToString(md)
	createWallet(stub, LoanChargesWalletIDsha, args[10])

	// Hashing LoanAccruedInterestWalletID
	LoanAccruedInterestWalletStr := args[11] + "LoanAccruedInterestWallet"
	hash.Write([]byte(LoanAccruedInterestWalletStr))
	md = hash.Sum(nil)
	LoanAccruedInterestWalletIDsha := hex.EncodeToString(md)
	createWallet(stub, LoanAccruedInterestWalletIDsha, args[11])

	//Checking existence of BuyerBusinessID
	println("Checking existence of BuyerBusinessID")
	chaincodeArgs = toChaincodeArgs("bisIDexists", args[13])
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("loancc: " + "BuyerBusinessID " + args[13] + " does not exits")
	}

	//Checking existence of SellerBusinessID
	println("Checking existence of SellerBusinessID")
	chaincodeArgs = toChaincodeArgs("bisIDexists", args[13])
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("loancc: " + "SellerBusinessID " + args[13] + " does not exits")
	}

	println("marshalling loaninfo")
	loan := loanInfo{args[1], args[2], args[3], sAmt, sDate, args[6], roi, dDate, vDate, "sanctioned", LoanDisbursedWalletIDsha, LoanChargesWalletIDsha, LoanAccruedInterestWalletIDsha, args[12], args[13]}
	loanBytes, err := json.Marshal(loan)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}
	stub.PutState(args[0], loanBytes)

	println("changing inst status")
	//argsList := []string{args[1], args[13], "sanctioned"}
	//argsListStr := strings.Join(argsList, ",")
	chaincodeArgs = toChaincodeArgs("updateInstrumentStatus", args[1], args[13], "sanctioned")
	response = stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("loancc: " + response.Message)
	}
	return shim.Success([]byte("Successfully added loan info into ledger"))
}

func createWallet(stub shim.ChaincodeStubInterface, walletID string, amt string) pb.Response {
	chaincodeArgs := toChaincodeArgs("newWallet", walletID, amt)
	response := stub.InvokeChaincode("walletcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("loancc: " + "Unable to create new wallet from business")
	}
	return shim.Success([]byte("created new wallet from business"))
}

func loanIDexists(stub shim.ChaincodeStubInterface, loanID string) pb.Response {
	ifExists, _ := stub.GetState(loanID)
	if ifExists != nil {
		fmt.Println(ifExists)
		return shim.Error("loancc: " + "LoanId " + loanID + " exits. Cannot create new ID")
	}
	return shim.Success(nil)
}

func getLoanStatus(stub shim.ChaincodeStubInterface, loanID string) pb.Response {
	fmt.Println("loancc: inside getLoanStatus")
	loanBytes, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + loanID)
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "Error unmarshiling in loanstatus(loan):" + err.Error())
	}

	return shim.Success([]byte(loan.LoanStatus))
}

func getLoanSancAmt(stub shim.ChaincodeStubInterface, loanID string) pb.Response {
	fmt.Println("loancc: inside getLoanSancAmt")
	loanBytes, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + loanID)
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "Error unmarshiling in loanstatus(loan):" + err.Error())
	}
	fmt.Println(loan.SanctionAmt)
	sancAmtString := strconv.FormatInt(loan.SanctionAmt, 10)
	fmt.Println(sancAmtString)
	return shim.Success([]byte(sancAmtString))
}

func getWalletID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "Unable to parse into loan the structure (loanWalletValues)" + err.Error())
	}

	walletID := ""

	switch args[1] {
	case "accrued":
		walletID = loan.LoanAccruedInterestWalletID
	case "charges":
		walletID = loan.LoanChargesWalletID
	case "disbursed":
		walletID = loan.LoanDisbursedWalletID
	default:
		return shim.Error("loancc: " + "There is no wallet of this type in Loan :" + args[1])
	}

	return shim.Success([]byte(walletID))
}

func getLoanInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("loancc: " + "Invalid number of arguments in getLoanInfo (required:1) given:" + xLenStr)

	}

	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	loanString := fmt.Sprintf("%+v", loan)
	fmt.Printf("Loan Info:%s\n ", loanString)

	return shim.Success(loanBytes)
}

func getSellerID(stub shim.ChaincodeStubInterface, loanID string) pb.Response {

	loanBytes, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID (getSellerID): " + loanID)
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	return shim.Success([]byte(loan.SellerBusinessID))
}

func getBuyerID(stub shim.ChaincodeStubInterface, loanID string) pb.Response {

	loanBytes, err := stub.GetState(loanID)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID (getSellerID): " + loanID)
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	}

	return shim.Success([]byte(loan.BuyerBusinessID))
}

func updateLoanInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/*
		Updating the variables for loan structure
	*/
	//args = strings.Split(args[0], ",")
	loanBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("loancc: " + err.Error())
	} else if loanBytes == nil {
		return shim.Error("loancc: " + "No data exists on this loanID: " + args[0])
	}

	loan := loanInfo{}
	err = json.Unmarshal(loanBytes, &loan)
	if err != nil {
		return shim.Error("loancc: " + "error in unmarshiling loan: in updateLoanInfo" + err.Error())
	}

	// To change the LoanStatus from "sanction" to "disbursed"
	if args[2] == "disbursement" {
		if (loan.LoanStatus != "sanctioned") && (loan.LoanStatus != "part disbursed") {
			return shim.Error("loancc: " + "Loan is not Sanctioned, so cannot be disbursed/ part Disbursed : " + loan.LoanStatus)
		}
		//Updating Loan status for disbursement
		loan.LoanStatus = args[1]
		loanBytes, _ := json.Marshal(loan)
		err = stub.PutState(args[0], loanBytes)
		if err != nil {
			return shim.Error("loancc: " + "Error in loan updation " + err.Error())
		}

		//Calling instrument chaincode to update the status
		//argsList := []string{loan.InstNum, loan.SellerBusinessID, "disbursed"}
		//argsListStr := strings.Join(argsList, ",")
		chaincodeArgs := toChaincodeArgs("updateInstrumentStatus", loan.InstNum, loan.SellerBusinessID, "disbursed")
		response := stub.InvokeChaincode("instrumentcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return shim.Error("loancc: " + response.Message)
		}
		return shim.Success([]byte("sanction updated succesfully"))

	} else if (args[1] == "repayment") && ((args[2] == "collected") || (args[2] == "part collected")) {
		if loan.LoanStatus != "disbursed" {
			return shim.Error("loancc: " + "Loan is not Sanctioned, so cannot be disbursed")
		}
		//Updating Loan status for repayment
		loan.LoanStatus = args[2]
		loanBytes, _ = json.Marshal(loan)
		err = stub.PutState(args[0], loanBytes)
		if err != nil {
			return shim.Error("loancc: " + "Error in loan status updation " + err.Error())
		}

		return shim.Success([]byte("Successfully updated loan status with data from repayment"))
	}
	return shim.Error("loancc: " + "Invalid info for update loan")
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("loancc: "+"Error starting Loan chaincode: %s\n", err)
	}
}
