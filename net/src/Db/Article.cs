using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace App.Db
{
    public class Article
    {
        [Key]
        [Required]
        [Column("articleId")]
        public long Id { get; set; }

        [Required]
        [Column("title")]
        public string Title { get; set; }

        [Required]
        [Column("userId")]
        public long UserId { get; set; }

        [Required]
        [Column("content")]
        public string Content { get; set; }
    }
}
