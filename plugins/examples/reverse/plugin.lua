metadata = {
    name = "reverse",
    version = "1.0.0",
    author = "BaconC2",
    description = "Reverses the provided string arguments",
    capabilities = {"command_handler"},
    dependencies = {}
}

function initialize()
    return nil
end

function reverse_string(str)
    local reversed = ""
    for i = #str, 1, -1 do
        reversed = reversed .. str:sub(i, i)
    end
    return reversed
end

function execute(cmd)
    local args = cmd.args or {}
    
    if #args == 0 then
        return {
            status = "failed",
            output = "Error: No arguments provided",
            result_type = "error"
        }
    end
    
    local input = table.concat(args, " ")
    local reversed = reverse_string(input)
    
    return {
        status = "completed",
        output = reversed,
        result_type = "text"
    }
end

function cleanup()
    return nil
end
