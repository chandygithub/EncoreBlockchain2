package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type bankInfo struct {
	//bankID			  string
	BankName              string
	BankBranch            string
	Bankcode              string
	BankWalletID          string
	BankAssetWalletID     string
	BankChargesWalletID   string
	BankLiabilityWalletID string
	TDSreceivableWalletID string
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "writeBankInfo" {
		return writeBankInfo(stub, args)
	} else if function == "getBankInfo" {
		return getBankInfo(stub, args)
	} else if function == "getWalletID" {
		return getWalletID(stub, args)
	}
	return shim.Error("No function named " + function + " in Bank")

}

func writeBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 9 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("Invalid number of arguments in writeBankInfo (required:9) given:" + xLenStr)
	}
	//args[0] -> bankID
	bank := bankInfo{args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]}
	bankBytes, err := json.Marshal(bank)
	if err != nil {
		return shim.Error("Unable to Marshal the json file " + err.Error())
	}

	ifExists, err := stub.GetState(args[0])
	if ifExists != nil {
		fmt.Println(ifExists)
		return shim.Error("BankId " + args[0] + " exits. Cannot create new ID")
	}

	err = stub.PutState(args[0], bankBytes)

	return shim.Success([]byte("Succefully written into the ledger"))
}

func getBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("Invalid number of arguments in getBankInfo (required:1) given:" + xLenStr)
	}

	bankInfoBytes, err := stub.GetState(args[0])

	if err != nil {
		return shim.Error("Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		return shim.Error("Data does not exist for " + args[0])
	}

	bank := bankInfo{}
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		return shim.Error("Uable to paser into the json format")
	}
	x := fmt.Sprintf("%+v", bank)
	fmt.Printf("BankInfo : %s\n", x)
	return shim.Success(nil)
}

func getWalletID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("Invalid number of arguments in getWalletID(bank) (required:2) given:" + xLenStr)
	}
	bankInfoBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		return shim.Error("Data does not exist for " + args[0])
	}
	bank := bankInfo{}
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		return shim.Error("Uable to paser into the json format")
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
	}
	return shim.Success([]byte(walletID))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("Error starting Bank chaincode: %s\n", err)
	}

}
