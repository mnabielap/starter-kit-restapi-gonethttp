import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils

token = utils.load_config("access_token")
user_id = utils.load_config("user_id")

if not user_id:
    print("No user_id found in secrets.json. Please run auth_register.py first.")
    sys.exit(1)

URL = f"{utils.BASE_URL}/users/{user_id}"

response = utils.send_and_print(
    url=URL,
    method="GET",
    headers={
        "Authorization": f"Bearer {token}"
    },
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)