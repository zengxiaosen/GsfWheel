package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	TAB = "    "
	CLI_PROXY = "CliProxy"
	CLI_ITF_STRUCT = "CliItfStruct"
)

func main() {

	pbFilePath := "/Users/zeng/go/src/gsf/src/pbCollection/book.proto"
	pb2ItfFile := getPb2ItfFile(pbFilePath)

	var buffer bytes.Buffer

	//part 0
	impoStr := genImportStr()

	//part 1
	itfStr := genItfStr(pbFilePath)
	itfStructStr := genItfStructStr(pbFilePath)
	itfDoFuncStr := genDoItfFuncStr(pbFilePath)

	buffer.WriteString(itfStr + "\n")
	buffer.WriteString(itfStructStr + "\n")
	buffer.WriteString(itfDoFuncStr + "\n")


	//part 2
	implProxyStructStr := genImplProxyStructStr(pbFilePath)
	implProxyFuncStr := genImplProxyFuncStr(pbFilePath)

	buffer.WriteString(implProxyStructStr + "\n")
	buffer.WriteString(implProxyFuncStr + "\n")


	//part3
	grpcConnStr := genGrpcConnStr()
	buffer.WriteString(grpcConnStr + "\n")

	pb2ItfStr := buffer.String()
	fmt.Println(pb2ItfStr)

	dstFile,err := os.Create("src/pbRpcTool/" + pb2ItfFile)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()

	dstFile.WriteString("package main\n" +
		impoStr + "\n" +
		pb2ItfStr + "\n")

	fmt.Println(impoStr)

}


/*
import (
	"context"
	"google.golang.org/grpc"
	"../pbCollection"
)
 */

func genImportStr() string {
	var buffer bytes.Buffer
	buffer.WriteString("import (\n")
	buffer.WriteString(TAB + "\"context\"" + "\n")
	buffer.WriteString(TAB + "\"google.golang.org/grpc\"" + "\n")
	buffer.WriteString(TAB + "\"../pbCollection\"" + "\n")
	buffer.WriteString(")")

	return buffer.String()
}

func getPb2ItfFile(pbFilePath string) string {
	pbName := pbNameFromPath(pbFilePath)
	fileName := pbName + "Pb2Itf.go"
	return fileName
}


/*

func getGrpcConn() *grpc.ClientConn {
	serviceAddress := "127.0.0.1:50052"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	return conn
}
 */
func genGrpcConnStr() string {
	var buffer bytes.Buffer
	buffer.WriteString("func GetGrpcConn() *grpc.ClientConn {\n")
	buffer.WriteString(TAB + "serviceAddress := \"127.0.0.1:50052\"\n")
	buffer.WriteString(TAB + "conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())\n")
	buffer.WriteString(TAB + "if err != nil {\n")
	buffer.WriteString(TAB + TAB + "panic(\"connect error\")\n")
	buffer.WriteString(TAB + "}\n")
	buffer.WriteString(TAB + "return conn\n")
	buffer.WriteString("}")
	return buffer.String()
}



/*
func (bookCliProxy *BookCliProxy) GetBookInfo (params *book.BookInfoParams) *book.BookInfo {
	bookInfo, _ := bookCliProxy.BookService.GetBookInfo(context.Background(), params)
	return bookInfo
}

func (bookCliProxy *BookCliProxy) GetBookList (params *book.BookListParams) *book.BookList {
	bookList, _ := bookCliProxy.BookService.GetBookList(context.Background(), params)
	return bookList
}
 */
func genImplProxyFuncStr(pbFilePath string) string {
	var buffer bytes.Buffer
	pbName := pbNameFromPath(pbFilePath)
	pbName = Lower(pbName)

	rpcInfoArray := Ioutil(pbFilePath)
	serviceName := getPbServiceName(pbFilePath)

	for _, rpcInfo := range rpcInfoArray {
		buffer.WriteString("func (" + pbName + CLI_PROXY + " *" + Capitalize(pbName) + CLI_PROXY + ") ")

		inParam := rpcInfo.RpcInParam
		outParam := rpcInfo.RpcOutParam
		methodName := rpcInfo.RpcMethodName

		buffer.WriteString(Capitalize(methodName) + " (params *" + Lower(pbName) + "." + inParam + ") *" +
			Lower(pbName) + "." + outParam + "{\n" )
		buffer.WriteString(TAB + Lower(outParam) + ", _ := " + Lower(pbName) + "CliProxy." + serviceName + "." +
			methodName + "(context.Background(), params)" +"\n")
		buffer.WriteString(TAB + "return " + Lower(outParam) + "\n")

		buffer.WriteString("}\n")
	}
	return buffer.String()
}


/*
type BookCliProxy struct {
	BookService book.BookServiceClient
}
 */
func genImplProxyStructStr(pbFilePath string) string {
	var buffer bytes.Buffer
	pbName := pbNameFromPath(pbFilePath)
	buffer.WriteString("type " + Capitalize(pbName) + CLI_PROXY + " " + "struct" + "{\n")
	//获取service name
	serviceName := getPbServiceName(pbFilePath)
	buffer.WriteString(TAB + serviceName + " "+ Lower(pbName) + "." + serviceName + "Client" + "\n" )
	buffer.WriteString("}")
	return buffer.String()
}

/*
func (bookCliItfStruct *BookCliItfStruct) DoGetBookInfo(bookCliItf BookCliItf, bookInfoParams *book.BookInfoParams) *book.BookInfo {
	bookInfo := bookCliItf.GetBookInfo(bookInfoParams)
	return bookInfo
}

func (bookCliItfStruct *BookCliItfStruct) DoGetBookList(bookCliItf BookCliItf, bookListParams *book.BookListParams) *book.BookList {
	bookList := bookCliItf.GetBookList(bookListParams)
	return bookList
}
 */
