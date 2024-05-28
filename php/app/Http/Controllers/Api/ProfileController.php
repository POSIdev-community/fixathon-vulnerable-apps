<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\ModelController;
use Illuminate\Routing\Controller;
use Illuminate\Support\Facades\URL;
use Illuminate\Support\Facades\View;

class ProfileController extends Controller
{
    public function __construct()
    {
        $token = request()->cookie('jwt');
        if(!empty($token)){
            request()->headers->set('Authorization', 'Bearer '. $token);
        }
        $this->middleware('auth:api')->except('load_profile_by_id');
    }

    public static function load_profile_by_id($id)
    {
        $author = ModelController::get_user_by_id($id);

        if(empty($author)){
            error_log("Unable to load profile by id $id", 3, '.' . DIRECTORY_SEPARATOR . 'phdays_log.txt');
            return response()->json(['error' => 'Profile not found'], 404);
        }
        $author = array_shift($author)->username;

        $articles = ModelController::get_articles_by_user_id($id);

        $articles_save_data = [];
        if(!empty($articles))
        {
            foreach($articles as $article)
            {
                array_push($articles_save_data, [
                    'title' => $article->title,
                    'content' => $article->content,
                ]);
            }
        }

        return View::make('author_template')
            ->with('author', $author)
            ->with('articles', $articles_save_data);
    }

    public static function load_my_profile()
    {
        if(session_status() != PHP_SESSION_ACTIVE) {
            session_start();
        }

        $author = $_SESSION['username'];
        $user_id = $_SESSION['userId'];

        $articles = ModelController::get_articles_by_user_id($user_id);
        $articles_save_data = [];
        if(!empty($articles))
        {
            foreach($articles as $article)
            {
                array_push($articles_save_data, [
                    'title' => $article->title,
                    'content' => $article->content,
                ]);
            }
        }

        $path_to_photo = URL::asset(PhotoController::get_path_to_photo_by_id($user_id));

        return View::make('my_profile', [
            'author' => $author,
            'articles' => $articles_save_data,
            'path_to_photo' => $path_to_photo]);
    }
}
