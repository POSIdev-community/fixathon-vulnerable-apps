using System;
using System.Collections.Generic;
using System.IdentityModel.Tokens.Jwt;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Security.Claims;
using System.Text;
using System.Threading.Tasks;
using App.Db;
using App.Models;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.Data.Sqlite;
using Microsoft.IdentityModel.Tokens;
using Newtonsoft.Json;
using NUnit.Framework;

namespace App.Tests
{
    public class Tests
    {
        private const string Cookie = "Cookie";
        private static string RefsDirectory { get; } = Path.Combine("..", "..", "..", "refs");
        
        private HttpClient httpClient;
        private WebApplicationFactory<Startup> application;
        private SqliteConnection connection;
        private string AuthCookie;

        [OneTimeSetUp]
        public void OneTimeSetup()
        {
            connection = new SqliteConnection("Data Source=InMemorySample;Mode=Memory;Cache=Shared");
            connection.Open();
            using var command = connection.CreateCommand();
            command.CommandText = File.ReadAllText(Path.Combine("..", "..", "..", "..", "..", "..", "db_init.sql"));
            command.ExecuteNonQuery();
        }

        [SetUp]
        public void Setup()
        {
            application = new CustomWebApplicationFactory<Startup>();
            httpClient = application.CreateClient(new WebApplicationFactoryClientOptions
            {
                AllowAutoRedirect = false
            });

            AuthCookie = $"jwt={CreateToken(1, "GalacticExplorer")}";
        }

        [Test]
        public async Task AuthController_Login_CorrectUser()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/login");
            httpRequest.Content = new FormUrlEncodedContent(new Dictionary<string, string>()
            {
                ["username"] = "GalacticExplorer",
                ["password"] = "explorer123",
                ["redirect_to"] = "/",
            });

            var response = await httpClient.SendAsync(httpRequest);
            Assert.AreEqual(HttpStatusCode.Redirect, response.StatusCode);
            Assert.IsTrue(response.Headers.Contains("Set-Cookie"));
            var cookie = response.Headers.GetValues("Set-Cookie").First();
            Assert.IsTrue(cookie.StartsWith("jwt="));
            cookie = cookie.Split(";")[0];
            var token = new JwtSecurityTokenHandler().ReadJwtToken(cookie.Substring(4, cookie.Length - 4));
            var userId = token.Claims.FirstOrDefault(c => c.Type == "sub");
            Assert.AreEqual("1", userId?.Value);
        }

        [Test]
        public async Task AuthController_Login_BadNameOrPass()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/login");
            httpRequest.Content = new FormUrlEncodedContent(new Dictionary<string, string>()
            {
                ["username"] = null,
                ["password"] = "explorer123",
                ["redirect_to"] = "/",
            });

            var response = await httpClient.SendAsync(httpRequest);
            Assert.AreEqual(HttpStatusCode.BadRequest, response.StatusCode);

