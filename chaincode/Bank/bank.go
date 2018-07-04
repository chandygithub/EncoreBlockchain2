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
	//bankID			  string
	BankName              string `json:"BankName"`
	BankBranch            string `json:"BankBranch"`
	Bankcode              string `json:"BankCode"`
	BankWalletID          string `json:"BankMainWalletID"`
	BankAssetWalletID     string `json:"BankAssetWalletID"`
	BankChargesWalletID   string `json:"BankChargesWalletID"`
	BankLiabilityWalletID string `json:"BankLiabilityWalletID"`
	TDSreceivableWalletID string `json:"TDSreceivableWalletID"`
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "writeBankInfo" {
		return c.writeBankInfo(stub, args)
	} else if function == "getBankInfo" {
		return c.getBankInfo(stub, args)
	} else if function == "getWalletID" {
		return c.getWalletID(stub, args)
	}
	return shim.Error("Invalid Smart Contract function name in bankcc.")

}

func (c *chainCode) writeBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 9 {
		return shim.Error("Invalid number of arguments for writing bank info")
	}
	//args[0] -> bankID
	var bank = bankInfo{args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]}

	//check if it exits alreadty
	ifExists, err := stub.GetState(args[0])
	if ifExists != nil {
		fmt.Println(ifExists)
		return shim.Error("BankId " + args[0] + " exits. Cannot create new ID")
	}

	//marshal the bank to put on the ledger
	bankBytes, err := json.Marshal(bank)
	if err != nil {
		return shim.Error("Unable to Marshal the json file " + err.Error())
	}
	err = stub.PutState(args[0], bankBytes)

	fmt.Println("Successfully Written Bank" + args[1] + "into the ledger")
	return shim.Success([]byte("Successfully Written Bank" + args[1] + "into the ledger"))
}

func (c *chainCode) getBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Requird only one field for getting bank info")
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
	return shim.Success([]byte(x))
}

func (c *chainCode) getWalletID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
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

	var walletID string

	switch args[1] {
	case "main":
		walletID = bank.BankWalletID
	case "asset":
		walletID = bank.BankAssetWalletID
	case "charges":
		walletID = bank.BankChargesWalletID
	case "liability":
		walletID = bank.BankLiabilityWalletID
	case "tdsReceivable":
		walletID = bank.TDSreceivableWalletID
	}
	return shim.Success([]byte(walletID))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s\n", err)
	}
}
