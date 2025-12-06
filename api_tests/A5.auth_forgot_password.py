import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils

email = utils.load_config("user_email") or "testuser@example.com"

payload = {
    "email": email
}

URL = f"{utils.BASE_URL}/auth/forgot-password"

response = utils.send_and_print(
    url=URL,
    method="POST",
    body=payload,
    headers={"Content-Type": "application/json"},
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)