using Amazon;
using Amazon.SimpleSystemsManagement;
using Amazon.SimpleSystemsManagement.Model;

namespace EbusinessServer.Common
{
    public static class ParameterHelper
    {
        private readonly static AmazonSimpleSystemsManagementClient _ssmClient;

        static ParameterHelper()
        {
            _ssmClient = new AmazonSimpleSystemsManagementClient(RegionEndpoint.APSouth1);
        }
        public static async Task<string> GetParameterAsync(string parameterName)
        {
            var request = new GetParameterRequest
            {
                Name = parameterName,
                WithDecryption = true // Set to true if your parameter value is encrypted
            };

            var response = await _ssmClient.GetParameterAsync(request);
            return response.Parameter.Value;
        }
    }
}
