local vout_history = {}
local vin_history = {}
local time_history = {}
local target_v = 0
local nominal_vin = 0

local window_size = 0.5    
local current_max_t = 0.1  

local min_v_in, max_v_in = 0, 0
local max_v_out = 20

function love.load(arg)
    nominal_vin = tonumber(arg[2]) or 127
    target_v = tonumber(arg[3]) or 0
    
    min_v_in = nominal_vin - 10    
    max_v_in = nominal_vin + 10

    love.window.setMode(1250, 750) 
    love.window.setTitle("SMPS Monitor: Entrada vs Saída")
    io.stdin:setvbuf("no")
end

function love.update(dt)
    local line = io.read("*l")
    if line then
        local vout, vin, t = line:match("([^,]+),([^,]+),([^,]+)")
        if vout and vin and t then
            local t_val = tonumber(t)
            
            table.insert(vout_history, tonumber(vout))
            table.insert(vin_history, tonumber(vin))
            table.insert(time_history, t_val)

            if t_val > current_max_t then
                current_max_t = t_val
            end
            
            if #time_history > 5000 then
                table.remove(vout_history, 1)
                table.remove(vin_history, 1)
                table.remove(time_history, 1)
            end

            if tonumber(vout) > max_v_out * 0.95 then 
                max_v_out = tonumber(vout) * 1.1 
            end
        end
    end
end

function love.draw()
    love.graphics.clear(0.03, 0.03, 0.04)
    local graph_x, graph_w = 80, 850

    local t_end = current_max_t
    local t_start = math.max(0, t_end - window_size)

    draw_graph(graph_x, 50, graph_w, 300, "Entrada (Vin) - Janela Deslizante", min_v_in, max_v_in, vin_history, {0.2, 0.6, 1}, 1, t_start, t_end)

    draw_graph(graph_x, 400, graph_w, 300, "Saída (Vout) - Janela Deslizante", 0, max_v_out, vout_history, {0, 1, 0.5}, 2, t_start, t_end)

    draw_telemetry_panel(970, 50, 250, 650)
end

function draw_graph(x, y, w, h, label, minV, maxV, data, color, stepV, t_start, t_end)
    if not t_start or not t_end then return end

    love.graphics.setColor(0.08, 0.08, 0.1)
    love.graphics.rectangle("fill", x, y, w, h)
    love.graphics.setColor(0.3, 0.3, 0.3)
    love.graphics.rectangle("line", x, y, w, h)

    local rangeV = maxV - minV
    local rangeT = t_end - t_start
    if rangeT <= 0 then rangeT = 0.001 end

    for v = minV, maxV, stepV do
        local ly = (y + h) - ((v - minV) * (h / rangeV))
        if ly >= y and ly <= y + h then
            love.graphics.setColor(1, 1, 1, 0.05)
            love.graphics.line(x, ly, x + w, ly)
            love.graphics.setColor(0.6, 0.6, 0.6)
            love.graphics.print(string.format("%.0fV", v), x - 45, ly - 7)
        end
    end

    local divs = 5
    for i = 0, divs do
        local lx = x + (i * (w / divs))
        local t_tick = t_start + (i * (rangeT / divs))
        love.graphics.setColor(1, 1, 1, 0.05)
        love.graphics.line(lx, y, lx, y + h)
        love.graphics.setColor(0.5, 0.5, 0.5)
        love.graphics.print(string.format("%.2fs", t_tick), lx - 15, y + h + 8)
    end

    love.graphics.setScissor(x, y, w, h)
    if #data > 1 then
        love.graphics.setColor(color)
        love.graphics.setLineWidth(2)
        local points = {}
        for i = 1, #data do
            local tv = time_history[i]
            if tv >= t_start - 0.05 and tv <= t_end + 0.05 then
                local px = x + ((tv - t_start) * (w / rangeT))
                local py = (y + h) - ((data[i] - minV) * (h / rangeV))
                table.insert(points, px)
                table.insert(points, py)
            end
        end
        if #points >= 4 then love.graphics.line(points) end
    end
    love.graphics.setScissor()
    
    love.graphics.setColor(1, 1, 1)
    love.graphics.print(label, x, y - 25)
end

function draw_telemetry_panel(x, y, w, h)
    love.graphics.setColor(0.1, 0.1, 0.15)
    love.graphics.rectangle("fill", x, y, w, h, 10)
    love.graphics.setColor(0.2, 0.6, 1, 0.5)
    love.graphics.rectangle("line", x, y, w, h, 10)

    local last_vout = vout_history[#vout_history] or 0
    local last_vin = vin_history[#vin_history] or 0
    local last_t = time_history[#time_history] or 0

    love.graphics.setColor(1, 1, 1)
    love.graphics.print("TELEMETRIA SMPS", x + 20, y + 20)
    love.graphics.setColor(0.7, 0.7, 0.7)
    love.graphics.print(string.format("Tempo: %.4f s", last_t), x + 20, y + 60)
    love.graphics.setColor(0.2, 0.6, 1)
    love.graphics.print(string.format("Vin Real: %.2f V", last_vin), x + 20, y + 100)
    love.graphics.setColor(0, 1, 0.5)
    love.graphics.print(string.format("Vout: %.4f V", last_vout), x + 20, y + 140)
    love.graphics.setColor(1, 0.3, 0.3)
    love.graphics.print(string.format("Erro: %.4f V", target_v - last_vout), x + 20, y + 180)
end