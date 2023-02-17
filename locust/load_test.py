from locust import HttpUser, between, task
import uuid
import random

hash_to_string = {}
hashes = []

class WebsiteUser(HttpUser):
    host = "http://localhost"

    # def get_go_node(self) -> str:
    #     if random.random() < 0.5:
    #         return 'go'
    #     return 'node'
    
    # main page
    @task(1)
    def get_mainpage(self):
        string = str(uuid.uuid1())[:random.randint(0, 7)]
        result = self.client.post(f"/8000").json()
        assert result['status'] == False

    # sign up
    @task(19)
    def post_signup(self):
        result = self.client.post(f"/8000/signup").json()
        # assert result['status'] == True
        
    # login
    @task(10)
    def post_login(self):
        result = self.client.post(f"/8000/login").json()
        # assert result['found'] == True
        # string = result['string']

    # search
    @task(10)
    def get_dashboard(self):
        result = self.client.get(f"/8001/search").json()

    # dashboard
    @task(10)
    def get_dashboard(self):
        result = self.client.get(f"/8001/dashboard").json()