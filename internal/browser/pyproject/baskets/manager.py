from django.db.models.manager import Manager


class CampaignManager(Manager):
    def actives(self):
        return self.filter(is_active=True)
