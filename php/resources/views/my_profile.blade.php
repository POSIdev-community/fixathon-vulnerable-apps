<!DOCTYPE html>
<html>
<head>
    <title>My Profile</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/cosmicjs/5.6.0/cosmic.min.css">
</head>
<body>

    <h1>{{ $author }}</h1>
    <h2>Profile Photo</h2>
    <img src="{{ $path_to_photo }}" alt="Profile Photo">

    <form enctype="multipart/form-data" action="/api/auth/profile/upload_photo" method="post">
        <label for="profile-photo">Upload Profile Photo:</label>
        <input type="file" id="profile-photo" name="profile-photo">
        <button type="submit">Upload</button>
    </form>

    <form action="/api/auth/profile/upload_photo_url" method="post">
        <label for="profile-photo-url">Upload Profile Photo from URL:</label>
        <input type="url" id="profile-photo-url" name="profile-photo-url">
        <button type="submit">Upload</button>
    </form>

    @foreach ($articles as $article)
        <li><b>{{ $article['title'] }}</b></li>
        <p>{{ $article['content'] }}</p>
    @endforeach

</body>
</html>
