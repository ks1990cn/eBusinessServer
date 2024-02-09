using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;

namespace EbusinessServer.Common
{
    public static class JWTConfiguration
    {
        static async Task CreateJWTConfigurationAsync(WebApplicationBuilder builder)
        {
            string clientId = await ParameterHelper.GetParameterAsync("clientIdUserPool");
            string region = await ParameterHelper.GetParameterAsync("region");
            string userPoolId = await ParameterHelper.GetParameterAsync("userPoolId");
            builder.Services.AddAuthentication(options =>
            {
                options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
                options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
            }).AddJwtBearer(options =>
            {
                options.Authority = $"https://cognito-idp.{region}.amazonaws.com/{userPoolId}";
                options.Audience = clientId;
            });
        }
    }
}
