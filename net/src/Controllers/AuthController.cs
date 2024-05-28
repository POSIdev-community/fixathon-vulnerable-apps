using System;
using System.IdentityModel.Tokens.Jwt;
using System.Linq;
using System.Security.Claims;
using App.Db;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Microsoft.IdentityModel.Tokens;

namespace App.Controllers
{
    [ApiController]
    [Route("api")]
    public class AuthController : ControllerBase
    {
        private readonly AppDbContext dbContext;
        private readonly ILogger<AuthController> logger;

        public AuthController(AppDbContext dbContext, ILogger<AuthController> logger)
        {
            this.dbContext = dbContext;
            this.logger = logger;
        }

        [HttpPost("login")]
        public IActionResult Login()
        {
            string username = Request.Form["username"];
            string password = Request.Form["password"];
            string redirectTo = Request.Form["redirect_to"];

            // Validate form input
            if (string.IsNullOrWhiteSpace(username) || string.IsNullOrWhiteSpace(password))
            {
                return BadRequest(new { Error = "Parameters can't be empty" });
            }

            // Check for username and password match
            var user = dbContext.Users
                .FirstOrDefault(u => u.Name == username && u.Password == password);
            if (user == null)
            {
                return Unauthorized(new { Error = "Authentication failed" });
            }

            var token = CreateToken(user.Id, username);
            if (token is null)
            {
                return Problem(detail: "Error creating token", statusCode: 500);
            }

            logger.LogInformation("Token created: {0}", token);
            HttpContext.Response.Cookies.Append("jwt", token, new CookieOptions()
            {
                MaxAge = TimeSpan.FromHours(1),
                HttpOnly = false,
                Secure = false,
            });

            if (redirectTo is null)
            {
                return Ok();
            }

            return Redirect(redirectTo);
        }

        [HttpPost("logout")]
        public IActionResult Logout(string redirect)
        {
            HttpContext.Response.Cookies.Delete("jwt");

            if (redirect is null)
            {
                return Ok();
            }

            return Redirect(redirect);
        }

        private string CreateToken(long userId, string username)
        {
            var authClaims = new Claim[]
            {
                new(ClaimTypes.Name, username),
                new(ClaimTypes.NameIdentifier, userId.ToString()),
                new(JwtRegisteredClaimNames.Jti, Guid.NewGuid().ToString()),
                new(JwtRegisteredClaimNames.Sub, userId.ToString()),
            };

            var authSigningKey = new SymmetricSecurityKey(JwtSecret.Default);

            var token = new JwtSecurityToken(
                issuer: "phdays-app",
                expires: DateTime.Now.AddHours(1),
                claims: authClaims,
                signingCredentials: new SigningCredentials(authSigningKey, SecurityAlgorithms.HmacSha256)
            );

            return new JwtSecurityTokenHandler().WriteToken(token);
        }
    }
}
