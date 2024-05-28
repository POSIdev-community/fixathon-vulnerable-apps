<!--This page is for filling from server-->
<!DOCTYPE html>
<html>
<head>
    <title>Article Template</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/cosmicjs/5.6.0/cosmic.min.css">
</head>
<body>
    <article>
        <h1>{{ $title }}</h1>
        <p><a href="/profile/{{ $author_id }}">{{ $author }}</a></p>
        <p>{{ $content }}</p>
    </article>
</body>
</html>
