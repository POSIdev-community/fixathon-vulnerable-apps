using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net.Http;
using System.Threading.Tasks;
using App.Db;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Microsoft.VisualBasic.FileIO;

namespace App.Controllers
{
    [ApiController]
    [Authorize]
    [Route("api/profile")]
    public class ProfileController : ControllerBase
    {
        private readonly AppDbContext dbContext;
        private readonly ILogger<ProfileController> logger;

        public ProfileController(AppDbContext dbContext, ILogger<ProfileController> logger)
        {
            this.dbContext = dbContext;
            this.logger = logger;
        }

        [HttpPost("upload_photo_url")]
        public async Task<IActionResult> UploadPhotoUrl()
        {
            var username = User.Identity?.Name;
            if (username == null) { return Unauthorized(); }

            var user = dbContext.Users.FirstOrDefault(u => u.Name == username);
            if (user == null) { return Unauthorized(); }

            string photoUrl = Request.Form["profile-photo-url"];
            //https://uxwing.com/wp-content/themes/uxwing/download/peoples-avatars/corporate-user-icon.png

            using var client = new HttpClient();
            var response = await client.GetAsync(photoUrl);

            if (!response.IsSuccessStatusCode)
            {
                return BadRequest(await response.Content.ReadAsStringAsync());
            }

            var fileStream = await response.Content.ReadAsStreamAsync();
            SaveFile(fileStream, user.Id);
            return Redirect("/myprofile");
        }

        [HttpPost("upload_photo")]
        public IActionResult UploadPhoto(IFormFile file)
        {
            var username = User.Identity?.Name;
            if (username == null) { return Unauthorized(); }

            var user = dbContext.Users.FirstOrDefault(u => u.Name == username);
            if (user == null) { return Unauthorized(); }

            var tempFileName = Path.Combine(SpecialDirectories.Temp, file.Name);
            ProcessImage(tempFileName);
            using var fileStream = System.IO.File.OpenRead(tempFileName);
            SaveFile(fileStream, user.Id);

            return Redirect("/myprofile");
        }

        private void ProcessImage(string filePath)
        {
            var process = new ProcessStartInfo()
            {
                FileName =
                    $"cmd.exe magick mogrify -format jpg -quality 50\" //\"magick mogrify -resize 50%% {filePath}"
            };

            Process.Start(process)?.WaitForExit();
        }

        private void SaveFile(Stream stream, long userId)
        {
            using var photo = System.IO.File.Create(Path.Combine("static", $"profile_photo{userId}.png"));
            stream.CopyTo(photo);
        }
    }
}
