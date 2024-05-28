<?php

use App\Http\Controllers\Api\ArticleController;
use App\Http\Controllers\Api\AuthController;
use App\Http\Controllers\Api\PhotoController;
use App\Http\Controllers\Api\ProfileController;
use Illuminate\Support\Facades\Route;

Route::get('articles', [ArticleController::class, 'get_all_articles']);
Route::post('search', [ArticleController::class, 'get_articles_by_keyword']);

Route::post('auth/login', [AuthController::class, 'login']);
Route::post('auth/logout', [AuthController::class, 'logout']);
Route::get('auth/my_profile', [ProfileController::class, 'load_my_profile']);

Route::post('auth/profile/upload_photo_url', [PhotoController::class, 'upload_photo_by_url']);
Route::post('auth/profile/upload_photo', [PhotoController::class, 'upload_photo']);

Route::post('auth/article_create', [ArticleController::class, 'create_article']);
