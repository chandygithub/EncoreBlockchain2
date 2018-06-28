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

	if function == "putNewBusinessInfo" { //Inserting a New Business information
		return putNewBusinessInfo(stub, args)
	} else if function == "getBusinessInfo" { // To view a Business information
		return getBusinessInfo(stub,
			args)
	}
	return shim.Success(nil)
}

func putNewBusinessInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 11 {
		return shim.Error("Invalid number of arguments. Needed 11 arguments")
	}

	businessLimitConv, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}
	// CONVERTING STRING INTO 64 BIT 2 CONVERTION
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
		//fmt.Println("Required only one argument")
		return shim.Error("Required only one argument")
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
	jsonString := fmt.Sprintf("%+v", parsedBusinessInfo)
	return shim.Success([]byte(jsonString))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s\n", err)
	}

}
