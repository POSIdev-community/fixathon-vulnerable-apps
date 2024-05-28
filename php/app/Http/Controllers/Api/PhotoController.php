<?php

namespace App\Http\Controllers\Api;

use Illuminate\Routing\Controller;

class PhotoController extends Controller
{
    public function __construct()
    {
        $token = request()->cookie('jwt');
        if(!empty($token)){
            request()->headers->set('Authorization', 'Bearer '.$token);
        }
        $this->middleware('auth:api')->except('get_path_to_photo_by_id');
    }

    public static function get_path_to_photo_by_id($id)
    {
        $path_to_custom_photo = implode(DIRECTORY_SEPARATOR, ['.', 'Photos', 'profile_photo' . $id . '.jpg']);
        $path_to_default_photo = implode(DIRECTORY_SEPARATOR, ['.', 'Photos', 'default_profile_photo.jpg']);
        return file_exists($path_to_custom_photo) ? $path_to_custom_photo : $path_to_default_photo;
    }

    public static function upload_photo_by_url()
    {
        $photo_url = request()->post('profile-photo-url');
        if(empty($photo_url)){
            return response()->json(['error' => 'No photo url provided'], 400);
        }

        if(session_status() != PHP_SESSION_ACTIVE) {
            session_start();
        }

        $file_name = 'profile_photo' . $_SESSION['userId'] . '.jpg';
        $new_path_to_photo = '.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . $file_name;

        $file = fopen($new_path_to_photo, 'w+');
        fwrite($file, file_get_contents($photo_url));
        fclose($file);

        self::convert_image_to_jpg($new_path_to_photo, $file_name);

        header('Location:api/auth/my_profile');

        return response(status: 302);
    }

    public static function upload_photo()
    {
        $file_data = request()->file('profile-photo');
        if(empty($file_data)){
            return response()->json(['error' => 'No photo provided'], 400);
        }

        if(session_status() != PHP_SESSION_ACTIVE) {
            session_start();
        }

        $file_name = 'profile_photo' . $_SESSION['userId'] . '.jpg';
        $new_path_to_photo = '.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . $file_name;

        rename($file_data->getPathname(), $new_path_to_photo);
        self::convert_image_to_jpg($new_path_to_photo, $file_name);

        header('Location:api/auth/my_profile');

        return response(status: 302);
    }

    private static function convert_image_to_jpg($path_to_image, $image_name)
    {
        $path_to_new_image = str_replace($image_name, 'new_' . $image_name, $path_to_image);

        //TODO: remove 'magick ' word from lines bellow if tests failed
        exec('magick mogrify -format jpg ' . $path_to_image .
            '&& magick convert ' . $path_to_image . ' -resize 200x200 ' . $path_to_new_image);
    }
}
