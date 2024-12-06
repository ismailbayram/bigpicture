from decimal import Decimal
from django.utils.translation import ugettext_lazy as _

from django.db.models import Q, F
from django.db.transaction import atomic

from baskets.models import Basket, BasketItem, Campaign
from baskets.enums import BasketStatus
from baskets.exceptions import (PrimaryProductsQuantityMustOne,
                                BasketInvalidException,
                                BasketEmptyException,
                                BasketAmountLessThanZeroException)
from cars.exceptions import CustomerHasNoSelectedCarException
import cars.*


class BasketService:
    def get_or_create_basket(self, customer_profile):
        """
        :param customer_profile: CustomerProfile
        :return: Basket
        """
        if not customer_profile.selected_car:
            raise CustomerHasNoSelectedCarException
        try:
            basket = Basket.objects.get(customer_profile=customer_profile,
                                        status=BasketStatus.active,
                                        car=customer_profile.selected_car)
            self._check_basket_items(basket)
        except Basket.DoesNotExist:
            basket = Basket.objects.create(customer_profile=customer_profile,
                                           status=BasketStatus.active,
                                           car=customer_profile.selected_car)
        return basket

    def apply_discounts(self, basket):
        """
        :param basket: Basket
        :return: Basket
        """
        self.clean_discounts(basket)
        campaigns = Campaign.objects.actives()
        for campaign in campaigns:
            strategy = campaign.promotion_type.get_strategy(campaign=campaign,
                                                            basket=basket)
            if strategy.check():
                strategy.apply()
            else:
                message = campaign.message.format(**{'remaining': strategy.remaining})
                basket.campaign_messages.append(message)
        return basket

    def _check_basket_items(self, basket):
        """
        :param basket: Basket
        :return: bool
        """
        invalid_items = basket.basketitem_set.exclude(product__is_active=True,
                                                      product__store__is_active=True,
                                                      product__store__is_approved=True)
        if invalid_items:
            msg = _('Some services that are invalid have been removed from your basket: ')

            for bi in invalid_items:
                msg += bi.product.name + ', '
                bi.delete()

            basket.warning_messages.append(msg)

        return not invalid_items.exists()

    def add_basket_item(self, basket, product):
        """
        :param basket: Basket
        :param product: Product
        :return: BasketItem
        """
        if not basket.is_empty and basket.basketitem_set.filter(~Q(product__store=product.store)).exists():
            self.clean_basket(basket)

        try:
            basket_item = basket.basketitem_set.get(product=product)
            if basket_item.product.is_primary:
                raise PrimaryProductsQuantityMustOne
            basket_item.quantity = F('quantity') + 1
            basket_item.save()
        except BasketItem.DoesNotExist:
            primary_count = basket.basketitem_set.filter(product__is_primary=True).count()
            if primary_count > 0 and product.is_primary:
                raise PrimaryProductsQuantityMustOne
            basket_item = BasketItem.objects.create(basket=basket, product=product)

        self._check_basket_items(basket)
        return basket_item

    def clean_discounts(self, basket):
        """
        :param basket: Basket
        :return: basket
        """
        basket.campaign_messages = []
        basket.discountitem_set.all().delete()
        return basket

    def clean_basket(self, basket):
        """
        :param basket: Basket
        :return: Basket
        """
        basket.basketitem_set.all().delete()
        self.clean_discounts(basket)
        return basket

    def delete_basket_item(self, basket, product):
        """
        :param basket: Basket
        :param product: Product
        :return: None
        """
        basket.basketitem_set.filter(product=product).delete()
        self._check_basket_items(basket)

    @atomic
    def complete_basket(self, basket):
        """
        :param basket: Basket
        :return: Basket
        """
        if basket.is_empty:
            raise BasketEmptyException
        is_basket_valid = self._check_basket_items(basket)
        if not is_basket_valid:
            raise BasketInvalidException
        if basket.get_net_amount() < Decimal('0.00'):
            raise BasketAmountLessThanZeroException

        for bi in basket.basketitem_set.all():
            bi.amount = bi.get_price()
            bi.save()
        basket.status = BasketStatus.completed
        basket.save()

        return basket

    def _create_quotation_items(self, quotation: Quotation | Order, basket_items: list[BasketItem]) -> None:
        QuotationItem.objects.bulk_create(
            [
                QuotationItem(
                    quotation=quotation,
                    product_remote_id=bi.product_remote_id,
                    quantity=bi.quantity,
                    division=bi.division,
                    price=bi.price,
                    data=bi.data,
                    currency=bi.currency,
                )
                for bi in basket_items
            ]
        )

   def check_discount_code(self, code: str) -> Campaign:
        discount_code = (
            DiscountCode.objects.filter(
                Q(user=None) | Q(user=self.basket.user),
                code=code,
                campaign__is_active=True,
                campaign__start_date__lte=timezone.now(),
                campaign__end_date__gte=timezone.now(),
            )
            .select_related("user", "campaign", "campaign__condition", "campaign__benefit")
            .first()
        )

        if (
            not discount_code
            or (discount_code.campaign.start_date > timezone.now() or discount_code.campaign.end_date < timezone.now())
            or self._is_max_usage_per_user_exceeded(discount_code.campaign)
        ):
            raise DiscountCodeInvalidException

        condition = get_condition(discount_code.campaign.condition, self.basket)
        if not condition.is_satisfied():
            raise DiscountCodeFailedException(params=[condition.failed_message])

        return discount_code.campaign

    def apply_campaigns(self) -> [CampaignResult]:
        self.basket.clear_discounts()
        campaign_results: [CampaignResult] = []

        campaigns = list(self.get_campaigns())

        if self.basket.discount_code:
            try:
                discount_code_campaign = self.check_discount_code(self.basket.discount_code)
                campaigns.insert(0, discount_code_campaign)
            except (DiscountCodeInvalidException, DiscountCodeFailedException):
                self.basket.discount_code = None
                self.basket.save()

        for campaign in campaigns:
            if self._is_max_usage_per_user_exceeded(campaign):
                continue

            condition = get_condition(campaign.condition, self.basket)
            benefit = get_benefit(campaign.benefit, self.basket)

            if condition.is_satisfied():
                result = benefit.apply()
                campaign_results.append(result)
                self.basket.add_discount(result)
            else:
                failed_message = f"{condition.failed_message} {benefit.courage_message}"
                self.basket.add_failed_message(failed_message)

        return campaign_results
