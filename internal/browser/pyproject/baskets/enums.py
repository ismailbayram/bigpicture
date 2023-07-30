from enumfields import Enum


class BasketStatus(Enum):
    active = 'active'
    completed = 'completed'


class PromotionType(Enum):
    one_free_in_nine = 'one_free_in_nine'

    @property
    def get_strategy(self):
        from baskets.strategies import OneFreeInNineStrategy

        if self.value == 'one_free_in_nine':
            return OneFreeInNineStrategy
