from io import BytesIO
from uuid import uuid4

from PIL import Image

from django.core.files.base import ContentFile
from django.conf import settings
from django.views.decorators.cache import cache_page as dj_cache


def ordereddict_to_dict(value):
    for k, v in value.items():
        if isinstance(v, dict):
            value[k] = ordereddict_to_dict(v)
    return dict(value)

def thumbnail_file_name_by_orginal_name(orginal_name, thumb_name):
    """
    :param orginal_name: String
    :param thumb_name: String
    :return: String
    ex:
    ("kadir.jpg", "20x30")  # kadir_20x30.jpg
    """
    pure_name = "".join(orginal_name.split(".")[0:-1])
    ext_name = orginal_name.split(".")[-1]
    return "{0}_{1}.{2}".format(pure_name, thumb_name, ext_name)


def generate_file_name(instance, *args, **kwargs):
    ext = instance.image.name.split(".")[-1]
    return "{0}.{1}".format(uuid4().hex, ext)


def compress_image(image, do_square=False):
    """
    :param image: ContentFile
    :param do_square: boolean
    :return: ContentFile

    This function compress the image to jpeg and return new image
    if do_square==True then image will resize the square version
    """
    not_saved_pure_name = "".join(image.name.split('.')[0:-1])
    pil_image = Image.open(image)
    pil_image.convert('RGB')

    if do_square:
        edge_size = min(pil_image.size)
        pil_image = pil_image.crop((0, 0, edge_size, edge_size))

    f = BytesIO()

    fill_color = '#ffffff'
    if pil_image.mode in ('RGBA', 'LA'):
        background = Image.new(pil_image.mode[:-1], pil_image.size, fill_color)
        background.paste(pil_image, pil_image.split()[-1])
        pil_image = background

    pil_image.save(f, "JPEG", quality=90)
    image = ContentFile(f.getvalue())
    image.name = "{0}.{1}".format(not_saved_pure_name, "jpeg")


    return image


def cache_page(timeout=settings.CACHE_TTL):
    return dj_cache(timeout)
