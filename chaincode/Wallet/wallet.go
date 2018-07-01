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

type walletsInfo struct {
	Balance int64
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "newWallet" {
		return c.newWallet(stub, args)
	} else if function == "getWallet" {
		return c.getWallet(stub, args)
	} else if function == "updateWallet" {
		return c.updateWallet(stub, args)
	}
	return shim.Success(nil)

}

func (c *chainCode) newWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Invalid number of arguments, need ID and amt")
	}

	ifExists, err := stub.GetState(args[0])
	if ifExists != nil {
		return shim.Error("WalletId " + args[0] + " exits. Cannot create new ID")
	}

	bal64, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	bal := walletsInfo{bal64}
	balBytes, _ := json.Marshal(bal)
	err = stub.PutState(args[0], balBytes)
	return shim.Success([]byte("Successfully added to the wallet"))
}

func (c *chainCode) getWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Need only wallet ID to get wallet info")
	}
	balBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if balBytes == nil {
		return shim.Error("No data exists on this WalletId: " + args[0])
	}
	bal := walletsInfo{}
	err = json.Unmarshal(balBytes, &bal)
	if err != nil {
		return shim.Error(err.Error())
	}
	//balString := fmt.Sprintf("%+v", bal)
	//fmt.Printf("Wallet %s : %s\n", args[0], balString)

	balStr := strconv.FormatInt(bal.Balance, 10)

	return shim.Success([]byte(balStr))
}

func (c *chainCode) updateWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Invalid number of arguments")
	}
	balBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if balBytes == nil {
		return shim.Error("No data exists on this WalletId: " + args[0])
	}
	bal := walletsInfo{}
	err = json.Unmarshal(balBytes, &bal)
	if err != nil {
		return shim.Error(err.Error())
	}
	bal.Balance, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error("Error in Wallet updation parse int" + err.Error())
	}

	balBytes, _ = json.Marshal(bal)
	err = stub.PutState(args[0], balBytes)
	if err != nil {
		return shim.Error("Error in Wallet updation " + err.Error())
	}
	fmt.Printf("Balance for %s : %d\n", args[0], bal.Balance)
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to start the chaincode")
	}
}
