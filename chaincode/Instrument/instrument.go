package main

import (
	"crypto/sha256"
	"encoding/hex"
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

type instrumentInfo struct {
	//Instrument ID for storing is auto generated
	InstrumentRefNo string    `json:"RefNo"`         //[0]
	InstrumentDate  time.Time `json:"Date"`          //[1]
	SellBusinessID  string    `json:"SellerID"`      //[2]
	BuyBusinsessID  string    `json:"BuyerID"`       //[3]
	InsAmount       string    `json:"Amount"`        //[4]// use int64 for convertion
	InsStatus       string    `json:"Status"`        // not required
	InsDueDate      time.Time `json:"DueDate"`       //[5]
	ProgramID       string    `json:"ProgramID"`     //[6]
	PPRid           string    `json:"PPRID"`         //[7]
	UploadBatchNo   string    `json:"UploadBatchNo"` //[8]
	ValueDate       time.Time `json:"ValueDate"`     //[9]
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	indexName := "InstrumentRefNo~SellBusinessID~InsAmount"
	inst := instrumentInfo{}

	refNoSellIDkey, err := stub.CreateCompositeKey(indexName, []string{inst.InstrumentRefNo, inst.SellBusinessID, inst.InsAmount})
	if err != nil {
		return shim.Error("instrumetcc: " + "Composite key InstrumentRefNo~SellBusinessID~InsAmount can not be created (instrument)")
	}
	value := []byte{0x00}
	stub.PutState(refNoSellIDkey, value)
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "enterInstrument" {
		//Used to enter new instrument data
		return enterInstrument(stub, args)
	} else if function == "getInstrument" {
		//used to retrieve the instrument data
		return getInstrument(stub, args)
	} else if function == "updateInstrumentStatus" { //instrumentStatus
		//Updates instrument status accordingly
		return updateInstrumentStatus(stub, args)
	} else if function == "getInstrumentAmt" {
		return getInstrumentAmt(stub, args)
	}

	return shim.Error("instrumetcc: " + "No function named " + function + " in Instrumentsssss")

}

func enterInstrument(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("instrumetcc: " + "Invalid number of arguments in enterInstrument (required:10) given:" + xLenStr)

	}

	// Checking existence of Instrument Reference No. – Supplier ID pair
	refNoSellIDiterator, _ := stub.GetStateByPartialCompositeKey("InstrumentRefNo~SellBusinessID~InsAmount", []string{args[0], args[2]})
	refNoSellIDdata, _ := refNoSellIDiterator.Next()
	if refNoSellIDdata != nil {
		return shim.Error("instrumetcc: " + "Instrument Reference No. – Supplier ID pair already exists")
	}
	defer refNoSellIDiterator.Close()

	//Checking existence of ProgramID
	chaincodeArgs := toChaincodeArgs("programIDexists", args[6])
	response := stub.InvokeChaincode("programcc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("instrumetcc: " + "ProgramId " + args[6] + " does not exits")
	}

	//Checking existence of pprID
	chaincodeArgs = toChaincodeArgs("pprIDexists", args[7])
	response = stub.InvokeChaincode("pprcc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("instrumetcc: " + "PprId " + args[7] + " does not exits")
	}

	//Checking existence of SellerBusinessID
	chaincodeArgs = toChaincodeArgs("bisIDexists", args[2])
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("instrumetcc: " + "BusinessId " + args[2] + " does not exits")
	}

	//Checking existence of BuyerBusinessID
	chaincodeArgs = toChaincodeArgs("bisIDexists", args[3])
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("instrumetcc: " + "BusinessId " + args[3] + " does not exits")
	}

	//InstrumentDate -> instDate
	instDate, err := time.Parse("02/01/2006", args[1])
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	}

	_, err = strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	}

	//InsDueDate -> insDate
	insDueDate, err := time.Parse("02/01/2006", args[5])
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	}
	if insDueDate.Weekday().String() == "Sunday" {
		fmt.Println("Since the due date falls on sunday, due date is extended to Monday(instrument) : ", insDueDate.AddDate(0, 0, 1))
	}
	insDueDate = insDueDate.AddDate(0, 0, 1)
	//Converting the incoming date from Dd/mm/yy:hh:mm:ss to Dd/mm/yyThh:mm:ss for parsing
	vString := args[9][:10] + "T" + args[9][11:] //removing the ":" part from the string

	//ValueDate -> vDate
	vDate, err := time.Parse("02/01/2006T15:04:05", vString)
	if err != nil {
		return shim.Error("instrumetcc: " + "error in parsing the date and time (instrument)" + err.Error())
	}

	// Hashing for key to store in ledger
	hash := sha256.New()
	instID := strings.ToLower(args[0] + args[2])
	hash.Write([]byte(instID))
	md := hash.Sum(nil)
	instIDsha := hex.EncodeToString(md)

	inst := instrumentInfo{args[0], instDate, args[2], args[3], args[4], "open", insDueDate, args[6], args[7], args[8], vDate}
	instBytes, err := json.Marshal(inst)
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	}
	stub.PutState(instIDsha, instBytes)

	return shim.Success([]byte("Successfully added instrument to the ledger"))
}

func updateInstrumentStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/*
		args[0] -> instrument reference number
		args[1] -> seller ID
		args[2] -> status
	*/
	key := strings.ToLower(args[0] + args[1])
	hash := sha256.New()
	instID := strings.ToLower(key)
	hash.Write([]byte(instID))
	md := hash.Sum(nil)
	instIDsha := hex.EncodeToString(md)

	instBytes, err := stub.GetState(instIDsha)
	if err != nil {
		return shim.Error("instrumetcc: " + "Unable to fetch instrument info for status updation")
	}
	inst := instrumentInfo{}
	err = json.Unmarshal(instBytes, &inst)
	if err != nil {
		return shim.Error("instrumetcc: " + "Error in unmarshaling the instrument (updateInsStatus)")
	}
	/*
	 updated sequentially Open > Sanctioned > Overdue>Settled or Open > Sanctioned > Settled
	*/
	if (args[2] == "sanctioned") && (inst.InsStatus != "open") {
		return shim.Error("instrumetcc: " + "Instrument status cannot be sanctioned as it is not open")
	} else if (args[2] == "overdue") && (inst.InsStatus != "sanctioned") {
		return shim.Error("instrumetcc: " + "Instrument status cannot be overdue as it is not sanctioned")
	} else if (args[2] == "settled") && ((inst.InsStatus != "overdue") && (inst.InsStatus != "sanctioned")) {
		return shim.Error("instrumetcc: " + "Instrument status cannot be settled as it is not overdue or sanctioned")
	}
	inst.InsStatus = args[2]
	instBytes, _ = json.Marshal(inst)
	stub.PutState(key, instBytes)

	return shim.Success([]byte("Instrument status updated successfully"))

}

func getInstrument(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("instrumetcc: " + "Invalid number of arguments in getInstrument (required:2) given:" + xLenStr)

	}
	/*
		args[0] -> InstrumentRefNo
		args[1] -> SellBusinessID
	*/
	hash := sha256.New()
	instID := strings.ToLower(args[0] + args[1])
	hash.Write([]byte(instID))
	md := hash.Sum(nil)
	instIDsha := hex.EncodeToString(md)

	insBytes, err := stub.GetState(instIDsha)
	if err != nil {
		return shim.Error("instrumetcc: " + err.Error())
	} else if insBytes == nil {
		return shim.Error("instrumetcc: " + "No data exists on this InstrumentID: " + args[0])
	}

	//ins := instrumentInfo{}
	//err = json.Unmarshal(insBytes, &ins)
	//insString := fmt.Sprintf("%+v", ins)
	return shim.Success(nil)
}

func getInstrumentAmt(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("invalid no. of arguments (requiered 2)")
	}

	/*fmt.Println("before refNoSellIDiterator")
	refNoSellIDiterator, err := stub.GetStateByPartialCompositeKey("InstrumentRefNo~SellBusinessID~InsAmount", args)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("before refNoSellIDiterator")
	refNoSellIDdata, err := refNoSellIDiterator.Next()
	if err != nil {
		return shim.Error(err.Error())
	}
	//fmt.Println(refNoSellIDdata.Key)
	_, insAmt, err := stub.SplitCompositeKey(refNoSellIDdata.Key)
	if err != nil {
		return shim.Error("Error splitting the composite key " + err.Error())
	}
	fmt.Println(insAmt)
	*/
	hash := sha256.New()
	instID := strings.ToLower(args[0] + args[1])
	hash.Write([]byte(instID))
	md := hash.Sum(nil)
	instIDsha := hex.EncodeToString(md)

	insBytes, err := stub.GetState(instIDsha)
	if err != nil {
		return shim.Error(err.Error())
	} else if insBytes == nil {
		return shim.Error("No data exists on this InstrumentID: " + instIDsha)
	}

	ins := instrumentInfo{}
	err = json.Unmarshal(insBytes, &ins)
	if err != nil {
		return shim.Error("Cannot unmarshal json")
	}

	return shim.Success([]byte(ins.InsAmount))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("instrumetcc: "+"Error starting Instrument chaincode: %s\n", err)
	}
}