            httpRequest.Content = new FormUrlEncodedContent(new Dictionary<string, string>()
            {
                ["username"] = "GalacticExplorer",
                ["password"] = null,
                ["redirect_to"] = "/",
            });
            Assert.AreEqual(HttpStatusCode.BadRequest, response.StatusCode);
        }

        [Test]
        public async Task AuthController_Login_Unauthorized()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/login");
            httpRequest.Content = new FormUrlEncodedContent(new Dictionary<string, string>()
            {
                ["username"] = "GalacticExplorer",
                ["password"] = "123",
                ["redirect_to"] = "/",
            });

            var response = await httpClient.SendAsync(httpRequest);
            Assert.AreEqual(HttpStatusCode.Unauthorized, response.StatusCode);

            httpRequest.Content = new FormUrlEncodedContent(new Dictionary<string, string>()
            {
                ["username"] = "ExplorerGalactic",
                ["password"] = "explorer123",
                ["redirect_to"] = "/",
            });
            Assert.AreEqual(HttpStatusCode.Unauthorized, response.StatusCode);
        }

        [Test]
        public async Task AuthController_Logout_Redirect()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/logout?redirect=%2F");
            httpRequest.Headers.Add(Cookie, AuthCookie);

            var response = await httpClient.SendAsync(httpRequest);
            Assert.AreEqual(HttpStatusCode.Redirect, response.StatusCode);
            var cookie = response.Headers.GetValues("Set-Cookie").First();
            Assert.IsTrue(cookie.StartsWith("jwt=;"));
        }

        [Test]
        public async Task AuthController_Logout_WithoutRedirect()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/logout");
            httpRequest.Headers.Add(Cookie, AuthCookie);

            var response = await httpClient.SendAsync(httpRequest);
            response.EnsureSuccessStatusCode();
            var cookie = response.Headers.GetValues("Set-Cookie").First();
            Assert.IsTrue(cookie.StartsWith("jwt=;"));
        }

        [Test]
        [Order(1)]
        public async Task ArticlesController_GetAll()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Get, "/api/articles");

            var response = await httpClient.SendAsync(httpRequest);
            response.EnsureSuccessStatusCode();

            var json = await response.Content.ReadAsStringAsync();
            var articles = JsonConvert.DeserializeObject<Article[]>(json);
            Assert.AreEqual(5, articles.Length);
            Assert.AreEqual("Исследование Галактики Андромеды", articles[0].Title);
            Assert.AreEqual("Открытие пришельцев", articles[1].Title);
            Assert.AreEqual("Раскрытие тайн черных дыр", articles[2].Title);
            Assert.AreEqual("Путешествие к краю Вселенной", articles[3].Title);
            Assert.AreEqual("Поиск экзопланет", articles[4].Title);

            Assert.AreEqual(5, articles[4].Id);
            Assert.AreEqual(2, articles[4].UserId);
            Assert.AreEqual("Ученые выявили перспективного кандидата на обитаемую экзопланету, обращающуюся вокруг близкой звезды, разжигая надежды на обнаружение внеземной жизни.", articles[4].Content);
        }

        [Test]
        [Order(2)]
        public async Task ArticlesController_Search()
        {
            Assert.Fail("Not implemented")
            using var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/search");
            var search = JsonConvert.SerializeObject(new SearchViewModel { Search = "Вселенной" });
            httpRequest.Content = new StringContent(search, Encoding.UTF8, "text/json");
            var response = await httpClient.SendAsync(httpRequest);
            response.EnsureSuccessStatusCode();

            var json = await response.Content.ReadAsStringAsync();
            var articles = JsonConvert.DeserializeObject<ArticleViewModel[]>(json);
            Assert.AreEqual(2, articles.Length);
            Assert.AreEqual(3, articles[0].Id);
            Assert.AreEqual("CosmicAdventurer", articles[0].Author);
            Assert.AreEqual("/profiles/3", articles[0].AuthorProfileUrl);
            Assert.AreEqual(4, articles[1].Id);
            Assert.AreEqual("GalacticExplorer", articles[1].Author);
            Assert.AreEqual("/profiles/1", articles[1].AuthorProfileUrl);
        }

        [Test]
        public async Task ArticlesController_SearchNonExistWord()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/search");
            var search = JsonConvert.SerializeObject(new SearchViewModel { Search = "Unknown" });
            httpRequest.Content = new StringContent(search, Encoding.UTF8, "text/json");
            var response = await httpClient.SendAsync(httpRequest);
            response.EnsureSuccessStatusCode();

            var json = await response.Content.ReadAsStringAsync();
            var articles = JsonConvert.DeserializeObject<ArticleViewModel[]>(json);
            Assert.IsEmpty(articles);
        }

        [Test]
        [Order(3)]
        public async Task ArticlesController_Post_Correct()
        {
            using (var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/article_create"))
            {
                httpRequest.Headers.Add(Cookie, AuthCookie);
                httpRequest.Content = new FormUrlEncodedContent(new Dictionary<string, string>()
                {
                    ["title"] = "����� ���������",
                    ["content"] = "�� ����������",
                });

                var response = await httpClient.SendAsync(httpRequest);
                Assert.AreEqual(HttpStatusCode.Redirect, response.StatusCode);
            }

            using (var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/search"))
            {
                var search = JsonConvert.SerializeObject(new SearchViewModel { Search = "���������" });
                httpRequest.Content = new StringContent(search, Encoding.UTF8, "text/json");
                var response = await httpClient.SendAsync(httpRequest);
                response.EnsureSuccessStatusCode();

                var json = await response.Content.ReadAsStringAsync();
                var articles = JsonConvert.DeserializeObject<ArticleViewModel[]>(json);
                Assert.AreEqual(1, articles.Length);
                Assert.AreEqual("����� ���������", articles[0].Title);
                Assert.AreEqual("�� ����������", articles[0].Content);
                Assert.AreEqual(1, articles[0].UserId);
            }
        }

        [Test]
        [Order(3)]
        public async Task ArticlesController_Post_Unauthorized()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Post, "/api/article_create");
            httpRequest.Content = new FormUrlEncodedContent(new Dictionary<string, string>()
            {
                ["title"] = "����� ���������",
                ["content"] = "�� ����������",
            });

            var response = await httpClient.SendAsync(httpRequest);
            Assert.AreEqual(HttpStatusCode.Unauthorized, response.StatusCode);
        }

        [Test]
        [Order(0)]
        public async Task Pages_Article()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Get, "/articles/1");

            var response = await httpClient.SendAsync(httpRequest);
            response.EnsureSuccessStatusCode();
            var reference = File.ReadAllText(Path.Combine(RefsDirectory, "article.html"));
            Assert.AreEqual(reference, await response.Content.ReadAsStringAsync());
        }

        [Test]
        [Order(0)]
        public async Task Pages_Author()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Get, "/profiles/1");

            var response = await httpClient.SendAsync(httpRequest);
            response.EnsureSuccessStatusCode();
            var reference = File.ReadAllText(Path.Combine(RefsDirectory, "author.html"));
            Assert.AreEqual(reference, await response.Content.ReadAsStringAsync());
        }

        [Test]
        [Order(0)]
        public async Task Pages_MyProfile_Authorized()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Get, "/MyProfile");
            httpRequest.Headers.Add(Cookie, AuthCookie);

            var response = await httpClient.SendAsync(httpRequest);
            response.EnsureSuccessStatusCode();
            var reference = File.ReadAllText(Path.Combine(RefsDirectory, "my_profile.html"));
            var responseHtml = await response.Content.ReadAsStringAsync();
            Assert.AreEqual(reference, responseHtml);
        }

        [Test]
        [Order(0)]
        public async Task Pages_MyProfile_Unauthorized()
        {
            using var httpRequest = new HttpRequestMessage(HttpMethod.Get, "/MyProfile");

            var response = await httpClient.SendAsync(httpRequest);
            Assert.AreEqual(HttpStatusCode.Unauthorized, response.StatusCode);
        }

        private string CreateToken(long userId, string username)
        {
            var authClaims = new Claim[]
            {
                new(ClaimTypes.Name, username),
                new(ClaimTypes.NameIdentifier, userId.ToString()),
                new(JwtRegisteredClaimNames.Jti, Guid.NewGuid().ToString()),
                new(JwtRegisteredClaimNames.Sub, userId.ToString()),
            };

            var authSigningKey = new SymmetricSecurityKey(JwtSecret.Default);

            var token = new JwtSecurityToken(
                issuer: "phdays-app",
                expires: DateTime.Now.AddHours(1),
                claims: authClaims,
                signingCredentials: new SigningCredentials(authSigningKey, SecurityAlgorithms.HmacSha256)
            );

            return new JwtSecurityTokenHandler().WriteToken(token);
        }
    }
}