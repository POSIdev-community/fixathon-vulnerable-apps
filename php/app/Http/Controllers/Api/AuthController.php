<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\ModelController;
use Illuminate\Routing\Controller;
use Illuminate\Support\Facades\Cookie;
use Illuminate\Support\Facades\Session;

class AuthController extends Controller
{
    public function __construct()
    {
        $this->middleware('auth:api')->except('login');
    }

    public function login()
    {
        $credentials = request(['username', 'password']);

        if (!array_key_exists('username', $credentials) || !array_key_exists('password', $credentials)) {
            return response()->json(['error' => 'Username and password expected'], 400);
        }

        if (!$token = auth()->attempt($credentials)) {
            return response()->json(['error' => 'Wrong credentials'], 401);
        }

        session_start();

        $_SESSION['username'] = $credentials['username'];
        $user = ModelController::get_user_id_by_username($credentials['username']);
        $_SESSION['userId'] = array_shift($user)->id;

        $redirect_to = request('redirect_to');
        if(empty($redirect_to)){
            $redirect_to = '/';
        }
        header('Location:' . $redirect_to);

        return response(status: 302)->cookie('jwt', $token, 5);
    }

    public function logout()
    {
        auth()->logout();

        Session::remove('username');
        Session::remove('userId');
        Cookie::forget('jwt');

        return response()->json(['message' => 'Successfully logged out']);
    }
}
