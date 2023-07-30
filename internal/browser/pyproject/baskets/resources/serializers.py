from rest_framework import serializers

from api.fields import EnumField
from baskets.models import Basket, BasketItem, DiscountItem, Campaign
from baskets.enums import PromotionType
from cars.resources.serializers import CarSerializer
from products.resources.serializers import ProductSerializer
from products.models import Product
from products.enums import Currency


class CreateBasketItemSerializer(serializers.Serializer):
    product = serializers.PrimaryKeyRelatedField(queryset=Product.objects.filter(is_active=True,
                                                                                 store__is_approved=True,
                                                                                 store__is_active=True))


class DiscountItemSerializer(serializers.ModelSerializer):
    name = serializers.CharField(source='campaign.name')

    class Meta:
        model = DiscountItem
        fields = ('name', 'amount')


class BasketItemSerializer(serializers.ModelSerializer):
    price = serializers.DecimalField(source='get_price', max_digits=6,
                                     decimal_places=2, read_only=True)
    product = ProductSerializer()

    class Meta:
        model = BasketItem
        fields = ('product', 'price', 'quantity',)


class BasketSerializer(serializers.ModelSerializer):
    currency = EnumField(enum=Currency)
    total_amount = serializers.DecimalField(source='get_total_amount', max_digits=6,
                                            decimal_places=2, read_only=True)
    total_quantity = serializers.IntegerField(source='get_total_quantity')
    net_amount = serializers.DecimalField(source='get_net_amount', max_digits=6,
                                          decimal_places=2, read_only=True)
    car = CarSerializer()
    basketitem_set = BasketItemSerializer(many=True)
    warning_messages = serializers.ListField()
    campaign_messages = serializers.ListField()
    discountitem_set = DiscountItemSerializer(many=True)

    class Meta:
        model = Basket
        fields = ('currency', 'total_amount', 'total_quantity', 'net_amount', 'car',
                  'basketitem_set', 'discountitem_set', 'warning_messages', 'campaign_messages',)


class CampaignSerializer(serializers.ModelSerializer):
    promotion_type = EnumField(enum=PromotionType)

    class Meta:
        model = Campaign
        fields = ('pk', 'name', 'promotion_type')
