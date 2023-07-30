from rest_framework import status
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.response import Response

from api.permissions import HasGroupPermission
from users.enums import GroupType
from baskets.enums import PromotionType
from baskets.models import Campaign
from baskets.service import BasketService
from baskets.resources.serializers import (BasketSerializer, CreateBasketItemSerializer,
                                           CampaignSerializer)
from reservations.enums import ReservationStatus


class BasketViewSet(viewsets.ViewSet):
    permission_classes = (HasGroupPermission,)
    permission_groups = {
        'view_basket': [GroupType.customer],
        'add_item': [GroupType.customer],
        'delete_item': [GroupType.customer],
    }
    service = BasketService()

    @action(methods=['GET'], detail=False)
    def view_basket(self, request, *args, **kwargs):
        basket = self.service.get_or_create_basket(request.user.customer_profile)
        self.service.apply_discounts(basket)
        serializer = BasketSerializer(instance=basket)
        return Response({'basket': serializer.data}, status=status.HTTP_200_OK)

    @action(methods=['POST'], detail=False)
    def add_item(self, request, *args, **kwargs):
        serializer = CreateBasketItemSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        basket = self.service.get_or_create_basket(request.user.customer_profile)
        self.service.add_basket_item(basket, **serializer.validated_data)
        self.service.apply_discounts(basket)
        serializer = BasketSerializer(instance=basket)
        return Response({'basket': serializer.data}, status=status.HTTP_201_CREATED)

    @action(methods=['DELETE'], detail=False)
    def delete_item(self, request, *args, **kwargs):
        serializer = CreateBasketItemSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        basket = self.service.get_or_create_basket(request.user.customer_profile)
        self.service.delete_basket_item(basket, **serializer.validated_data)
        self.service.apply_discounts(basket)
        serializer = BasketSerializer(instance=basket)
        return Response({'basket': serializer.data}, status=status.HTTP_204_NO_CONTENT)


class CampaignViewSet(viewsets.ReadOnlyModelViewSet):
    queryset = Campaign.objects.filter(is_active=True)
    serializer_class = CampaignSerializer
    permission_classes = (HasGroupPermission,)
    permission_groups = {
        'retrieve': [GroupType.customer],
    }

    def retrieve(self, request, *args, **kwargs):
        campaign = self.get_object()
        customer_profile = request.user.customer_profile
        res_count = customer_profile.reservation_set.filter(status=ReservationStatus.completed).count() % 9
        if campaign.promotion_type == PromotionType.one_free_in_nine:
            return Response({'completed': res_count})
