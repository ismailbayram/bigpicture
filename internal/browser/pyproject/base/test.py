from model_mommy import mommy

from users.enums import GroupType
from users.service import UserService


class BaseTestViewMixin:
    def init_users(self):
        service = UserService()
        self.superuser = mommy.make('users.User', is_staff=True)
        self.superuser_token = service._create_token(self.superuser)

        data = {
            "first_name": "Customer 1",
            "last_name": "CusLast",
            "phone_number": "555111",
            "group_type": GroupType.customer
        }
        self.customer, self.customer_token = service.get_or_create_user(**data)
        self.customer_profile = self.customer.customer_profile

        data = {
            "first_name": "Customer 2",
            "last_name": "CusLast 2",
            "phone_number": "5551112",
            "group_type": GroupType.customer
        }
        self.customer2, self.customer2_token = service.get_or_create_user(**data)
        self.customer2_profile = self.customer2.customer_profile

        data = {
            "first_name": "Worker 1",
            "last_name": "WorkLast",
            "phone_number": "555222",
            "group_type": GroupType.worker
        }
        self.worker, self.worker_token = service.get_or_create_user(**data)
        self.worker_profile = self.worker.worker_profile

        data = {
            "first_name": "Worker 2",
            "last_name": "WorkLast",
            "phone_number": "555223",
            "group_type": GroupType.worker
        }
        self.worker2, self.worker2_token = service.get_or_create_user(**data)
        self.worker2_profile = self.worker2.worker_profile

        data = {
            "first_name": "Washer 1",
            "last_name": "WashLast",
            "phone_number": "555333",
            "group_type": GroupType.washer
        }
        self.washer, self.washer_token = service.get_or_create_user(**data)
        self.washer_profile = self.washer.washer_profile

        data = {
            "first_name": "Washer 2",
            "last_name": "WashLast",
            "phone_number": "555334",
            "group_type": GroupType.washer
        }
        self.washer2, self.washer2_token = service.get_or_create_user(**data)
        self.washer2_profile = self.washer2.washer_profile