func genDoItfFuncStr(pbFilePath string) string {
	var buffer bytes.Buffer
	pbName := pbNameFromPath(pbFilePath)
	pbName = Lower(pbName)
	fmt.Println("pbName: ", pbName)

	rpcInfoArray := Ioutil(pbFilePath)

	for _, rpcInfo := range rpcInfoArray {
		buffer.WriteString("func (" + pbName + CLI_ITF_STRUCT + " *" + Capitalize(pbName) + CLI_ITF_STRUCT + ") ")

		inParam := rpcInfo.RpcInParam
		outParam := rpcInfo.RpcOutParam
		methodName := rpcInfo.RpcMethodName
		//fmt.Println(inParam, ",", outParam, ",", methodName)

		buffer.WriteString("" + methodName + "(" + Lower(pbName) + "CliItf " + Capitalize(pbName) + "CliItf, " +
			Lower(inParam) + " *" + Lower(pbName) + "." + Capitalize(inParam) + ") *" + Lower(pbName) + "." +
			Capitalize(outParam) + "{\n" )
		buffer.WriteString(TAB + Lower(outParam) + " := " + Lower(pbName) + "CliItf." + Capitalize(methodName) + "(" +
			Lower(inParam) + ")\n" )
		buffer.WriteString(TAB + "return " + Lower(outParam) + "\n")

		buffer.WriteString("}\n")
	}
	return buffer.String()
}

/*
type BookCliItfStruct struct {

}
 */
func genItfStructStr(pbFilePath string) string {
	pbName := pbNameFromPath(pbFilePath)
	var buffer bytes.Buffer
	//第一行
	buffer.WriteString("type ")
	buffer.WriteString(Capitalize(pbName))
	buffer.WriteString(CLI_ITF_STRUCT + " struct {}")
	return buffer.String()
}

func pbNameFromPath(pbFilePath string) string {
	pbPathSplit := strings.Split(pbFilePath, "/")
	pbName := pbPathSplit[len(pbPathSplit)-1]
	pbName = strings.Split(pbName, ".")[0]
	return pbName
}

/*
type BookCliItf interface {
	GetBookInfo(params *book.BookInfoParams)(*book.BookInfo)
	GetBookList(params *book.BookListParams)(*book.BookList)
}
*/
func genItfStr(pbFilePath string) string {
	pbName := pbNameFromPath(pbFilePath)
	var buffer bytes.Buffer
	//第一行
	buffer.WriteString("type ")
	buffer.WriteString(Capitalize(pbName))
	buffer.WriteString("CliItf interface {\n")
	//rpc行
	rpcInfoArray := Ioutil(pbFilePath)
	for _, rpcInfo := range rpcInfoArray {
		buffer.WriteString(TAB)
		inParam := rpcInfo.RpcInParam
		outParam := rpcInfo.RpcOutParam
		methodName := rpcInfo.RpcMethodName
		buffer.WriteString(methodName)

		buffer.WriteString("(params *" + pbName + "." + inParam + ")")
		buffer.WriteString("(*" + pbName + "." + outParam + ")\n")
	}
	//收尾行
	buffer.WriteString("}")
	return buffer.String()
}

func Ioutil(pbFilePath string) []*RpcInfo {
	if contents, err := ioutil.ReadFile(pbFilePath); err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		//result := string(contents)
		resSplit := strings.Split(result, "\n")
		rpcInfoArray := getRpcInfoArray(resSplit)
		return rpcInfoArray
	}
	return nil
}

func getPbServiceName (pbFilePath string) string {
	if contents, err := ioutil.ReadFile(pbFilePath); err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := strings.Replace(string(contents), "\n", "", 1)
		//result := string(contents)
		resSplit := strings.Split(result, "\n")
		name := getServiceName(resSplit)
		return name
	}
	return ""
}

func getServiceName(resSplit []string) string {
	for _, s := range resSplit {
		if strings.Contains(s, "service") && strings.Contains(s, "{") {
			s = strings.Trim(s, "")
			array := strings.Split(s, " ")
			serviceName := array[1]
			return serviceName
			break
		}
	}
	return ""
}

func getRpcInfoArray(resSplit []string) []*RpcInfo {

	rpcInfoArray := []*RpcInfo{}
	for _, s := range resSplit {
		if strings.Contains(s, "rpc") {
			s = strings.Trim(s, " ")
			array := strings.Split(s, " ")
			/*
				array[1],array[2],array[3],array[4],array[5]
				GetBookInfo (BookInfoParams) returns (BookInfo) {}
			*/
			rpcMethod := array[1]
			inParam := strings.Replace(array[2], "(", "", 1)
			inParam = strings.Replace(inParam, ")", "", 1)
			outParam := strings.Replace(array[4], "(", "", 1)
			outParam = strings.Replace(outParam, ")", "", 1)
			rpcInfo := &RpcInfo{RpcMethodName: rpcMethod, RpcInParam: inParam, RpcOutParam: outParam}
			rpcInfoArray = append(rpcInfoArray, rpcInfo)
		}
	}

	return rpcInfoArray
}

func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				//fmt.Println("Not begins with lowercase letter .")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func Lower(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] < 97 || vv[i] > 122 {
				vv[i] += 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				//fmt.Println("Not begins with uppercase letter .")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

type RpcInfo struct {
	RpcMethodName string
	RpcInParam    string
	RpcOutParam   string
}
