require "socket"
local JSON = require("JSON")
local UUID = require("uuid")
time = socket.gettime() * 1000
UUID.randomseed(time)
math.randomseed(socket.gettime() * 1000)
math.random();
math.random();
math.random()

local function uuid()
    return UUID():gsub('-', '')
end

local gatewayPath = os.getenv("ENDPOINT")

local function get_user()
    local id = math.random(0, 500)
    local user_name = "Cornell_" .. tostring(id)
    local pass_word = ""
    for i = 0, 9, 1 do
        pass_word = pass_word .. tostring(id)
    end
    return user_name, pass_word
end

local function search_hotel()
    local in_date = math.random(9, 23)
    local out_date = math.random(in_date + 1, 24)

    local in_date_str = tostring(in_date)
    if in_date <= 9 then
        in_date_str = "2015-04-0" .. in_date_str
    else
        in_date_str = "2015-04-" .. in_date_str
    end

    local out_date_str = tostring(out_date)
    if out_date <= 9 then
        out_date_str = "2015-04-0" .. out_date_str
    else
        out_date_str = "2015-04-" .. out_date_str
    end

    local lat = 38.0235 + (math.random(0, 481) - 240.5) / 1000.0
    local lon = -122.095 + (math.random(0, 325) - 157.0) / 1000.0

    local method = "GET"

    local path = gatewayPath
    local param = {
        InstanceId = uuid(),
        CallerName = "",
        Async = true,
        Input = {
            Function = "search",
            Input = {
                Lat = lat,
                Lon = lon,
                InDate = in_date_str,
                OutDate = out_date_str
            },
        }
    }
    local body = JSON:encode(param)

    local headers = {}
    headers["Content-Type"] = "application/json"
    return wrk.format(method, path, headers, body)
end

local function recommend()
    local coin = math.random()
    local req_param = ""
    if coin < 0.33 then
        req_param = "dis"
    elseif coin < 0.66 then
        req_param = "rate"
    else
        req_param = "price"
    end

    local lat = 38.0235 + (math.random(0, 481) - 240.5) / 1000.0
    local lon = -122.095 + (math.random(0, 325) - 157.0) / 1000.0

    local method = "GET"
    local path = gatewayPath

    local param = {
        InstanceId = uuid(),
        CallerName = "",
        Async = true,
        Input = {
            Function = "recommend",
            Input = {
                Require = req_param,
                Lat = lat,
                Lon = lon,
            }
        }
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"
    return wrk.format(method, path, headers, body)
end

local function reserve()
    local in_date = math.random(9, 23)
    local out_date = in_date + math.random(1, 5)

    local in_date_str = tostring(in_date)
    if in_date <= 9 then
        in_date_str = "2015-04-0" .. in_date_str
    else
        in_date_str = "2015-04-" .. in_date_str
    end

    local out_date_str = tostring(out_date)
    if out_date <= 9 then
        out_date_str = "2015-04-0" .. out_date_str
    else
        out_date_str = "2015-04-" .. out_date_str
    end

    local hotel_id = tostring(math.random(1, 80))
    local user_id, password = get_user()
    local cust_name = user_id

    local num_room = "1"

    local method = "POST"
    local path = "http://localhost:5000/reservation?inDate=" .. in_date_str ..
            "&outDate=" .. out_date_str .. "&lat=" .. tostring(lat) .. "&lon=" .. tostring(lon) ..
            "&hotelId=" .. hotel_id .. "&customerName=" .. cust_name .. "&username=" .. user_id ..
            "&password=" .. password .. "&number=" .. num_room
    local headers = {}
    return wrk.format(method, path, headers, nil)
end

local function reserve_all()
    local method = "GET"
    local path = gatewayPath

    local param = {
        InstanceId = uuid(),
        CallerName = "",
        Async = true,
        Input = {
            Function = "reserve",
            Input = {
                userId = "user1",
                hotelId = tostring(math.random(0, 99)),
                flightId = tostring(math.random(0, 99)),
            }
        }
    }
    local body = JSON:encode(param)

    local headers = {}
    headers["Content-Type"] = "application/json"
    return wrk.format(method, path, headers, body)
end

local function user_login()
    local user_name, password = get_user()
    local method = "GET"
    local path = gatewayPath
    local param = {
        InstanceId = uuid(),
        CallerName = "",
        Async = true,
        Input = {
            Function = "user",
            Input = {
                Username = user_name,
                Password = password
            }
        }
    }
    local body = JSON:encode(param)
    local headers = {}
    headers["Content-Type"] = "application/json"

    return wrk.format(method, path, headers, body)
end

request = function()
    cur_time = math.floor(socket.gettime())
    local search_ratio = 0.6
    local recommend_ratio = 0.39
    local user_ratio = 0.005
    local reserve_ratio = 0.005

    --return search_hotel()
    --return recommend()
    --return user_login()
    --return reserve_all()
    local coin = math.random()
    if coin < search_ratio then
        return search_hotel()
    elseif coin < search_ratio + recommend_ratio then
        return recommend()
    elseif coin < search_ratio + recommend_ratio + user_ratio then
        return user_login()
    else
        return reserve_all()
    end
end
