FROM python:3

WORKDIR /usr/src/app


RUN pip install Django djangorestframework

COPY . .

RUN python manage.py migrate

CMD ["python", "manage.py", "runserver", "0.0.0.0:8080"]
