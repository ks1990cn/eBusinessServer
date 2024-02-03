using System.Net.WebSockets;

namespace CartBFFAPI.Middleware
{
    public class WebSocketMiddleware
    {
        private readonly RequestDelegate _next;
        private readonly Manager.IWebSocketManager _webSocketManager;

        public WebSocketMiddleware(RequestDelegate next, Manager.IWebSocketManager webSocketManager)
        {
            _next = next;
            _webSocketManager = webSocketManager;
        }

        public async Task Invoke(HttpContext context)
        {
            if (!context.WebSockets.IsWebSocketRequest)
            {
                await _next(context);
                return;
            }

            var cancellationToken = context.RequestAborted;
            var socket = await context.WebSockets.AcceptWebSocketAsync();

            await _webSocketManager.AddSocket(socket);

            var buffer = new byte[1024 * 4];
            while (socket.State == WebSocketState.Open)
            {
                var result = await socket.ReceiveAsync(new ArraySegment<byte>(buffer), cancellationToken);
                if (result.MessageType == WebSocketMessageType.Close)
                {
                    await _webSocketManager.RemoveSocket(socket);
                }
            }
        }
    }
}
