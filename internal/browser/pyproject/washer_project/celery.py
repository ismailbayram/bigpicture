from __future__ import absolute_import
import os
from celery import Celery
from django.conf import settings

from washer_project.celery_schedule_conf import CELERYBEAT_SCHEDULE

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'washer_project.settings')
app = Celery('washer_project')

app.config_from_object('django.conf:settings')
app.autodiscover_tasks(lambda: settings.INSTALLED_APPS)

app.conf.beat_schedule = CELERYBEAT_SCHEDULE
