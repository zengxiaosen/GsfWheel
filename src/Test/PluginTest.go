package main

import "fmt"

type pluginfunc interface {
	hello()
	world()
}

//定义一个类，来存放插件
type plugins struct {
	plist map[string]pluginfunc
}

//初始化插件
func (p *plugins)init()  {
	p.plist = make(map[string]pluginfunc)
}

//注册插件
func (p *plugins)register(name string, plugin pluginfunc) {
	p.plist[name] = plugin
}

//plugin1
type plugin1 struct {}
func (p *plugin1) hello (){
	fmt.Println("plugin1 hello")
}
func (p *plugin1) world (){
	fmt.Println("plugin1 world")
}
//plugin2
type plugin2 struct {}
func (p *plugin2) hello() {
	fmt.Println("plugin2 hello")
}
func (p *plugin2) world (){
	fmt.Println("plugin2 world")
}
//plugin3
type plugin3 struct {}
func (p *plugin3) hello() {
	fmt.Println("plugin3 hello")
}
func (p *plugin3) world (){
	fmt.Println("plugin3 world")
}

func main() {
	plugin := new(plugins)
	plugin.init()

	plugin1 := new(plugin1)
	plugin2 := new(plugin2)
	plugin3 := new(plugin3)
	plugin.register("plugin1",plugin1)
	plugin.register("plugin2",plugin2)
	plugin.register("plugin3",plugin3)
	for _,plugin := range(plugin.plist) {
		plugin.hello()
		plugin.world()
	}

}
