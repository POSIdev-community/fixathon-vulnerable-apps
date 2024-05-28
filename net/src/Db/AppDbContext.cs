using Microsoft.EntityFrameworkCore;

namespace App.Db
{
    public class AppDbContext : DbContext
    {
        public AppDbContext(DbContextOptions<AppDbContext> options) :
            base(options)
        {
        }

        public DbSet<User> Users { get; set; }

        public DbSet<Article> Articles { get; set; }
    }
}
