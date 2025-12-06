import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils
import random

random_id = random.randint(1000, 9999)
payload = {
    "name": "New User",
    "email": f"testuser{random_id}@example.com",
    "password": "password123",
    "role": "user"
}

URL = f"{utils.BASE_URL}/auth/register"

print(f"--- REGISTERING USER: {payload['email']} ---")

response = utils.send_and_print(
    url=URL,
    method="POST",
    body=payload,
    headers={"Content-Type": "application/json"},
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)

if response.status_code == 201:
    data = response.json()
    if data:
        target_data = data.get("results", data)

        user = target_data.get("user", {})
        tokens = target_data.get("tokens", {})
        
        user_id = user.get("id")
        access_token = tokens.get("access", {}).get("token")
        refresh_token = tokens.get("refresh", {}).get("token")

        if user_id and access_token:
            utils.save_config("user_id", user_id)
            utils.save_config("user_email", user.get("email"))
            utils.save_config("access_token", access_token)
            utils.save_config("refresh_token", refresh_token)
            print("\n[SUCCESS] Credentials SAVED to secrets.json")
        else:
            print("\n[ERROR] Structure matches but ID or Token missing.")
    else:
        print("\n[ERROR] No data in response.")
else:
    print(f"\n[ERROR] Registration failed with status {response.status_code}")