from locust import HttpUser, task, between, constant

class RateLimitTest(HttpUser):
    # wait_time = constant(1)
    @task
    def test_health_check_lock_token(self):
        headers = {'content-type': 'application/json'}
        headers = {'API_KEY': '4a97fe3c-ea92-489f-8d9a-4d93f2ea9ce1}
        self.client.get("/health", headers=headers, name="Health check token lock") 

    @task
    def test_health_check_lock_ip(self):
        headers = {'content-type': 'application/json'}
        headers = {'X-Forwarded-For': '192.168.9.9'}
        self.client.get("/health", headers=headers, name="Health check ip lock")