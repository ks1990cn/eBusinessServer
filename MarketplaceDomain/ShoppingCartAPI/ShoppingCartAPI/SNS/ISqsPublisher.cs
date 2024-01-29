
namespace ShoppingCartAPI.SNS
{
    public interface ISqsPublisher
    {
        Task PublishMessageAsync(string message);
    }
}