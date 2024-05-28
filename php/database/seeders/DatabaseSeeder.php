<?php

namespace Database\Seeders;

use App\Models\Article;
use App\Models\User;
use Illuminate\Database\Seeder;

class DatabaseSeeder extends Seeder
{
    /**
     * Seed the application's database.
     */
    public function run(): void
    {
        User::factory()->create([
            'username' => 'GalacticExplorer',
            'password' => 'explorer123',
        ]);

        User::factory()->create([
            'username' => 'StellarTraveler',
            'password' => 'traveler456',
        ]);

        User::factory()->create([
            'username' => 'CosmicAdventurer',
            'password' => 'adventurer789',
        ]);

        Article::factory()->create([
            'title' => 'Исследование Галактики Андромеды',
            'content' => 'Наше путешествие в Галактику Андромеды раскрыло захватывающие виды и таинственные явления.',
            'userId' => 1,
        ]);

        Article::factory()->create([
            'title' => 'Открытие пришельцев',
            'content' => 'Астрономы обнаружили сигналы из далекой звездной системы, возможно, указывающие на присутствие интеллектуальных форм жизни.',
            'userId' => 2,
        ]);

        Article::factory()->create([
            'title' => 'Раскрытие тайн черных дыр',
            'content' => 'Новые наблюдения проливают свет на загадочную природу черных дыр, вызывая сомнения в нашем понимании Вселенной.',
            'userId' => 3,
        ]);

        Article::factory()->create([
            'title' => 'Путешествие к краю Вселенной',
            'content' => 'Отправляйтесь в эпическое путешествие к самым дальним пределам космоса, где время и пространство искривляются в невообразимых масштабах.',
            'userId' => 1,
        ]);

        Article::factory()->create([
            'title' => 'Поиск экзопланет',
            'content' => 'Ученые выявили перспективного кандидата на обитаемую экзопланету, обращающуюся вокруг близкой звезды, разжигая надежды на обнаружение внеземной жизни.',
            'userId' => 2,
        ]);
    }
}
