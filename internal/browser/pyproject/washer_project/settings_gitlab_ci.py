from washer_project.settings import *

ES_HOST = {"host": "elasticsearch-unit", "port": 9200}
ES_STORE_INDEX = 'test_stores'
ES_RESERVATION_INDEX = 'test_reservations'

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': 'test',
        'HOST': 'postgres',
        'USER': 'test',
        'PASSWORD': 'test',
        'PORT': 5432,
    }
}


BROKER_URL = 'redis://redis:6379'
CELERY_RESULT_BACKEND = 'redis://redis:6379'
CACHES = {
    "default": {
        "BACKEND": "django_redis.cache.RedisCache",
        "LOCATION": "redis://redis:6379/1",
        "TIMEOUT": CACHE_TTL, # 5 minutes
        "OPTIONS": {
            "CLIENT_CLASS": "django_redis.client.DefaultClient"
        },
        "KEY_PREFIX": "_aracyika_"
    }
}