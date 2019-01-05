package main

import (
	"grpc-go-demo/src/pbCollection"
	"net"
	"context"
	"google.golang.org/grpc"
)

/**
创建BookServer结构 实现 BookServiceServer接口
type BookServiceServer interface {
   GetBookInfo(context.Context, *BookInfoParams) (*BookInfo, error)
   GetBookList(context.Context, *BookListParams) (*BookList, error)
}
*/
type BookServer struct {}
func (s *BookServer) GetBookInfo(ctx context.Context, in *book.BookInfoParams) (*book.BookInfo, error) {
	//请求详情时返回 书籍信息
	b := new(book.BookInfo)
	b.BookId = in.BookId
	b.BookName = "21天精通php"
	return b,nil
}

func (s *BookServer) GetBookList(ctx context.Context, in *book.BookListParams) (*book.BookList, error) {
	//请求列表时返回 书籍列表
	bl := new(book.BookList)
	bl.BookList = append(bl.BookList, &book.BookInfo{BookId:1,BookName:"21天精通php"})
	bl.BookList = append(bl.BookList, &book.BookInfo{BookId:2,BookName:"21天精通java"})
	return bl,nil
}

func main() {
	serviceAddress := ":50052"
	bookServer := new(BookServer)
	//创建tcp监听
	ls, _ := net.Listen("tcp", serviceAddress)
	//创建grpc服务
	gs := grpc.NewServer()
	//注册bookServer
	book.RegisterBookServiceServer(gs, bookServer)
	//启动服务
	gs.Serve(ls)
}

