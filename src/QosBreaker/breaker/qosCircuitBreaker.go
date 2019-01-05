package qosCircuitBreaker

/*
熔断器的三种状态：
1，关闭状态：服务正常，并维护一个失败率统计，当失败率达到阈值时，转到开启状态
2，开启状态：服务异常，调用fallback函数，一段时间之后，进入半开启状态
3，半开启状态：尝试恢复服务，失败率高于阈值，进入开启状态，低于阈值，进入关闭状态
 */

/*
Breaker is the base of circuit breaker. It maintains failure and success counters as well as the event subscribers
目前先用hyntrix
 */
type Breaker struct {



}


