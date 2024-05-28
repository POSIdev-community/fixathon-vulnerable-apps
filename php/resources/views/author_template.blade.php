<!--This page is for filling from server-->
<!DOCTYPE html>
<html>
<head>
    <title>Article Template</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/cosmicjs/5.6.0/cosmic.min.css">
</head>
<body>
    <article>
        <h1>{{ $author }}</h1>
        @foreach ($articles as $article)
            <li><b>{{ $article['title'] }}</b></li>
            <p>{{ $article['content'] }}</p>
        @endforeach
    </article>
</body>
</html>
