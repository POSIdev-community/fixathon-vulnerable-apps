using System.IO;
using App.Db;
using App.Tests;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.AspNetCore.TestHost;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;

public class CustomWebApplicationFactory<TStartup> : WebApplicationFactory<TStartup> where TStartup : class
{

    public IConfiguration Configuration { get; private set; }

    protected override void ConfigureWebHost(IWebHostBuilder builder)
    {
        builder.ConfigureAppConfiguration(config =>
        {
            Configuration = new ConfigurationBuilder()
                .AddJsonFile("appsettings.Testing.json")
                .Build();

            config.AddConfiguration(Configuration);
        });
        /*
        builder.ConfigureServices(services =>
        {
            var sc = services.AddScoped<AppDbContext, TestAppDbContext>();
            var sp = sc.BuildServiceProvider();
            var context = sp.GetService<AppDbContext>();
            context.Database.EnsureDeleted();
            context.Database.EnsureCreated();
            context.Users.Add(new User
            {
                Id = 1,
                Name = "myuser",
                Password = "mypassword"
            });
            context.SaveChanges();
        });*/
    }

    protected override IHostBuilder CreateHostBuilder()
    {
        var builder = Host.CreateDefaultBuilder()
            .ConfigureWebHostDefaults(x =>
            {
                x
                    .UseStartup<App.Startup>()
                    .UseTestServer()
                    ;
            });
        return builder;
    }
}