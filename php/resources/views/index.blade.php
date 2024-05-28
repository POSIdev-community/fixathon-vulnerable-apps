<!DOCTYPE html>
<html>
<head>
    <title>My Blog</title>
    <style>
        .login-button {
            position: absolute;
            top: 10px;
            right: 10px;
        }
    </style>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/cosmicjs/5.6.0/cosmic.min.css">
</head>
<body>
    <header class="cosmic-header">
        <h1>Cosmic Blogs</h1>
    </header>
    <button class="login-button" onclick="location.href='/login'">Login</button>
    <h1>Articles</h1>
    <ul id="article-list"></ul>

    <script>
        // Fetch articles from the API
        fetch('/api/articles')
            .then(response => response.json())
            .then(articles => {
                const articleList = document.getElementById('article-list');
                articles.forEach(article => {
                    const listItem = document.createElement('li');
                    const title = document.createElement('h2');
                    title.textContent = article.title;
                    listItem.appendChild(title);

                    const author = document.createElement('p');
                    const authorLink = document.createElement('a');
                    authorLink.href = article.authorProfileUrl;
                    authorLink.textContent = article.author;
                    author.appendChild(authorLink);
                    listItem.appendChild(author);

                    const content = document.createElement('p');
                    content.textContent = article.content;
                    listItem.appendChild(content);

                    articleList.appendChild(listItem);
                });
            })
            .catch(error => {
                console.error('Error fetching articles:', error);
            });
    </script>
</body>
</body>
</html>