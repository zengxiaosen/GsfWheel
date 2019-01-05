package redis

import (
	"errors"
	redigo "github.com/garyburd/redigo/redis"
	"time"
)

var (
	CommandTypeNotValid = errors.New("not valid command type")
	CommandArgsNumNotValid = errors.New("command args num not valid")
	BadPipelineFormat = errors.New("bad pipeline command format")
)

type CommandInfo struct {
	command string
	args []interface{}
}


type Redis struct {
	//pipeline
	commandList []CommandInfo
	command string
	key string
	//ip:port
	addr string
	timeout time.Duration
	network string

}

func New() *Redis {

	r := &Redis{
		timeout:time.Millisecond * 300,
		addr:"localhost:6379",
		network:"tcp",
	}
	return r
}

func GetRedisConn(c *Redis)  (redigo.Conn, interface{}) {
	return redigo.Dial(c.network, c.addr)
}

func (c *Redis) Pipeline(conn redigo.Conn, commandName string, args ...interface{}) (interface{}, error) {
	//使用Add方法添加的命令
	defer c.Clean()

	if len(c.commandList) > 0 {
		cmdCount := len(c.commandList)

		for _, v := range c.commandList {
			err := conn.Send(v.command, v.args...)

			if err != nil {
				return nil, err
			}
		}

		conn.Flush()

		var replys []interface{}

		for i := 0; i< cmdCount; i++ {
			v, err := conn.Receive()

			if err != nil {
				return nil, err
			}

			replys = append(replys, v)
		}

		return replys, nil

	}

	return c.PipelineArgs(conn, commandName, args...)

}

//采用数据包格式:{[]interface{}{"set","zhangsan",1},{"get","zhangsan"}}
func (c *Redis) PipelineArgs(conn redigo.Conn, _ string, args ...interface{}) (interface{}, error) {

	cmdCount := len(args)

	for _, val := range args {

		if commands, ok := val.([]interface{}); ok {

			if len(commands) <= 1 {
				return nil, CommandArgsNumNotValid
			}

			var cmd string

			if cmdStr, ok := commands[0].(string); ok {
				cmd = cmdStr
			} else {
				return nil, CommandTypeNotValid
			}

			leftArgs := commands[1:]

			err := conn.Send(cmd, leftArgs...)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, BadPipelineFormat
		}
	}
	conn.Flush()

	var replys []interface{}

	for i := 0; i < cmdCount; i++ {
		v, err := conn.Receive()

		if err != nil {
			return nil, err
		}

		replys = append(replys, v)
	}

	return replys, nil

}

func (c *Redis) Clean() {
	c.commandList = nil
}

func (c *Redis) Add(commandName string, args ...interface{}) error {
	var command CommandInfo
	command.command = commandName

	for _, v := range args {
		command.args = append(command.args, v)
	}

	c.commandList = append(c.commandList, command)

	return nil
}

