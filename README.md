# gsf: 一些golang 服务化工具 
## 针对的问题：业务和公共代码耦合度高，影响工程质量
## 解决的方式：基于golang的服务化工具
## 预期效果：达到的预期是中台逻辑聚拢，提升开发迭代速率，提升业务服务质量。
## 最终目的：服务治理


# 1,关于Grpc
http://www.grpc.io/docs/quickstart/go.html
protoc -I ./ helloworld.proto --go_out=plugins=grpc:.

# 2,切面工具 (做了非反射和反射两个方式)
https://github.com/xiaosenzeng/gsf/tree/master/src/Aop/AopUtil.go
https://github.com/xiaosenzeng/gsf/tree/master/src/Test/AopTest.go

# 3,运行时异常及自定义异常捕获处理工具
https://github.com/xiaosenzeng/gsf/tree/master/src/TryCatch/TryCatchFinally.go

# 4,反射的套路用法
https://github.com/xiaosenzeng/gsf/tree/master/src/Reflect/ReflectTool.go

# 5,代理模式接口化Grpc，优化rpc技术，开发者只需要根据pb定义对应的接口，然后把实现类和参数传入接口方法
https://github.com/xiaosenzeng/gsf/tree/master/src/NetworkIO/client/bookClient.go

# 6,基于5，研发了pb生成接口的工具，pb和golang的接口完全一一对应。脚本一键生成接口和实现类，用户只需要传入参数，无需写业务无关的代码
https://github.com/xiaosenzeng/gsf/tree/master/src/PbRpcTool/Pb2ItfImplTool.go  
pb规范:  
1）proto的文件名和service名对应，比如book.proto <--> BookService  
2）待补充  
生成接口文件样例:      
生成的接口文件如https://github.com/xiaosenzeng/gsf/tree/master/src/PbRpcTool/bookPb2Itf.go所示  
使用方式:  
直接将生成的接口文件内容拷贝到grpc的业务client代码端即可  
再业务client端利用接口结构体调用接口即可  
```
bookCliItfStruct := &BookCliItfStruct{}  
bookCliProxy := BookCliProxy{BookService:bookClient}  

第一条rpc 用户只需要写两行代码  
bookInfoReqParams := book.BookInfoParams{BookId:1}
BookInfo := bookCliItfStruct.GetBookInfo(&bookCliProxy, &bookInfoReqParams) 

fmt.Println("获取书籍详情")
fmt.Println("bookId: 1", " => ", "bookName:", BookInfo.BookName) 

第二条rpc 用户只需要写两行代码   
bookListReqParams := book.BookListParams{Page:1,Limit:10}
BookList := bookCliItfStruct.GetBookList(&bookCliProxy, &bookListReqParams)

fmt.Println("获取书籍列表")
for _, b := range BookList.BookList {
    fmt.Println("bookId: ", b.BookId, " => ", "bookName: ", b.BookName)
}
```
总结：rpc接口化，开发者可以更关注于本地业务，像调用本地方法一样调用本地服务   
# 7，基于grpc的网关层协议转换，http qapp rpc 三层互通
http->网关反响代理->grpc stub的数据流：    
http层：https://github.com/xiaosenzeng/gsf/tree/master/src/GateWay/gateWayProtocolTransLayer/helloGateWayHttpLayer.go  
网关server层：https://github.com/xiaosenzeng/gsf/tree/master/GateWay/src/gateWayServer/helloGateWayServer.go  
qapp层：待补充  
# 8，JWT安全验证
https://github.com/xiaosenzeng/gsf/tree/master/src/Jwt  
## JWT的流程(基于gin的JWT p0测试)：  
### 8.1，获取jwt的token：http://localhost:8088/auth?username=z&password=z  
### 8.2，携带token进行jwt验证：http://localhost:8088/hello?name=z&token=XXX  
XXX为8.1 生成的jwt token  
## gin的源码解构
![Image text](https://img-blog.csdnimg.cn/20181206135844644.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3plbmd4aWFvc2Vu,size_16,color_FFFFFF,t_70)
# 9，基于grpc的限流服务  
## 9.1目前比较经典的限流算法 https://github.com/xiaosenzeng/gsf/tree/master/src/Test/rateLimitAlgorithms.go 
计数器限流算法  
漏桶限流算法  
计数器限流算法  
## 9.2 grpc server端限流
1,client端请求到达grpc svr端 https://github.com/xiaosenzeng/gsf/tree/master/src/QosRateLimit/rateLimitServer/helloServer.go  
2,走拦截器
```
s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			qosRateLimit.LimitRateServerInterceptor(),
		)
	  )
    )
go-grpc-middleware:这个项目对grpc的interceptor进行了封装，支持多个拦截器的链式组装。   
```
3,拦截器的限流服务 https://github.com/xiaosenzeng/gsf/tree/master/src/QosRateLimit/rateLimit/interceptor.go  
4,分发令牌的核心算法服务 https://github.com/xiaosenzeng/gsf/tree/master/src/QosRateLimit/rateLimit/qosRateLimit.go  
# 10，基于grpc的服务熔断技术  
## 10.1 目前业界的主流的熔断方法：github.com/afex/hystrix-go/hystrix
```
//熔断器
hystrix.ConfigureCommand(
    //熔断器名字，可以用服务名称命名，一个名字对应一个熔断器，对应一份熔断策略
    "qosCircuitBreakerService",
    hystrix.CommandConfig{
        //超时时间 100ms
        Timeout: 100,
        //最大并发数，超过并发返回错误
        MaxConcurrentRequests: 2,
        //请求数量的阀值，用这些数量的请求来计算阈值
        RequestVolumeThreshold: 4,
        //错误率阈值，达到阈值，启动熔断器，25%
        ErrorPercentThreshold: 25,
        //熔断尝试恢复时间，1000ms
        SleepWindow: 1000,
    },
)
```
## 10.2 hystrix的两种熔断模式：阻塞和非阻塞，一般于多个资源需要并发访问的场景下使用非阻塞模式
## 10.3 服务熔断的p0测试，在rpc的client端切如拦截器
https://github.com/xiaosenzeng/gsf/tree/master/src/QosBreaker/breakerClient/helloClient.go  
https://github.com/xiaosenzeng/gsf/tree/master/src/QosBreaker/breakerServer/helloServer.go  
代码中test breaker模拟来1w并发请求  
# 11，基于grpc的服务注册发现与负载均衡技术
![Image text](https://segmentfault.com/img/bVKyon?w=554&h=226)
## 11.1 loadbalance中间件
https://github.com/xiaosenzeng/gsf/tree/master/src/QosLoadBalance/lbMdw   
## 11.2 loadbalance client server两端p0测试
https://github.com/xiaosenzeng/gsf/tree/master/src/QosLoadBalance/lbClient/helloClient.go  
https://github.com/xiaosenzeng/gsf/tree/master/src/QosLoadBalance/lbServer/helloServer.go     
# 12，基于grpc的分布式服务链路追踪
## 12.1 目前业界主流的链路追踪组件有google的Dapper，Twitter的zipkin和阿里的Eagleeye（鹰眼），京东商城的Hydra，eBay的CAL
![Image text](https://images0.cnblogs.com/blog/7438/201412/121609597285837.gif)
![Image text](https://images0.cnblogs.com/blog/7438/201412/121611003065809.png)
## 12.2 本工程基于grpc集成zipkin做链路追踪
https://github.com/xiaosenzeng/gsf/tree/master/src/QosTracer/tracerClient/helloClient.go  
https://github.com/xiaosenzeng/gsf/tree/master/src/QosTracer/tracerServer/helloServer.go  
# 13 Redis工具
## 13.1 Redis分布式锁工具
https://github.com/xiaosenzeng/gsf/tree/master/src/Storage/redis/DistributeLockutil.go  
## 13.2 Redis Pipeline工具
https://github.com/xiaosenzeng/gsf/tree/master/src/Storage/redis/PipelineUtil.go  
# 14 SQL ORM工具
https://github.com/xiaosenzeng/gsf/tree/master/src/Storage/mysql/xorm.go  
# 15 构建分布式定时器，cronJob,delayJob，最终构建统一调度系统
## 场景：跟grpc结合，grpc server任务的定时调度
## 分布式调度针对的问题
### 一，传统crontab的缺陷
1）配置任务时，需要ssh登陆脚本服务器进行操作  
2）服务器宕机，任务停止，单点问题  
3）排查问题低效，无法方便查看任务状态和错误输出  
### 二，开源项目的选型问题
1）版本迭代慢，相关企业维护精力有限
2）文档少，社区不太友好
## 国内外现状
1）当当网的Elastic Job  
2）点评网的XX－JOB
## 根据分布式调度针对的问题，提出分布式调度系统的解决方案 
1）可视化任务管理  
2）分布式架构，解决单点问题  
3）log采集跟踪  
## 根据分布式调度rpc过程中网络不稳定性导致数据不一致的问题，提出一种有效的分布式调度策略 
1）利用etcd同步全量任务列表到所有worker节点  
2）每个worker独立调度全量任务，无需与master产生直接rpc  
3）每个worker利用分布式锁抢占，解决并发调度相同任务的问题  
## 代码:
https://github.com/zengxiaosen/GsfWheel/tree/master/src/DistributeCrotab
## 界面
![Image text](https://img-blog.csdnimg.cn/20190105165235570.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3plbmd4aWFvc2Vu,size_16,color_FFFFFF,t_70)
## worker服务注册(使用租约是为了worker宕机让他自动下限，也就是利用租约到期来避免worker宕机还在列表中，过期后key就会被删除，从而下线)
1) 启动后获取本机网卡ip作为节点唯一标识  
2) 启动服务注册协程，首先创建lease并且自动续租
3) 带着lease注册到/cron/workers/{IP}下，供服务发现  
#


  


  


 







