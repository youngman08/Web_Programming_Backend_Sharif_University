FROM python:3

WORKDIR /usr/src/app

COPY . .

RUN pip install Django djangorestframework

RUN python manage.py migrate

CMD ["python", "manage.py", "runserver", "0.0.0.0:8080"]
