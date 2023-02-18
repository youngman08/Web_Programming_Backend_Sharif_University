from locust import HttpUser, between, task
import uuid
import random

hash_to_string = {}
hashes = []

class WebsiteUser(HttpUser):
    host = "http://localhost/"
    
    # main page
    @task(1)
    def get_mainpage(self):
        result = self.client.get(f"/").json()
        # assert result['status'] == True

    # sign up
    @task(19)
    def post_signup(self):
        for _ in range(100):
            result = self.client.post("api/auth/signup", json={"Name":f"parsa{_}", "Email":f"parsa{_}@gmail.com", "PassportNumber":"123" , "Password":"123"}).json()
        # assert result['status'] == True
        
    # login
    @task(10)
    def post_login(self):
        for _ in range(100):
            result = self.client.post("api/auth/login", json={"Email":f"parsa{_}@gmail.com", "Password":"123"}).json()
        # assert result['status'] == True

    # search
    @task(10)
    def get_search(self):
        for _ in range(100):
            result = self.client.get(f"api/ticket/search", json={"From": "Tehran", "To":"Mashhad", "date": "2022-10-10", "passengerCount":"2"}).json()
        # assert result['status'] == True

    # ticket transaction 
    @task
    def get_ticket_transaction(self):
        for _ in range(1, 100):
            result = self.client.get(f"api/ticket/transaction", json={"receipt_id": f"{_}", "amount":f"{_}"}).json()
        # assert result['status'] == True

    # dashboard
    @task(10)
    def get_dashboard(self):
        result = self.client.get("dashboard.html")
