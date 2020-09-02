require "socket"
local JSON = require("JSON")
local UUID = require("uuid")
time = socket.gettime() * 1000
math.randomseed(time)
UUID.randomseed(time)

local function uuid()
    return UUID():gsub('-', '')
end

local charset = { 'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', 'a', 's',
                  'd', 'f', 'g', 'h', 'j', 'k', 'l', 'z', 'x', 'c', 'v', 'b', 'n', 'm', 'Q',
                  'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P', 'A', 'S', 'D', 'F', 'G', 'H',
                  'J', 'K', 'L', 'Z', 'X', 'C', 'V', 'B', 'N', 'M', '1', '2', '3', '4', '5',
                  '6', '7', '8', '9', '0' }

local decset = { '1', '2', '3', '4', '5', '6', '7', '8', '9', '0' }

local function stringRandom(length)
    if length > 0 then
        return stringRandom(length - 1) .. charset[math.random(1, #charset)]
    else
        return ""
    end
end

local function decRandom(length)
    if length > 0 then
        return decRandom(length - 1) .. decset[math.random(1, #decset)]
    else
        return ""
    end
end

local num_users = 962

local function compose()
    local user_index = math.random(1, num_users)
    local username = "username" .. tostring(user_index)
    local user_id = tostring(user_index)

    local text = stringRandom(20)
    local num_user_mentions = math.random(0, 2)
    local num_urls = math.random(0, 2)
    local num_media = math.random(0, 2)
    local media_ids = {}
    local media_types = {}

    for i = 0, num_user_mentions, 1 do
        local user_mention_id
        while (true) do
            user_mention_id = math.random(1, num_users)
            if user_index ~= user_mention_id then
                break
            end
        end
        text = text .. " @username_" .. tostring(user_mention_id)
    end

    for i = 0, num_urls, 1 do
        text = text .. " http://" .. stringRandom(20)
    end

    for i = 0, num_media, 1 do
        local media_id = decRandom(10)
        table.insert(media_ids, media_id)
        table.insert(media_types, "png")
    end

    local path = "https://js1qk9h8q0.execute-api.us-east-1.amazonaws.com/default/txn-dev-Frontend"
    local method = "POST"
    local headers = {}
    local param = {
        InstanceId = uuid(),
        Input = {
            Function = "CreatePost",
            Input = {
                Text = text,
                UserName = username,
                UserId = user_id,
                MediaIds = media_ids,
                MediaTypes = media_types
            }
        }
    }

    local body = JSON:encode(param)
    headers["Content-Type"] = "application/json"

    return wrk.format(method, path, headers, body)
end

local function timeline()
    local user_index = math.random(1, 962)
    user_index = 1
    local user_id = tostring(user_index)

    local path = "https://3mctp9q03d.execute-api.us-east-1.amazonaws.com/default/txn-dev-Frontend"
    local method = "POST"
    local headers = {}
    local param = {
        InstanceId = uuid(),
        Input = {
            Function = "ReadUserTimeline",
            Input = {
                UserId = user_id,
                Start = 0,
                End = 2,
            }
        }
    }

    local body = JSON:encode(param)
    headers["Content-Type"] = "application/json"

    return wrk.format(method, path, headers, body)
end

request = function()
    cur_time = math.floor(socket.gettime())

    return compose()
end