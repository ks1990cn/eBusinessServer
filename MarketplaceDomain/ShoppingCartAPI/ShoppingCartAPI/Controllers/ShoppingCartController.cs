using Amazon.SimpleNotificationService.Util;
using Microsoft.AspNetCore.Mvc;
using ShoppingCartAPI.SNS;

// For more information on enabling Web API for empty projects, visit https://go.microsoft.com/fwlink/?LinkID=397860

namespace ShoppingCartAPI.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class ShoppingCartController : ControllerBase
    {
        private readonly ISqsPublisher _snsPublisher;
        public static List<int> ProductIds = new List<int>();
        public ShoppingCartController(ISqsPublisher snsPublisher)
        {
           _snsPublisher = snsPublisher;
        }
        // GET: api/<ShoppingCartController>
        [HttpGet]
        public IEnumerable<int> Get()
        {
            return ProductIds;
        }

        // GET api/<ShoppingCartController>/5
        [HttpGet("{id}")]
        public string Get(int id)
        {
            return "value";
        }

        // POST api/<ShoppingCartController>
        [HttpPost]
        public async Task Post([FromQuery] int productId)
        {
            ProductIds.Add(productId);
            await _snsPublisher.PublishMessageAsync(productId.ToString());
        }

        // PUT api/<ShoppingCartController>/5
        [HttpPut("{id}")]
        public void Put(int id, [FromBody] string value)
        {
        }

        // DELETE api/<ShoppingCartController>/5
        [HttpDelete("{id}")]
        public void Delete(int id)
        {
        }
    }
}
