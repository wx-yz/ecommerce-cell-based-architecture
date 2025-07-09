using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.DependencyInjection;
using Grpc.Core;
using Grpc.Net.Client;
using StackExchange.Redis;
using System.Text.Json;
using CartService.Services;

namespace CartService
{
    public class Program
    {
        public static void Main(string[] args)
        {
            CreateHostBuilder(args).Build().Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureWebHostDefaults(webBuilder =>
                {
                    webBuilder.UseStartup<Startup>();
                });
    }

    public class Startup
    {
        public void ConfigureServices(IServiceCollection services)
        {
            // Configure Redis
            var redisConnectionString = Environment.GetEnvironmentVariable("REDIS_ADDR") ?? "localhost:6379";
            services.AddSingleton<IConnectionMultiplexer>(provider =>
            {
                return ConnectionMultiplexer.Connect(redisConnectionString);
            });

            services.AddGrpc();
            services.AddScoped<CartServiceImpl>();
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }

            app.UseRouting();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapGrpcService<CartServiceImpl>();
                endpoints.MapGet("/", async context =>
                {
                    await context.Response.WriteAsync("Cart service is running");
                });
            });
        }
    }
}

namespace CartService.Services
{
    public class CartItem
    {
        public string ProductId { get; set; }
        public int Quantity { get; set; }
    }

    public class Cart
    {
        public string UserId { get; set; }
        public List<CartItem> Items { get; set; } = new List<CartItem>();
    }

    public class CartServiceImpl
    {
        private readonly IConnectionMultiplexer _redis;
        private readonly ILogger<CartServiceImpl> _logger;

        public CartServiceImpl(IConnectionMultiplexer redis, ILogger<CartServiceImpl> logger)
        {
            _redis = redis;
            _logger = logger;
        }

        public async Task<string> AddItemAsync(string userId, CartItem item)
        {
            try
            {
                var db = _redis.GetDatabase();
                var key = $"cart:{userId}";
                
                var cartJson = await db.StringGetAsync(key);
                Cart cart;
                
                if (cartJson.HasValue)
                {
                    cart = JsonSerializer.Deserialize<Cart>(cartJson);
                }
                else
                {
                    cart = new Cart { UserId = userId };
                }

                // Find existing item or add new one
                var existingItem = cart.Items.FirstOrDefault(i => i.ProductId == item.ProductId);
                if (existingItem != null)
                {
                    existingItem.Quantity += item.Quantity;
                }
                else
                {
                    cart.Items.Add(item);
                }

                var updatedCartJson = JsonSerializer.Serialize(cart);
                await db.StringSetAsync(key, updatedCartJson);
                
                _logger.LogInformation($"Added item {item.ProductId} to cart for user {userId}");
                return "OK";
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, $"Error adding item to cart for user {userId}");
                throw;
            }
        }

        public async Task<Cart> GetCartAsync(string userId)
        {
            try
            {
                var db = _redis.GetDatabase();
                var key = $"cart:{userId}";
                
                var cartJson = await db.StringGetAsync(key);
                
                if (cartJson.HasValue)
                {
                    return JsonSerializer.Deserialize<Cart>(cartJson);
                }
                
                return new Cart { UserId = userId };
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, $"Error getting cart for user {userId}");
                throw;
            }
        }

        public async Task<string> EmptyCartAsync(string userId)
        {
            try
            {
                var db = _redis.GetDatabase();
                var key = $"cart:{userId}";
                
                await db.KeyDeleteAsync(key);
                
                _logger.LogInformation($"Emptied cart for user {userId}");
                return "OK";
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, $"Error emptying cart for user {userId}");
                throw;
            }
        }
    }
}