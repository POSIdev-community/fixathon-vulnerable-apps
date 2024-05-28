<!DOCTYPE html>
<html>
<head>
    <title>Create article form</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/cosmicjs/5.6.0/cosmic.min.css">
</head>
<body>
    <h2>Article creation form</h2>
    <form action="/api/auth/article_create" method="post">
        <input type="hidden" name="redirect_to" value="/">
        <label for="title">Title:</label>
        <input type="text" id="title" name="title" required><br><br>
        <label for="content">Content:</label>
        <input type="text" id="content" name="content" required><br><br>
        <input type="submit" value="Create article">
    </form>
</body>
</html>
