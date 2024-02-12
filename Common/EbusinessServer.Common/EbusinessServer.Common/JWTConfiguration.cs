using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.IdentityModel.Tokens;

namespace EbusinessServer.Common
{
    public static class JWTConfiguration
    {
        public static async Task CreateJWTConfigurationAsync(WebApplicationBuilder builder)
        {
            string clientId = await ParameterHelper.GetParameterAsync("clientIdUserPool");
            string region = await ParameterHelper.GetParameterAsync("region");
            string userPoolId = await ParameterHelper.GetParameterAsync("userPoolId");
            builder.Services.AddTransient<JWTMiddleware>();
            builder.Services.AddAuthentication(options =>
            {
                options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
                options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
            }).AddJwtBearer(JwtBearerDefaults.AuthenticationScheme,options =>
            {
                options.Authority = $"https://cognito-idp.{region}.amazonaws.com/{userPoolId}";
                options.Audience = clientId;
                options.TokenValidationParameters = new TokenValidationParameters
                {
                    ValidateIssuer = true,
                    ValidateAudience = true,
                    ValidateLifetime = true,
                    ValidateIssuerSigningKey = true,
                    ValidIssuer = $"https://cognito-idp.{region}.amazonaws.com/{userPoolId}",
                    ValidAudience = clientId
                };
            });
        }
    }
}
