using Microsoft.AspNetCore.Builder;

namespace EbusinessServer.Common
{
    public static class JWTMiddlewareExtensions
    {
        public static IApplicationBuilder UseDownstreamRouteFinderMiddleware(this IApplicationBuilder builder)
        {
            return builder.UseMiddleware<JWTMiddleware>();
        }
    }
}
