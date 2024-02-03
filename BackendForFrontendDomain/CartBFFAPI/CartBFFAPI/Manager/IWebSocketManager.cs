using System.Net.WebSockets;

namespace CartBFFAPI.Manager
{
    public interface IWebSocketManager
    {
        Task AddSocket(WebSocket socket);
        Task RemoveSocket(WebSocket socket);
        Task SendToAllAsync(byte[] buffer);
    }
}