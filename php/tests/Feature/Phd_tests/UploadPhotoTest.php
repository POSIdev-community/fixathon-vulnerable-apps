<?php

namespace Tests\Feature\Phd_tests;

use Illuminate\Foundation\Testing\DatabaseMigrations;
use Illuminate\Http\UploadedFile;
use Illuminate\Testing\Assert;
use Tests\TestCase;

class UploadPhotoTest extends TestCase
{
    use DatabaseMigrations;

    protected $seed = true;

    public function test1_photo_url()
    {
        session_abort();

        $response = $this->post('/api/auth/login', ['username' => 'GalacticExplorer', 'password' => 'explorer123']);

        $response->assertStatus(302);

        $response = $this->post('/api/auth/profile/upload_photo_url', ['profile-photo-url' => 'https://upload.wikimedia.org/wikipedia/en/8/86/Avatar_Aang.png']);

        $response->assertStatus(302);

        Assert::assertFileExists('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'profile_photo1.jpg');

        $old_photo_size = getimagesize('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'profile_photo1.jpg');
        $new_photo_size = getimagesize('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'new_profile_photo1.jpg');

        unlink('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'profile_photo1.jpg');
        unlink('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'new_profile_photo1.jpg');

        self::assertTrue($old_photo_size != $new_photo_size);
    }

    public function test2_photo_file()
    {
        session_abort();

        $response = $this->post('/api/auth/login', ['username' => 'GalacticExplorer', 'password' => 'explorer123']);

        $response->assertStatus(302);

        $file = UploadedFile::fake()->image('./photos/photo_to_upload.png');

        $response = $this->post('/api/auth/profile/upload_photo', ['profile-photo' => $file]);

        $response->assertStatus(302);

        Assert::assertFileExists('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'profile_photo1.jpg');

        $old_photo_size = getimagesize('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'profile_photo1.jpg');
        $new_photo_size = getimagesize('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'new_profile_photo1.jpg');

        unlink('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'profile_photo1.jpg');
        unlink('.' . DIRECTORY_SEPARATOR . 'public' . DIRECTORY_SEPARATOR . 'Photos' . DIRECTORY_SEPARATOR . 'new_profile_photo1.jpg');

        self::assertTrue($old_photo_size != $new_photo_size);
    }
}
