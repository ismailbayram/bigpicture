# Generated by Django 2.2.2 on 2019-08-01 13:28

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = [
        ('baskets', '0001_initial'),
        ('cars', '0001_initial'),
        ('products', '0001_initial'),
    ]

    operations = [
        migrations.AddField(
            model_name='basketitem',
            name='product',
            field=models.ForeignKey(on_delete=django.db.models.deletion.PROTECT, to='products.Product'),
        ),
        migrations.AddField(
            model_name='basket',
            name='car',
            field=models.ForeignKey(on_delete=django.db.models.deletion.PROTECT, to='cars.Car'),
        ),
    ]
