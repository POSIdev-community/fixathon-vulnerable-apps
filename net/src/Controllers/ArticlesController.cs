using System.Collections.Generic;
using System.Linq;
using App.Db;
using App.Models;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;

namespace App.Controllers
{
    [ApiController]
    [Route("api/articles")]
    public class ArticlesController : ControllerBase
    {
        private readonly AppDbContext dbContext;
        private readonly ILogger<ArticlesController> logger;

        private static readonly string searchSql = @"
select a.articleId, a.title, a.content, a.userId, u.username from articles a
join users u on a.userId = u.userId
where a.title like '%{0}%' or a.content like '%{0}%'";

        public ArticlesController(AppDbContext dbContext, ILogger<ArticlesController> logger)
        {
            this.dbContext = dbContext;
            this.logger = logger;
        }

        [HttpGet]
        public IActionResult GetAll()
        {
            var articles = dbContext.Articles
                .Join(dbContext.Users,
                    a => a.UserId,
                    u => u.Id,
                    (a, u) => new
                    {
                        a.Id,
                        a.Title,
                        a.Content,
                        a.UserId,
                        Author = u.Name,
                    })
                .ToArray()
                .Select(a => new ArticleViewModel(a.Id, a.Title, a.Content, a.Author, a.UserId));
            
            return new JsonResult(articles);
        }

        [HttpPost("~/api/article_create")]
        [Authorize]
        public IActionResult Post()
        {
            var username = User.Identity?.Name;
            if (username == null)
            {
                return Unauthorized();
            }

            var user = dbContext.Users.FirstOrDefault(u => u.Name == username);
            if (user == null)
            {
                return Unauthorized();
            }

            var article = new Article
            {
                Title = Request.Form["title"],
                Content = Request.Form["content"],
                UserId = user.Id,
            };

            dbContext.Articles.Add(article);
            dbContext.SaveChanges();

            logger.LogInformation($"User {user.Id} created article {article.Title}");

            return Redirect($"/articles/{article.Id}");
        }

        [HttpPost("~/api/search")]
        public IActionResult Search([FromBody]SearchViewModel searchView)
        {
            // Found on stack overflow. How it works? Magic, lol
            using var dbConnection = dbContext.Database.GetDbConnection();
            dbConnection.Open();
            
            using var command = dbConnection.CreateCommand();
            command.CommandText = string.Format(searchSql, searchView.Search);
            using var reader = command.ExecuteReader();

            var articles = new List<ArticleViewModel>();
            
            while (reader.Read())
            {
                var id = reader.GetInt64(0);
                var title = reader.GetString(1);
                var content = reader.GetString(2);
                var userId = reader.GetInt64(3);
                var author = reader.GetString(4);

                articles.Add(new ArticleViewModel(id, title, content, author, userId));
            }

            return new JsonResult(articles);
        }
    }
}
