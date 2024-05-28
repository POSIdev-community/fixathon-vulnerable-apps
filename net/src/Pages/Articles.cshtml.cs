using App.Db;
using App.Models;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;

namespace App.Pages
{
    public class ArticlesModel : PageModel
    {
        private readonly AppDbContext dbContext;
        private readonly ILogger<ArticlesModel> logger;

        public ArticleViewModel Article { get; private set; }

        public ArticlesModel(AppDbContext dbContext, ILogger<ArticlesModel> logger)
        {
            this.dbContext = dbContext;
            this.logger = logger;
        }

        public IActionResult OnGet(long id)
        {
            var article = dbContext.Articles.Find(id);
            if (article == null)
            {
                return NotFound();
            }

            var author = dbContext.Users.Find(article.UserId);
            Article = new(article.Id, article.Title, article.Content, author.Name, article.UserId);

            return Page();
        }
    }
}
