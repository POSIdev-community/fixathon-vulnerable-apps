using System.Linq;
using App.Db;
using App.Models;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.RazorPages;
using Microsoft.Extensions.Logging;

namespace App.Pages
{
    public class ProfilesModel : PageModel
    {
        private readonly AppDbContext dbContext;
        private readonly ILogger<ProfilesModel> logger;

        public ProfilesModel(AppDbContext dbContext, ILogger<ProfilesModel> logger)
        {
            this.dbContext = dbContext;
            this.logger = logger;
        }

        public ProfileViewModel Profile { get; private set; }

        public IActionResult OnGet(long id)
        {
            var author = dbContext.Users.Find(id);
            if (author == null)
            {
                return NotFound();
            }

            var articles = dbContext.Articles.Where(a => a.UserId == id).ToArray();

            Profile = new ProfileViewModel
            {
                Id = id,
                Name = author.Name,
                Articles = articles
                    .Select(a => new ArticleViewModel(a.Id, a.Title, a.Content, author.Name, id))
                    .ToArray(),
            };

            return Page();
        }
    }
}