package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type bankInfo struct {
	BankName              string `json:"BankName"`
	BankBranch            string `json:"BankBranch"`
	Bankcode              string `json:"BankCode"`
	BankWalletID          string `json:"MainWallet"`      //will take the values for the respective wallet from the user
	BankAssetWalletID     string `json:"AssetWallet"`     //will take the values for the respective wallet from the user
	BankChargesWalletID   string `json:"ChargesWallet"`   //will take the values for the respective wallet from the user
	BankLiabilityWalletID string `json:"LiabilityWallet"` //will take the values for the respective wallet from the user
	TDSreceivableWalletID string `json:"TDSWallet"`       //will take the values for the respective wallet from the user
}

// toChaincodeArgs returns byte arrau of string of arguments, so it can be passed to other chaincodes
func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	bank := bankInfo{}
	indexName := "Bankcode~BankBranch"
	codeBranchKey, err := stub.CreateCompositeKey(indexName, []string{bank.Bankcode, bank.BankBranch})
	if err != nil {
		return shim.Error("Unable to create composite key Bankcode~BankBranch in bankcc")
	}
	value := []byte{0x00}
	stub.PutState(codeBranchKey, value)
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "writeBankInfo" {
		//Creates a new Bank Information
		return writeBankInfo(stub, args)
	} else if function == "getBankInfo" {
		//Retrieves the Bank information
		return getBankInfo(stub, args)
	} else if function == "getWalletID" {
		//Returns the walletID for the required wallet type
		return getWalletID(stub, args)
	} else if function == "bankIDexists" {
		//To check the BankId existence
		return bankIDexists(stub, args[0])
	}
	return shim.Error("No function named " + function + " in Banksssssssss")

}

// To write the bank info into the ledger
func writeBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//Checking argument length
	if len(args) != 9 {
		xLenStr := strconv.Itoa(len(args)) //needed?!
		return shim.Error("Invalid number of arguments in writeBankInfo (required:9) given:" + xLenStr)
	}

	//Checking Bank ID existence
	response := bankIDexists(stub, args[0])
	if response.Status != shim.OK {
		return shim.Error("bankcc: " + response.Message)
	}

	//Checking existence of Bank code (!*!)
	codeBranchIterator, err := stub.GetStateByPartialCompositeKey("Bankcode~BankBranch", []string{args[3]})
	codeBranchData, err := codeBranchIterator.Next()
	if codeBranchData != nil {
		return shim.Error("bankcc : " + "Bank code already exist: " + args[3])
	}
	defer codeBranchIterator.Close()

	hash := sha256.New()

	//TODO:	check for collisions

	// Hashing bankWalletId
	BankWalletStr := args[3] + "BankWallet"
	hash.Write([]byte(BankWalletStr))
	md := hash.Sum(nil)
	BankWalletIDsha := hex.EncodeToString(md)
	createWallet(stub, BankWalletIDsha, "1000")

	// Hashing bankAssetWalletId
	BankAssetWalletStr := args[3] + "BankAssetWallet"
	hash.Write([]byte(BankAssetWalletStr))
	md = hash.Sum(nil)
	BankAssetWalletIDsha := hex.EncodeToString(md)
	createWallet(stub, BankAssetWalletIDsha, "1000")

	// Hashing BankChargesWalletID
	BankChargesWalletStr := args[3] + "BankChargesWallet"
	hash.Write([]byte(BankChargesWalletStr))
	md = hash.Sum(nil)
	BankChargesWalletIDsha := hex.EncodeToString(md)
	createWallet(stub, BankChargesWalletIDsha, "1000")

	// Hashing BankLiabilityWalletID
	BankLiabilityWalletStr := args[3] + "BankLiabilityWallet"
	hash.Write([]byte(BankLiabilityWalletStr))
	md = hash.Sum(nil)
	BankLiabilityWalletIDsha := hex.EncodeToString(md)
	createWallet(stub, BankLiabilityWalletIDsha, "1000")

	// Hashing TDSreceivableWalletID
	TDSreceivableWalletStr := args[3] + "TDSreceivableWallet"
	hash.Write([]byte(TDSreceivableWalletStr))
	md = hash.Sum(nil)
	TDSreceivableWalletIDsha := hex.EncodeToString(md)
	createWallet(stub, TDSreceivableWalletIDsha, "1000")

	//args[0] -> bankID | creating a bank struct obj and writing it to the ledger
	bank := bankInfo{args[1], args[2], args[3], BankWalletIDsha, BankAssetWalletIDsha, BankChargesWalletIDsha, BankLiabilityWalletIDsha, TDSreceivableWalletIDsha}
	bankBytes, err := json.Marshal(bank)
	if err != nil {
		return shim.Error("Unable to Marshal the json file " + err.Error())
	}

	err = stub.PutState(args[0], bankBytes)
	if err != nil {
		return shim.Error("bankcc : " + err.Error())
	}

	return shim.Success([]byte("Succefully written into the ledger"))
}

func bankIDexists(stub shim.ChaincodeStubInterface, bankID string) pb.Response {
	ifExists, _ := stub.GetState(bankID)
	if ifExists != nil {
		fmt.Println(ifExists) //needed!?
		return shim.Error("bankcc : " + "BankId " + bankID + " exits. Cannot create new ID")
	}
	return shim.Success(nil)
}

func createWallet(stub shim.ChaincodeStubInterface, walletID string, amt string) pb.Response {
	//Calling wallet Chaincode to create new wallet
	chaincodeArgs := toChaincodeArgs("newWallet", walletID, amt)
	response := stub.InvokeChaincode("walletcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("bankcc : " + "Unable to create new wallet from bank")
	}
	return shim.Success([]byte("created new wallet from bank"))
}

func getBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("bankcc : " + "Invalid number of arguments in getBankInfo (required:1) given:" + xLenStr)
	}

	// getting bank info from the ledger, args[0] -> bankID
	bankInfoBytes, err := stub.GetState(args[0])

	if err != nil {
		return shim.Error("bankcc : " + "Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		return shim.Error("bankcc : " + "Data does not exist for " + args[0])
	}

	// unmarshalling bankInfoBytes into bank structure
	bank := bankInfo{}
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		return shim.Error("bankcc : " + "Uable to paser into the json format")
	}
	x := fmt.Sprintf("%+v", bank)
	fmt.Printf("BankInfo : %s\n", x)
	return shim.Success(nil)
}

func getWalletID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("bankcc : " + "Invalid number of arguments in getWalletID(bank) (required:2) given:" + xLenStr)
	}
	bankInfoBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("bankcc : " + "Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		return shim.Error("bankcc : " + "Data does not exist for " + args[0])
	}
	bank := bankInfo{}
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		return shim.Error("bankcc : " + "Uable to paser into the json format")
	}

	walletID := ""

	switch args[1] {
	case "main":
		walletID = bank.BankWalletID
	case "asset":
		walletID = bank.BankAssetWalletID
	case "charges":
		walletID = bank.BankChargesWalletID
	case "liability":
		walletID = bank.BankLiabilityWalletID
	case "tds":
		walletID = bank.TDSreceivableWalletID
	}

	return shim.Success([]byte(walletID))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("bankcc : "+"Error starting Bank chaincode: %s\n", err)
	}

}
