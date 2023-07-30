from django.utils.translation import ugettext_lazy as _

from enumfields import Enum


class CarType(Enum):
    sedan = 'sedan'
    hatchback = 'hatchback'
    suv = 'suv'
    commercial = 'commercial'

    class Labels:
        sedan = _('Sedan')
        hatchback = _('Hatchback')
        suv = _('Suv')
        commercial = _('Ticari Araç')
