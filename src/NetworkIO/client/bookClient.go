package main


import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gsf/src/pbCollection"
)
func main()  {
	conn := GetGrpcConn()
	defer conn.Close()
	bookClient := book.NewBookServiceClient(conn)


	/*
	面向接口编程，优化了rpc
	只需定义bookCliItfStruct接口结构体
	直接调用接口结构体对应的DoGetBookInfo，DoGetBookList接口方法

	开发者需要做的是传入一个实现了接口的结构体即可
	 */
	bookCliItfStruct := &BookCliItfStruct{}
	bookCliProxy := BookCliProxy{BookService:bookClient}

	//
	bookInfoReqParams := book.BookInfoParams{BookId:1}
	BookInfo := bookCliItfStruct.GetBookInfo(&bookCliProxy, &bookInfoReqParams)

	fmt.Println("获取书籍详情")
	fmt.Println("bookId: 1", " => ", "bookName:", BookInfo.BookName)

	//
	bookListReqParams := book.BookListParams{Page:1,Limit:10}
	BookList := bookCliItfStruct.GetBookList(&bookCliProxy, &bookListReqParams)

	fmt.Println("获取书籍列表")
	for _, b := range BookList.BookList {
		fmt.Println("bookId: ", b.BookId, " => ", "bookName: ", b.BookName)
	}

}



type BookCliItf interface {
	GetBookInfo(params *book.BookInfoParams)(*book.BookInfo)
	GetBookList(params *book.BookListParams)(*book.BookList)
}

type BookCliItfStruct struct {}
func (bookCliItfStruct *BookCliItfStruct) GetBookInfo(bookCliItf BookCliItf, bookInfoParams *book.BookInfoParams) *book.BookInfo{
	bookInfo := bookCliItf.GetBookInfo(bookInfoParams)
	return bookInfo
}
func (bookCliItfStruct *BookCliItfStruct) GetBookList(bookCliItf BookCliItf, bookListParams *book.BookListParams) *book.BookList{
	bookList := bookCliItf.GetBookList(bookListParams)
	return bookList
}

type BookCliProxy struct{
	BookService book.BookServiceClient
}
func (bookCliProxy *BookCliProxy) GetBookInfo (params *book.BookInfoParams) *book.BookInfo{
	bookInfo, _ := bookCliProxy.BookService.GetBookInfo(context.Background(), params)
	return bookInfo
}
func (bookCliProxy *BookCliProxy) GetBookList (params *book.BookListParams) *book.BookList{
	bookList, _ := bookCliProxy.BookService.GetBookList(context.Background(), params)
	return bookList
}

func GetGrpcConn() *grpc.ClientConn {
	serviceAddress := "127.0.0.1:50052"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	return conn
}

