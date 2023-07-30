import logging
from django.utils.timezone import datetime
from django.conf import settings
from django.core.cache import cache
from django.utils.encoding import force_str

from washer_project.celery import app

logger = logging.getLogger(__name__)


class LockTask(app.Task):
    """this abstract class ensures the same tasks run only once at a time"""
    abstract = True
    TTL = getattr(settings, 'DEFAULT_TASK_LOCK_TTL', 60 * 15)

    def __init__(self, *args, **kwargs):
        super(LockTask, self).__init__(*args, **kwargs)

    def generate_lock_cache_key(self, *args, **kwargs):
        args_key = [force_str(arg) for arg in args]
        kwargs_key = ['{}_{}'.format(k, force_str(v)) for k, v in
                      sorted(kwargs.items())]
        return '_'.join([self.name] + args_key + kwargs_key)

    def __call__(self, *args, **kwargs):
        """check task"""
        lock_cache_key = (self.request.headers or {}).pop('cache_key', None)
        if not lock_cache_key:
            lock_cache_key = self.generate_lock_cache_key(*args, **kwargs)
        lock_time = datetime.now().isoformat()
        lock_acquired = cache.set(lock_cache_key, lock_time, nx=True,
                                  timeout=self.TTL)

        if lock_acquired:
            try:
                return self.run(*args, **kwargs)
            finally:
                cache.delete(lock_cache_key)
        else:
            logger.info('Task %s is already running..' % self.name)
