using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;

namespace EbusinessServer.Common
{
    public static class JWTConfiguration
    {
        static void CreateJWTConfiguration(WebApplicationBuilder builder)
        {
            builder.Services.AddAuthentication(options =>
            {
                options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
                options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
            }).AddJwtBearer(options =>
            {
                options.Authority = "https://cognito-idp.{region}.amazonaws.com/{userPoolId}";
                options.Audience = "{clientId}";
            });
        }
    }
}
