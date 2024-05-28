
import subprocess

def convert_to_png(image_path):
    return subprocess.run(f'magick mogrify -format png {image_path}', shell=True)