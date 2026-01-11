metadata = {
    name = "echo",
    version = "1.0.0",
    author = "BaconC2",
    description = "Echoes back the provided arguments",
    capabilities = {"command_handler"},
    dependencies = {}
}

function initialize()
    return nil
end

function execute(cmd)
    local args = cmd.args or {}
    
    if #args == 0 then
        return {
            status = "completed",
            output = "Echo: (empty)",
            result_type = "text"
        }
    end
    
    local output = "Echo: " .. table.concat(args, " ")
    
    return {
        status = "completed",
        output = output,
        result_type = "text"
    }
end

function cleanup()
    return nil
end
