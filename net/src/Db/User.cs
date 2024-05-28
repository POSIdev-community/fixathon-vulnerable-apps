using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace App.Db
{
    public class User
    {
        [Key]
        [Required]
        [Column("userId")]
        public long Id { get; set; }

        [Required]
        [Column("username")]
        public string Name { get; set; }

        [Required]
        [Column("password")]
        public string Password { get; set; }
    }
}
