package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/zheng-ji/goSnowFlake"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"strings"
)

func main() {

	//初始化分布式id
	if err := initUUID(); err != nil {
		return
	}

	///fmt.Println("注册")
	//注册
	mongoTest1()

	//fmt.Println("登录")
	//登录
	//mongoTest2()

	//fmt.Println("改密码")
	//改密码
	//mongoTest3()

	//fmt.Println("账号显示")
	//账号显示
	mongoTest4()

	//fmt.Println("修改信息(用户名，头像)")
	//修改信息(用户名，头像)
	//mongoTest5()

	//测试查询不存在
	mongoTest6()

}

//[{295560024959225856 13533401988 123456789111}
// {295721397102055424 13533401988 1234}
// {295721415594741760 13533401988 1234}
// {295721570767212544 13533401988 1234}
// {295721653726351360 13533401988 1234}
// {295721694390128640 13533401988 1234}
// {295721787059081216 13533401988 1234}]

func mongoTest6() {
	account := "13533401988"
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	accountInfoCol := session.DB("test").C("account_basic_info")
	accountInfo := &[]AccountBasicInfo{}

	accountInfoCol.Find(bson.M{"account": account}).All(accountInfo)
	if len(*accountInfo) < 1 {
		fmt.Println("不存在")
	}
	fmt.Println("存在")
	fmt.Println(*accountInfo)

}

func mongoTest5() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	account := "13533401988"
	userName := "newUserName"
	headImg := "newHeadimg"
	//pwd := "123456789111"
	//identifyCode := "123456"

	accountInfoCol := session.DB("test").C("account_basic_info")
	accountInfo := &AccountBasicInfo{}

	usrInfoCol := session.DB("test").C("user_info")
	usrInfo := &UserInfo{}

	accountInfoCol.Find(bson.M{"account": account}).One(accountInfo)

	fmt.Println(*accountInfo)

	uid := accountInfo.User_Id
	log.Println("test5, user_id: ", uid)
	usrInfoCol.Find(bson.M{"user_id": uid}).One(usrInfo)

	usrInfoUpdate := &UserInfo{}
	deepCopy(usrInfoUpdate, usrInfo)
	usrInfoUpdate.Nickname = userName
	usrInfoUpdate.Headimgurl = headImg

	usrByte, err := bson.Marshal(usrInfo)
	if err != nil {
		panic(err)
	}
	mmap := bson.M{}
	err = bson.Unmarshal(usrByte, mmap)
	if err != nil {
		panic(err)
	}

	usrByteUpadate, err := bson.Marshal(usrInfoUpdate)
	if err != nil {
		panic(err)
	}
	mmapUpdate := bson.M{}
	err = bson.Unmarshal(usrByteUpadate, mmapUpdate)
	if err != nil {
		panic(err)
	}

	fmt.Println("mmap:", mmap)
	fmt.Println("mmapUpdate:", mmapUpdate)
	usrInfoCol.Update(mmap, mmapUpdate)
}

func mongoTest4() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	account := "13533401988"
	//identifyCode := "123456"

	accountInfoCol := session.DB("test").C("account_basic_info")
	accountInfo := &AccountBasicInfo{}
	accountInfoCol.Find(bson.M{"account": account}).One(accountInfo)
	accountByte, err := bson.Marshal(accountInfo)
	if err != nil {
		panic(err)
	}
	mmap := bson.M{}
	err = bson.Unmarshal(accountByte, mmap)
	if err != nil {
		panic(err)
	}

	fmt.Println("accountInfo: ", mmap)

}

func mongoTest3() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	account := "13533401988"
	pwd := "123456789111"
	//identifyCode := "123456"

	accountInfoCol := session.DB("test").C("account_basic_info")
	accountInfo := &AccountBasicInfo{}

	accountInfoCol.Find(bson.M{"account": account}).One(accountInfo)

	fmt.Println(accountInfo)

	if strings.EqualFold(accountInfo.Pwd, pwd) {
		fmt.Println("密码一样")
		return
	}

	fmt.Println("密码不一样")

	accountInfoUpdate := &AccountBasicInfo{}
	deepCopy(accountInfoUpdate, accountInfo)

	accountInfoUpdate.Pwd = pwd
	accountByte, err := bson.Marshal(accountInfo)
	if err != nil {
		panic(err)
	}
	mmap := bson.M{}
	err = bson.Unmarshal(accountByte, mmap)
	if err != nil {
		panic(err)
	}

	accountByteUpadate, err := bson.Marshal(accountInfoUpdate)
	if err != nil {
		panic(err)
	}
	mmapUpdate := bson.M{}
	err = bson.Unmarshal(accountByteUpadate, mmapUpdate)
	if err != nil {
		panic(err)
	}

	fmt.Println("mmap:", mmap)
	fmt.Println("mmapUpdate:", mmapUpdate)
	accountInfoCol.Update(mmap, mmapUpdate)

}

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func mongoTest2() {

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	account := "13533401988"
	pwd := "12341111111111"
	//identifyCode := "123456"

	accountInfoCol := session.DB("test").C("account_basic_info")
	accountInfo := &AccountBasicInfo{}
	accountInfoCol.Find(bson.M{"account": account}).One(accountInfo)

	fmt.Println(accountInfo)

	if strings.EqualFold(accountInfo.Pwd, pwd) {
		fmt.Println("密码正确")
	}

}

func initUUID() error {

	iw, err := goSnowFlake.NewIdWorker(1)
	if err != nil {
		log.Fatal(err)
		panic(err)
		return err
	}
	iwCxt = iw
	return nil
}

/*
一个ID 由三部分参与“或”运算组合而成，分别是：

1.毫秒级别的时间戳
2. 机器 workerid
3.以及为了解决冲突的序列号
*/

func mongoTest1() {

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	usrInfoCol := session.DB("test").C("user_info")
	accountInfoCol := session.DB("test").C("account_basic_info")

	//distribute uuid

	id, err := iwCxt.NextId()
	if err != nil {
		panic(err)
	}
	log.Println("uuid:", id)
	uid := strconv.FormatInt(id, 10)

	/*
		1.1手机号+验证码
		1.2输入密码
		1.3输入用户名
		1.4输入公司名称
		1.5注册成功

	*/

	//新用户,可注册
	err = usrInfoCol.Insert(bson.M{"org": "zengxiaosen", "role": "", "user_id": uid, "nickname": "zengxiaosen", "headimgurl": "headingimgurl"})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = accountInfoCol.Insert(bson.M{"user_id": uid, "account": "13533401988", "pwd": "1234"})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Println("注册成功")

	accountInfo := &AccountBasicInfo{}
	accountInfoCol.Find(bson.M{"account": "13533401988"}).One(accountInfo)

	fmt.Println(accountInfo)

}

type UserInfo struct {
	Org        string "bson:`org`"
	Role       string "bson:`role`"
	User_Id    string "bson:`user_id`"
	Nickname   string "bson:`nickname`"
	Headimgurl string "bson:`headimgurl`"
}

type AccountBasicInfo struct {
	User_Id string "bson:`user_id`"
	Account string "bson:`account`"
	Pwd     string "bson:`pwd`"
}

var (
	iwCxt *goSnowFlake.IdWorker
)
