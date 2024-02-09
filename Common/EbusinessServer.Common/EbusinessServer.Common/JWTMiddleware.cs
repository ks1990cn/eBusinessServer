using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Http;

namespace EbusinessServer.Common
{
    internal class JWTMiddleware : IMiddleware
    {
        public async Task InvokeAsync(HttpContext context, RequestDelegate next)
        {
            if (context.Request.Headers.TryGetValue("Authorization", out var token))
            {
                var result = await context.AuthenticateAsync(JwtBearerDefaults.AuthenticationScheme);
                if (result.Succeeded)
                {
                    await next(context);
                    return;
                }
            }

            context.Response.StatusCode = StatusCodes.Status401Unauthorized;
        }
    }
}
