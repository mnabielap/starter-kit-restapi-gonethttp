import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils
import random

token = utils.load_config("access_token")
random_id = random.randint(1000, 9999)

payload = {
    "name": "Created Via API",
    "email": f"created{random_id}@example.com",
    "password": "password123",
    "role": "user"
}

URL = f"{utils.BASE_URL}/users"

response = utils.send_and_print(
    url=URL,
    method="POST",
    body=payload,
    headers={
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}"
    },
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)