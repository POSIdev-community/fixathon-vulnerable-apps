<?php

namespace Tests\Feature\Phd_tests;

use Illuminate\Foundation\Testing\DatabaseMigrations;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class ModelTest extends TestCase
{
    use DatabaseMigrations;

    protected $seed = true;

    public function test1_get_all_articles(): void
    {
        $response = $this->get('/api/articles');
        $response->assertSimilarJson(
            [
                [
                    "title" => "Исследование Галактики Андромеды",
                    "content" => "Наше путешествие в Галактику Андромеды раскрыло захватывающие виды и таинственные явления.",
                    "author" => "GalacticExplorer",
                    "authorProfileUrl" => "/profile/1"
                ],
                [
                    "title" => "Открытие пришельцев",
                    "content" => "Астрономы обнаружили сигналы из далекой звездной системы, возможно, указывающие на присутствие интеллектуальных форм жизни.",
                    "author" => "StellarTraveler",
                    "authorProfileUrl" => "/profile/2"
                ],
                [
                    "title" => "Раскрытие тайн черных дыр",
                    "content" => "Новые наблюдения проливают свет на загадочную природу черных дыр, вызывая сомнения в нашем понимании Вселенной.",
                    "author" => "CosmicAdventurer",
                    "authorProfileUrl" => "/profile/3"
                ],
                [
                    "title" => "Путешествие к краю Вселенной",
                    "content" => "Отправляйтесь в эпическое путешествие к самым дальним пределам космоса, где время и пространство искривляются в невообразимых масштабах.",
                    "author" => "GalacticExplorer",
                    "authorProfileUrl" => "/profile/1"
                ],
                [
                    "title" => "Поиск экзопланет",
                    "content" => "Ученые выявили перспективного кандидата на обитаемую экзопланету, обращающуюся вокруг близкой звезды, разжигая надежды на обнаружение внеземной жизни.",
                    "author" => "StellarTraveler",
                    "authorProfileUrl" =>  "/profile/2"
                ]
            ]);
    }

    public function test2_add_new_article(): void
    {
        $this->assertDatabaseMissing('articles', ['articleId' => 6, 'title' => 'new_title', 'content' => 'new_content', 'userId' => 1]);

        $this->post('/api/auth/login', ['username' => 'GalacticExplorer', 'password' => 'explorer123']);
        $this->post('/api/auth/article_create', ['title' => 'new_title', 'content' => 'new_content']);

        $this->assertDatabaseHas('articles', ['articleId' => 6, 'title' => 'new_title', 'content' => 'new_content', 'userId' => 1]);
    }

    public function test3_search_by_keyword(): void
    {
        $response = $this->postJson(
            '/api/search',
            ['search' => 'путешествие']
        );

        $response->assertSimilarJson(
            [
                [
                    "title" => "Исследование Галактики Андромеды",
                    "content" => "Наше путешествие в Галактику Андромеды раскрыло захватывающие виды и таинственные явления.",
                    "author" => "GalacticExplorer",
                    "id" => 1
                ],
                [
                    "title" => "Путешествие к краю Вселенной",
                    "content" => "Отправляйтесь в эпическое путешествие к самым дальним пределам космоса, где время и пространство искривляются в невообразимых масштабах.",
                    "author" => "GalacticExplorer",
                    "id" => 4
                ]
            ]);
    }
}
