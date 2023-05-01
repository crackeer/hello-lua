function max(num1, num2)

    if (num1 > num2) then
       result = num1;
    else
       result = num2;
    end
 
    return result;
 end
 print("两值比较最大值为 ",max(10,4))
 print("两值比较最大值为 ",max(5,6))

function sign(header, params, config, api)
    params['ts'] = os.time()
    local list = sort(params)
    local joinList ={}
    for i, v in ipairs(list) do
        joinList[i] = v[1]..'='..v[2]
    end

    local origin = table.concat(joinList, ',')..config['ak']
    local md5 = require("./md5")
    params['sign'] = md5.sumhexa(origin)
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

local input = {}
input["ccc"] = 11
input["uuu"] = 9
input["bbb"] = 10

local config = {}
config['ak'] = 'shs'

local header, params = sign({}, input, config, {})
for key, value in pairs(params) do
    print(key, value)
end
print("参数类型 ",header, params)
