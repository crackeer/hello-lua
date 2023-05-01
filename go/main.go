package main

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

const luaFile = `function sign(header, params, config, api)
params['ts'] = os.time()
local list = sort(params)
local joinList ={}
for i, v in ipairs(list) do
	joinList[i] = v[1]..'='..v[2]
end

local origin = table.concat(joinList, ',')..config['ak']
params['sign'] = origin
return header, params
end

function sort(data)
local list = {}
local index= 1
for key, val in pairs(data) do
	list[index] = key
	index = index + 1
end
table.sort(list, function(a, b) 
	return a < b
end
)
index = 1
local ret = {}
for i, value in ipairs(list) do
	local tmp = {}
	tmp[1] = value
	tmp[2] = data[value]
	ret[index] = tmp
	index = index + 1
end
return ret
end
`

func main() {
	// 创建一个lua解释器实例
	l := lua.NewState()
	defer l.Close()
	// 需要执行的lua代码
	err := l.DoString(luaFile)
	if err != nil {
		log.Println(err)
	}
	// 执行具体的lua脚本

	input := NewTable(l, map[string]interface{}{
		"simple": "ok",
		"file":   "file",
	})
	header := NewTable(l, map[string]interface{}{})
	config := NewTable(l, map[string]interface{}{
		"ak": "ak",
	})

	err = l.CallByParam(lua.P{
		Fn:      l.GetGlobal("sign"), // 获取info函数引用
		NRet:    2,                   // 指定返回值数量
		Protect: true,                // 如果出现异常，是panic还是返回err
	}, header, input, config) // 传递输入参数n=1
	if err != nil {
		panic(err)
	}
	// 获取返回结果
	ret := l.Get(-1)
	// 从堆栈中删除返回值
	l.Pop(1)
	// 打印返回结果
	fmt.Println(ret)
	if table, ok := ret.(*lua.LTable); ok {
		fmt.Println("Ok")
		table.ForEach(func(l1, l2 lua.LValue) {
			fmt.Println(l1.String(), l2.String())
		})
	}
}

func NewTable(l *lua.LState, data map[string]interface{}) *lua.LTable {
	if data == nil {
		return nil
	}
	input := l.NewTable()
	for key, value := range data {
		if tmp, ok := value.(string); ok {
			input.RawSetString(key, lua.LString(tmp))
			continue
		}
		if tmp, ok := value.(float64); ok {
			input.RawSetString(key, lua.LNumber(tmp))
			continue
		}
		if tmp, ok := value.(map[string]interface{}); ok {
			input.RawSetString(key, NewTable(l, tmp))
			continue
		}
		if tmp, ok := value.(bool); ok {
			input.RawSetString(key, lua.LBool(tmp))
			continue
		}
	}
	return input
}

/*
func CompileLua(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}
*/
