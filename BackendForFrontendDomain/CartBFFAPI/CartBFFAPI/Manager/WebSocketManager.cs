using Amazon;
using Amazon.SQS;
using Amazon.SQS.Model;
using CartBFFAPI.Helpers;
using System.Collections.Concurrent;
using System.Net.WebSockets;
using System.Text;

namespace CartBFFAPI.Manager
{
    public class WebSocketManager : IWebSocketManager
    {
        private readonly AmazonSQSClient _sqsClient;
        private readonly CancellationToken _cancellationToken;
        private readonly ConcurrentDictionary<WebSocket, bool> _sockets = new ConcurrentDictionary<WebSocket, bool>();

        public WebSocketManager()
        {
            _sqsClient = new AmazonSQSClient(RegionEndpoint.APSouth1);
            _cancellationToken = new CancellationToken();
            StartListening();
        }

        private void StartListening()
        {
            Task.Run(async () =>
            {
                var receiveMessageRequest = new ReceiveMessageRequest
                {
                    QueueUrl = await ParameterHelper.GetParameterAsync("cartqueueurl"),
                    WaitTimeSeconds = 0, // Long-polling
                    MaxNumberOfMessages = 1
                };

                while (!_cancellationToken.IsCancellationRequested)
                {
                    try
                    {
                        var response = await _sqsClient.ReceiveMessageAsync(receiveMessageRequest, _cancellationToken);
                        foreach (var message in response.Messages)
                        {
                            await SendToAllAsync(Encoding.UTF8.GetBytes(message.Body));

                            // Delete the message from the queue
                            await _sqsClient.DeleteMessageAsync(receiveMessageRequest.QueueUrl, message.ReceiptHandle);
                        }
                        Console.WriteLine("Successfully Sent!");
                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine($"Error receiving message from SQS: {ex.Message}");
                    }
                }
            });
        }

        public async Task AddSocket(WebSocket socket)
        {
            _sockets.TryAdd(socket, true);
        }

        public async Task RemoveSocket(WebSocket socket)
        {
            _sockets.TryRemove(socket, out _);
            await socket.CloseAsync(WebSocketCloseStatus.NormalClosure, "Connection closed", _cancellationToken);
        }

        public async Task SendToAllAsync(byte[] buffer)
        {
            foreach (var socket in _sockets.Keys)
            {
                if (socket.State == WebSocketState.Open)
                {
                    await socket.SendAsync(buffer, WebSocketMessageType.Text, true, _cancellationToken);
                }
            }
        }
    }
}
