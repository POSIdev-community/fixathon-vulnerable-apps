<?php

namespace Tests\Feature\Phd_tests;

use Illuminate\Foundation\Testing\DatabaseMigrations;
use Tests\TestCase;

class AuthTest extends TestCase
{
    use DatabaseMigrations;

    protected $seed = true;

    public function test1_login_success(): void
    {
        $response = $this->post('/api/auth/login', ['username' => 'GalacticExplorer', 'password' => 'explorer123']);
        $response->assertStatus(302);
    }

    public function test2_login_wrong_credentials(): void
    {
        $response = $this->post('/api/auth/login', ['username' => 'random_user', 'password' => 'random_password']);
        $response->assertStatus(401);
    }

    public function test3_login_no_credentials(): void
    {
        $response = $this->post('/api/auth/login');
        $response->assertStatus(400);
    }

    public function test4_access_my_profile_success(): void
    {
        $this->post('/api/auth/login', ['username' => 'GalacticExplorer', 'password' => 'explorer123']);
        $response = $this->get('/api/auth/my_profile');
        $response->assertStatus(200);
    }

    public function test5_logout(): void
    {
        $this->post('/api/auth/login', ['username' => 'GalacticExplorer', 'password' => 'explorer123']);
        $response = $this->get('/api/auth/my_profile');
        $response->assertStatus(200);

        $response = $this->post('/api/auth/logout');
        $response->assertSimilarJson(['message' => 'Successfully logged out']);

        $response = $this->get('/api/auth/my_profile');
        $response->assertStatus(302);
    }
}
