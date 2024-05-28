<!DOCTYPE html>
<html>
<head>
    <title>Articles Search</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/cosmicjs/5.6.0/cosmic.min.css">
</head>
<body>
    <header class="cosmic-header">
        <h1>Cosmic Blogs</h1>
    </header>

    <main>
        <h1>Articles Search</h1>
        <form>
            <input type="text" name="search" placeholder="Search articles...">
            <button type="submit">Search</button>
        </form>

        <div id="search-results">
            <!-- Display search results here -->
        </div>
    </main>

    <script>
        document.querySelector('form').addEventListener('submit', async (event) => {
            event.preventDefault();

            const searchInput = document.querySelector('input[name="search"]');
            const searchValue = searchInput.value;

            try {
                const response = await fetch('/api/search', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ search: searchValue })
                });

                const data = await response.json();

                const searchResults = document.getElementById('search-results');
                searchResults.innerHTML = '';

                data.forEach(article => {
                    const articleElement = document.createElement('div');
                    articleElement.innerHTML = `
                        <h2><a href="/article/${article.id}">${article.title}</a></h2>
                        <p><strong>Author:</strong> ${article.author}</p>
                        <p>${article.content}</p>
                    `;
                    searchResults.appendChild(articleElement);
                });
            } catch (error) {
                console.error('Error:', error);
            }
        });

        const jwtCookie = document.cookie.split(';').find(cookie => cookie.trim().startsWith('jwt='));
        const logoutButton = document.createElement('button');
        logoutButton.textContent = 'Logout';

        if (jwtCookie) {
            logoutButton.addEventListener('click', async () => {
                try {
                    const response = await fetch('/api/logout', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    });

                    if (response.ok) {
                        // Redirect to home page
                        window.location.href = '/';
                    } else {
                        console.error('Logout failed');
                    }
                } catch (error) {
                    console.error('Error:', error);
                }
            });
            document.body.appendChild(logoutButton);
        }
    </script>
    </script>
    </main>

    <footer>
        <!-- Add your footer content here -->
    </footer>

    <script>
        // Add your JavaScript code here
    </script>
</body>
</html>
