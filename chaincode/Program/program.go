package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type programInfo struct {
	ProgramName        string    `json:"ProgramName"`
	ProgramAnchor      string    `json:"ProgramAnchor"`
	ProgramType        string    `json:"ProgramType"`
	ProgramStartDate   time.Time `json:"ProgramStartDate"`
	ProgramEndDate     time.Time `json:"ProgramEndDate"`
	ProgramLimit       int64     `json:"ProgramLimit"`
	ProgramROI         float64   `json:"ProgramROI"`
	ProgramExposure    string    `json:"ProgramExposure"`
	DiscountPercentage float64   `json:"DiscountPercentage"`
	DiscountPeriod     int       `json:"DiscountPeriod"`
	SanctionAuthority  string    `json:"SanctionAuthority"`
	SanctionDate       time.Time `json:"SanctionDate"`
	RepaymentAcNum     string    `json:"RepaymentAccountNumber"`
	RepaymentWalletID  string    `json:"RepaymentWalletID"`
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "writeProgram" {
		return writeProgram(stub, args)
	} else if function == "getProgram" {
		return getProgram(stub, args)
	}
	return shim.Success(nil)
}

func writeProgram(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 13 {
		return shim.Error("Invalid number of arguments")
	}

	//args[0] -> programID ; Key for the structure, must be passed by the user

	pTypes := map[string]bool{
		"ar": true,
		"ap": true,
		"df": true,
	}

	//Checking whether the given argument is a valid type
	pTypeLower := strings.ToLower(args[3])
	if !pTypes[pTypeLower] {
		return shim.Error("Invalid program type" + pTypeLower)
	}

	//ProgramStartDate -> pSDate
	pSDate, err := time.Parse("02/01/2006", args[4])
	if err != nil {
		return shim.Error(err.Error())
	}

	//ProgramEndDate -> pEDate
	pEDate, err := time.Parse("02/01/2006", args[5])
	if err != nil {
		return shim.Error(err.Error())
	}

	pLimit, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return shim.Error("Invalid Program limit " + args[6])
	}

	pROI, err := strconv.ParseFloat(args[7], 32)
	if err != nil {
		return shim.Error("Invalid Rate of Interest")
	}

	pExposure := map[string]bool{
		"buyer":  true,
		"seller": true,
	}

	pExposureLower := strings.ToLower(args[8])

	if !pExposure[pExposureLower] {
		return shim.Error("Invalid Program Exposure " + pExposureLower)
	}

	dPercentage, err := strconv.ParseFloat(args[9], 32)
	if err != nil {
		return shim.Error("Invalid discount percentage")
	}

	dPeriod, err := strconv.Atoi(args[10])
	if err != nil {
		return shim.Error("Invalid discount period")
	}

	//SanctionDate -> sDate
	sDate, err := time.Parse("02/01/2006", args[12])
	if err != nil {
		return shim.Error(err.Error())
	}

	pInfo := programInfo{args[1], args[2], pTypeLower, pSDate, pEDate, pLimit, pROI, pExposureLower, dPercentage, dPeriod, args[11], sDate, args[13], args[14]}
	programInfoBytes, _ := json.Marshal(pInfo)
	err = stub.PutState(args[0], programInfoBytes)
	return shim.Success(nil)
}

func getProgram(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Invalid number of arguments")
	}

	pInfo := programInfo{}
	pInfoBytes, err := stub.GetState(args[0])

	if err != nil {
		return shim.Error(err.Error())
	} else if pInfoBytes == nil {
		return shim.Error("No information on this programID: " + args[0])
	}

	err = json.Unmarshal(pInfoBytes, &pInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	printProgramInfo := fmt.Sprintf("%+v", pInfo)

	return shim.Success([]byte(printProgramInfo))

}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("Unable to initiate the chaincode")
	}
}
