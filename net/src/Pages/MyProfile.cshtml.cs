using System.Linq;
using App.Db;
using App.Models;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;

namespace App.Pages
{
    [Authorize]
    public class MyProfileModel : PageModel
    {
        private readonly AppDbContext dbContext;
        private readonly ILogger<ArticlesModel> logger;

        public ProfileViewModel Profile { get; private set; }

        public MyProfileModel(AppDbContext dbContext, ILogger<ArticlesModel> logger)
        {
            this.dbContext = dbContext;
            this.logger = logger;
        }

        public IActionResult OnGet()
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

            var articles = dbContext.Articles.Where(a => a.UserId == user.Id).ToArray();

            Profile = new ProfileViewModel
            {
                Name = user.Name,
                Articles = articles
                    .Select(a => new ArticleViewModel(a.Id, a.Title, a.Content, user.Name, user.Id))
                    .ToArray(),
            };

            return Page();
        }
    }
}
