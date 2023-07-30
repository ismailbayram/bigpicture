from django.test import TestCase

from users.service import UserService
from users.enums import GroupType
from cars.models import Car
from cars.service import CarService
from cars.enums import CarType


class CarServiceTest(TestCase):
    def setUp(self):
        user_service = UserService()

        data = {
            "first_name": "Customer 1",
            "last_name": "CusLast",
            "phone_number": "555111",
            "group_type": GroupType.customer
        }
        self.customer, self.customer_token = user_service.get_or_create_user(**data)
        self.service = CarService()

        self.car1 = self.service.create_car(
            licence_plate="09 TK 40",
            car_type=CarType.sedan,
            customer_profile=self.customer.customer_profile
        )
        self.car2 = self.service.create_car(
            licence_plate="09 TK 41",
            car_type=CarType.sedan,
            customer_profile=self.customer.customer_profile
        )
        self.car3 = self.service.create_car(
            licence_plate="09 TK 42",
            car_type=CarType.suv,
            customer_profile=self.customer.customer_profile
        )

    def test_create_car(self):
        data = {
            'licence_plate': "07 Z 8686",
            'car_type': CarType.sedan,
            'customer_profile': self.customer.customer_profile
        }
        self.customer.customer_profile.cars.all().delete()
        self.car1 = self.service.create_car(**data)
        self.assertIsInstance(self.car1, Car)
        self.assertEqual(self.car1.licence_plate, "07 Z 8686")
        self.assertTrue(self.car1.is_selected)

    def test_update_car(self):
        data = {
            "car": self.car1,
            "licence_plate": "06 AT 88",
            "car_type": CarType.suv,
        }
        car = self.service.update_car(**data)
        self.assertEqual(car.licence_plate, "06 AT 88")
        self.assertEqual(car.car_type, CarType.suv)

        data = {
            "car": self.car1,
            "licence_plate": "22 a 22",
        }
        car = self.service.update_car(**data)
        self.assertEqual(car.licence_plate, "22 a 22")

        data = {
            "car": self.car1,
            "car_type": CarType.suv,
        }
        car = self.service.update_car(**data)
        self.assertEqual(car.car_type, CarType.suv)

    def test_select_car(self):
        self.assertTrue(self.car1.is_selected)
        self.assertFalse(self.car2.is_selected)
        self.assertFalse(self.car3.is_selected)

        self.service.select_car(self.car1, self.customer.customer_profile)
        self.car1.refresh_from_db()
        self.car2.refresh_from_db()
        self.car3.refresh_from_db()
        self.assertTrue(self.car1.is_selected)
        self.assertFalse(self.car2.is_selected)
        self.assertFalse(self.car3.is_selected)

        self.service.select_car(self.car2, self.customer.customer_profile)
        self.car1.refresh_from_db()
        self.car2.refresh_from_db()
        self.car3.refresh_from_db()
        self.assertFalse(self.car1.is_selected)
        self.assertTrue(self.car2.is_selected)
        self.assertFalse(self.car3.is_selected)

        self.service.select_car(self.car3, self.customer.customer_profile)
        self.car1.refresh_from_db()
        self.car2.refresh_from_db()
        self.car3.refresh_from_db()
        self.assertFalse(self.car1.is_selected)
        self.assertFalse(self.car2.is_selected)
        self.assertTrue(self.car3.is_selected)

    def test_deactivate_car(self):
        self.customer.customer_profile.cars.all().delete()
        car1 = self.service.create_car(
            licence_plate="09 TK 40",
            car_type=CarType.sedan,
            customer_profile=self.customer.customer_profile
        )
        car2 = self.service.create_car(
            licence_plate="09 TK 41",
            car_type=CarType.sedan,
            customer_profile=self.customer.customer_profile
        )
        self.service.deactivate_car(car1)
        car2.refresh_from_db()
        self.assertTrue(car2.is_selected)
