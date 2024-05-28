<!DOCTYPE html>
<html>
<head>
    <title>Login Form</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/cosmicjs/5.6.0/cosmic.min.css">
</head>
<body>
    <header class="cosmic-header">
        <h1>Welcome to Cosmic Blogs</h1>
    </header>
    <h2>Login</h2>
    <form action="api/auth/login" method="post">
        <input type="hidden" name="redirect_to" value="/">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" required><br><br>
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required><br><br>
        <input type="submit" value="Login">
    </form>
</body>
</html>
