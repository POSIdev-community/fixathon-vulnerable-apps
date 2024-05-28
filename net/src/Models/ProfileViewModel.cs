namespace App.Models
{
    public class ProfileViewModel
    {
        public long Id { get; init; }

        public string Name { get; init; }

        public ArticleViewModel[] Articles { get; init; }
    }
}
