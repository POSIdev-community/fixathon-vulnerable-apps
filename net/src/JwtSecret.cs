using System.Text;

namespace App
{
    internal static class JwtSecret
    {
        public static byte[] Default { get; } = Encoding.UTF8.GetBytes("strong_secret_key");
    }
}
