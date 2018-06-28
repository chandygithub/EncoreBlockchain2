package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type C1Info struct{}

func (c *C1Info) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Successsssssssssss")
	return shim.Success(nil)
}

func (c *C1Info) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("inside the invoke function")
	if function == "printThis" {
		fmt.Println("inside the if statement")
		/*x := c.PrintThis(stub, args)
		fmt.Println("After the function call")*/
		return c.PrintThis(stub, args) //x

	}
	fmt.Println("outside if")

	return shim.Success(nil)
}

func (c *C1Info) PrintThis(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Printf("Printing %s\n", args[0])
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(C1Info))
	if err != nil {
		fmt.Println("Unable to start the chaincode")
	}
	//x := a.PrintThis(shim.ChaincodeStubInterface, "cool")
}
