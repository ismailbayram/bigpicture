from api.exceptions import ProjectBaseException
from cars import codes


class CustomerHasNoSelectedCarException(ProjectBaseException):
    code = codes.car_100_0
