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

type businessInfo struct {
	BusinessName              string
	BusinessAcNo              string
	BusinessLimit             int64
	BusinessWalletID          string //Hash
	BusinessLoanWalletID      string
	BusinessLiabilityWalletID string
	MaxROI                    float64
	MinROI                    float64
	NumberOfPrograms          int
	BusinessExposure          int64
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "putNewBusinessInfo" {
		return putNewBusinessInfo(stub, args)
	} else if function == "getBusinessInfo" {
		return getBusinessInfo(stub, args)
	} else if function == "getWalletID" {
		return getWalletID(stub, args)
	}
	return shim.Error("No function named " + function + " in Business")
}

func putNewBusinessInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 11 {
		return shim.Error("Invalid number of arguments. Needed 11 arguments")
	}

	businessLimitConv, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}

	maxROIconvertion, err := strconv.ParseFloat(args[7], 32)
	if err != nil {
		fmt.Printf("Invalid Maximum ROI: %s\n", args[7])
		return shim.Error(err.Error())
	}

	minROIconvertion, err := strconv.ParseFloat(args[8], 32)
	if err != nil {
		fmt.Printf("Invalid Minimum ROI: %s\n", args[8])
		return shim.Error(err.Error())
	}

	numOfPrograms, err := strconv.Atoi(args[9])
	if err != nil {
		fmt.Printf("Number of programs should be integer: %s\n", args[9])
	}

	businessExposureConv, err := strconv.ParseInt(args[10], 10, 64)
	if err != nil {
		fmt.Printf("Invalid business exposure: %s\n", args[10])
	}

	ifExists, err := stub.GetState(args[0])
	if ifExists != nil {
		fmt.Println(ifExists)
		return shim.Error("BusinessId " + args[0] + " exits. Cannot create new ID")
	}

	newInfo := &businessInfo{args[1], args[2], businessLimitConv, args[4], args[5], args[6], maxROIconvertion, minROIconvertion, numOfPrograms, businessExposureConv}
	newInfoBytes, _ := json.Marshal(newInfo)
	err = stub.PutState(args[0], newInfoBytes) // businessID = args[0]
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func getBusinessInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("Invalid number of arguments in getBusinessInfo (required:1) given:" + xLenStr)
	}

	parsedBusinessInfo := businessInfo{}
	businessIDvalue, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get the business information: " + err.Error())
	} else if businessIDvalue == nil {
		return shim.Error("No information is avalilable on this businessID " + args[0])
	}

	err = json.Unmarshal(businessIDvalue, &parsedBusinessInfo)
	if err != nil {
		return shim.Error("Unable to parse businessInfo into the structure " + err.Error())
	}
	jsonString := fmt.Sprintf("%+v", parsedBusinessInfo)
	fmt.Printf("Business Info: %s\n", jsonString)
	return shim.Success(nil)
}

func getWalletID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("Invalid number of arguments in getWalletId(business) (required:2) given:" + xLenStr)
	}

	parsedBusinessInfo := businessInfo{}
	businessIDvalue, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get the business information: " + err.Error())
	} else if businessIDvalue == nil {
		return shim.Error("No information is avalilable on this businessID " + args[0])
	}

	err = json.Unmarshal(businessIDvalue, &parsedBusinessInfo)
	if err != nil {
		return shim.Error("Unable to parse into the structure " + err.Error())
	}

	walletID := ""

	switch args[1] {
	case "main":
		walletID = parsedBusinessInfo.BusinessWalletID
	case "loan":
		walletID = parsedBusinessInfo.BusinessLoanWalletID
	case "liability":
		walletID = parsedBusinessInfo.BusinessLiabilityWalletID
	}

	return shim.Success([]byte(walletID))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("Error starting Business chaincode: %s\n", err)
	}

}
