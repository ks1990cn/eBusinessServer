using Amazon;
using Amazon.SQS;
using Amazon.SQS.Model;
using System.Text.Json;
using EbusinessServer.Common;
namespace ShoppingCartAPI.SNS
{
    public class SqsPublisher : ISqsPublisher
    {
        private readonly AmazonSQSClient _sqsClient;
        public SqsPublisher()
        {
            _sqsClient = new AmazonSQSClient(RegionEndpoint.APSouth1);
        }

        public async Task PublishMessageAsync(string message)
        {
            var request = new SendMessageRequest()
            {
                QueueUrl = await ParameterHelper.GetParameterAsync("cartqueueurl"),
                MessageBody = JsonSerializer.Serialize(message)
            };
            try
            {
                var response = await _sqsClient.SendMessageAsync(request);
                Console.WriteLine("Message published. MessageId: " + response.MessageId);
            }
            catch (Exception ex)
            {
                Console.WriteLine("Error publishing message: " + ex.Message);
            }
        }

    }
}
