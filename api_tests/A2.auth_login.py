import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils

email = "admin@example.com"
password = "password123"

payload = {
    "email": email,
    "password": password
}

URL = f"{utils.BASE_URL}/auth/login"

print(f"--- LOGGING IN AS: {email} ---")

response = utils.send_and_print(
    url=URL,
    method="POST",
    body=payload,
    headers={"Content-Type": "application/json"},
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)

if response.status_code == 200:
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
            print("\n[SUCCESS] Login successful. Credentials SAVED to secrets.json")
        else:
            print("\n[ERROR] JSON structure recognized, but keys (id/token) missing.")
            print(f"Debug Data: {data}")
    else:
        print("\n[ERROR] Empty JSON response.")
else:
    print(f"\n[ERROR] Login failed with status {response.status_code}")