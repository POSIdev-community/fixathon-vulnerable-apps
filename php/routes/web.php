<?php

use App\Http\Controllers\Api\ProfileController;
use App\Http\Controllers\Api\ArticleController;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return view('index');
});

Route::get('/login', function () {
    return view('login');
})->name('login');

Route::get('/search', function () {
    return view('search');
});

Route::get('/article_create', function () {
    return view('article_create');
});

Route::get('/profile/{id}', [ProfileController::class, 'load_profile_by_id']);
Route::get('/article/{id}', [ArticleController::class, 'load_article_by_id']);
