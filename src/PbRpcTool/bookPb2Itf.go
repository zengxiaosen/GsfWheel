package main
import (
    "context"
    "google.golang.org/grpc"
    "../pbCollection"
)
type BookCliItf interface {
    GetBookInfo(params *book.BookInfoParams)(*book.BookInfo)
    GetBookList(params *book.BookListParams)(*book.BookList)
}
type BookCliItfStruct struct {}
func (bookCliItfStruct *BookCliItfStruct) DoGetBookInfo(bookCliItf BookCliItf, bookInfoParams *book.BookInfoParams) *book.BookInfo{
    bookInfo := bookCliItf.GetBookInfo(bookInfoParams)
    return bookInfo
}
func (bookCliItfStruct *BookCliItfStruct) DoGetBookList(bookCliItf BookCliItf, bookListParams *book.BookListParams) *book.BookList{
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

