from celery.schedules import crontab


CELERYBEAT_SCHEDULE = {
    'check-expired-reservations': {
        'task': 'reservations.tasks.check_expired_reservations',
        'schedule': crontab(minute='*/1'),
    },
    'create-next-week-day-reservations': {
        'task': 'reservations.tasks.create_next_week_day',
        'schedule': crontab(hour='4', minute='0'),
    },
    'notify_stores_for_increasing': {
        'task': 'store.tasks.notify_stores_for_increasing',
        'schedule': crontab(hour='12', minute='0'),
        # her gün 12 de
    },
}
