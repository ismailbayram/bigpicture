from rest_framework.serializers import ModelSerializer

from cars.models import Car
from cars.enums import CarType
from api.fields import EnumField


class CarSerializer(ModelSerializer):
    car_type = EnumField(enum=CarType)

    class Meta:
        model = Car
        fields = ('pk', 'licence_plate', 'car_type', 'is_selected', 'is_active')
