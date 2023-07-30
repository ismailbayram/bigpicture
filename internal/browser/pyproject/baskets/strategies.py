from baskets.models import DiscountItem, BasketItem
from reservations.enums import ReservationStatus


class BasePromotionStrategy:
    def __init__(self, campaign, basket):
        self.basket = basket
        self.campaign = campaign
        self.remaining = 0

    def check(self):
        raise NotImplementedError()

    def get_discount_amount(self, basket_item):
        raise NotImplementedError()

    def apply(self):
        raise NotImplementedError()


class OneFreeInNineStrategy(BasePromotionStrategy):
    def check(self):
        res_count = self.basket.customer_profile.reservation_set.filter(status=ReservationStatus.completed).count()
        if res_count != 0 and res_count % 9 == 0:
            return True
        self.remaining = 9 - (res_count % 9)
        return False

    def get_discount_amount(self, basket_item):
        product = basket_item.product
        return product.price(basket_item.basket.car.car_type)

    def apply(self):
        try:
            basket_item = self.basket.basketitem_set.get(product__is_primary=True)
            if basket_item.discount_item:
                return
        except BasketItem.DoesNotExist:
            return
        DiscountItem.objects.create(campaign=self.campaign, basket=self.basket, basket_item=basket_item,
                                    amount=self.get_discount_amount(basket_item))
