import json
from model_mommy import mommy
from django.test import TestCase
from rest_framework.reverse import reverse_lazy
from rest_framework import status

from base.test import BaseTestViewMixin
from cars.enums import CarType
from cars.service import CarService
from products.service import ProductService
from baskets.service import BasketService


class BasketViewSetTest(TestCase, BaseTestViewMixin):
    def setUp(self):
        self.service = BasketService()
        self.product_service = ProductService()
        self.car_service = CarService()
        self.init_users()
        self.store = mommy.make('stores.Store', washer_profile=self.washer_profile,
                                is_approved=True, is_active=True)
        self.store2 = mommy.make('stores.Store', washer_profile=self.washer2_profile,
                                 is_approved=True, is_active=True)
        self.product1 = self.product_service.create_primary_product(self.store)
        self.product2 = self.product_service.create_product(name='Parfume', store=self.store,
                                                            washer_profile=self.store.washer_profile)
        self.product3 = self.product_service.create_primary_product(self.store2)
        self.car = self.car_service.create_car(licence_plate="34FH3773", car_type=CarType.sedan,
                                               customer_profile=self.customer_profile)

    def test_view_basket(self):
        url = reverse_lazy('api:basket')

        response = self.client.get(url)
        self.assertEqual(response.status_code, status.HTTP_401_UNAUTHORIZED)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.customer2_token}'}
        response = self.client.get(url, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_406_NOT_ACCEPTABLE)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.customer_token}'}
        response = self.client.get(url, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        jresponse = json.loads(response.content)
        self.assertEqual(jresponse['basket']['total_amount'], '0.00')
        self.assertEqual(jresponse['basket']['total_quantity'], 0)
        self.assertEqual(jresponse['basket']['car']['pk'], self.car.pk)
        self.assertEqual(len(jresponse['basket']['basketitem_set']), 0)
        self.assertEqual(len(jresponse['basket']['discountitem_set']), 0)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.washer_token}'}
        response = self.client.get(url, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_403_FORBIDDEN)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.worker_token}'}
        response = self.client.get(url, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_403_FORBIDDEN)

    def test_add_item(self):
        url = reverse_lazy('api:basket')
        data = {
            "product": self.product1.pk
        }

        response = self.client.post(url, data=data)
        self.assertEqual(response.status_code, status.HTTP_401_UNAUTHORIZED)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.customer_token}'}
        response = self.client.post(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_201_CREATED)
        jresponse = json.loads(response.content)
        self.assertEqual(jresponse['basket']['total_amount'], str(self.product1.price(self.car.car_type)))
        self.assertEqual(jresponse['basket']['total_quantity'], 1)
        self.assertEqual(len(jresponse['basket']['basketitem_set']), 1)

        response = self.client.post(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_406_NOT_ACCEPTABLE)

        data.update({"product": self.product2.pk})
        response = self.client.post(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_201_CREATED)
        jresponse = json.loads(response.content)
        self.assertEqual(jresponse['basket']['total_amount'],
                         str(self.product1.price(self.car.car_type) + self.product2.price(self.car.car_type)))
        self.assertEqual(jresponse['basket']['total_quantity'], 2)
        self.assertEqual(len(jresponse['basket']['basketitem_set']), 2)

        self.product_service.delete_product(self.product2)
        response = self.client.get(url, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        jresponse = json.loads(response.content)
        self.assertNotEqual(len(jresponse['basket']['warning_messages']), 0)

        data.update({"product": self.product3.pk})
        response = self.client.post(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_201_CREATED)
        jresponse = json.loads(response.content)
        self.assertEqual(jresponse['basket']['total_amount'], str(self.product3.price(self.car.car_type)))
        self.assertEqual(jresponse['basket']['total_quantity'], 1)
        self.assertEqual(len(jresponse['basket']['basketitem_set']), 1)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.washer_token}'}
        response = self.client.post(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_403_FORBIDDEN)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.worker_token}'}
        response = self.client.post(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_403_FORBIDDEN)

    def test_delete_item(self):
        url = reverse_lazy('api:basket')
        data = {
            "product": self.product1.pk
        }

        response = self.client.delete(url, data=data)
        self.assertEqual(response.status_code, status.HTTP_401_UNAUTHORIZED)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.customer_token}'}
        response = self.client.post(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_201_CREATED)

        response = self.client.delete(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_204_NO_CONTENT)
        response = self.client.get(url, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        jresponse = json.loads(response.content)
        self.assertEqual(jresponse['basket']['total_amount'], '0.00')
        self.assertEqual(jresponse['basket']['total_quantity'], 0)
        self.assertEqual(len(jresponse['basket']['basketitem_set']), 0)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.washer_token}'}
        response = self.client.delete(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_403_FORBIDDEN)

        headers = {'HTTP_AUTHORIZATION': f'Token {self.worker_token}'}
        response = self.client.delete(url, data=data, content_type='application/json', **headers)
        self.assertEqual(response.status_code, status.HTTP_403_FORBIDDEN)
