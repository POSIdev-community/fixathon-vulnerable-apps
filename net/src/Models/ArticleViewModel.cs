namespace App.Models
{

    public class ArticleViewModel
    {
        public ArticleViewModel(long id, string title, string content,
            string author, long userId)
        {
            Id = id;
            Title = title;
            Content = content;
            Author = author;
            UserId = userId;
            AuthorProfileUrl = $"/profiles/{userId}";
        }

        public long Id { get; }

        public long UserId { get; }

        public string Title { get; }

        public string Content { get; }

        public string Author { get; }

        public string AuthorProfileUrl { get; }
    }
}
