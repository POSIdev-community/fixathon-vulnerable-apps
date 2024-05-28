<?php

namespace App\Http\Controllers;

use PDO;

class ModelController extends Controller
{
    public static function add_article($title, $content, $userId)
    {
        self::query("insert into articles (title, content, userId) values ('$title', '$content', $userId);", execute: false);
    }
    public static function get_all_articles()
    {
        return self::query('select * from articles');
    }

    public static function get_articles_by_user_id($user_id)
    {
        return self::query('select * from articles where userId = :userId',
            ['userId' => $user_id]);
    }

    public static function get_article_by_article_id($article_id)
    {
        return self::query('select * from articles where articleId = :articleId',
            [':articleId' => $article_id]);
    }

    public static function get_articles_by_keyword($keyword)
    {
        return self::query('select * from articles where title like :keyword or content like :keyword',
            [':keyword' => "% $keyword %"]);
    }

    public static function get_user_by_id($id)
    {
        return self::query('select * from users where id = :userId',
            ['userId' => $id]);
    }

    public static function get_user_id_by_username($username)
    {
        return self::query('select id from users where username = :username',
            [':username' => $username]);
    }

    private static function query($query, array $bindings = [], bool $execute = true)
    {
        $connection = new PDO(implode(DIRECTORY_SEPARATOR, ['sqlite:.', 'database', 'database.sqlite']));
        $statement = $connection->query($query);
        if($execute) {
            $statement->execute($bindings);
        }
        return $statement->fetchAll(PDO::FETCH_CLASS);
    }
}
