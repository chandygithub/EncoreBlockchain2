package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type bankInfo struct {
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
	}
	return shim.Success([]byte("All done"))

}

func writeBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 9 {
		return shim.Error("Invalid number of arguments")
	}
	//args[0] -> bankID
	bank := bankInfo{args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]}
	bankBytes, err := json.Marshal(bank)
	if err != nil {
		return shim.Error("Unable to Marshal the json file " + err.Error())
	}
	err = stub.PutState(args[0], bankBytes)

	return shim.Success([]byte("Succefully written into the ledger"))
}

func getBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Requird only one field")
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
	return shim.Success([]byte(x))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s\n", err)
	}

}
