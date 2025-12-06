import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils

refresh_token = utils.load_config("refresh_token")

if not refresh_token:
    print("No refresh_token found in secrets.json")
    sys.exit(1)

payload = {
    "refreshToken": refresh_token
}

URL = f"{utils.BASE_URL}/auth/logout"

response = utils.send_and_print(
    url=URL,
    method="POST",
    body=payload,
    headers={"Content-Type": "application/json"},
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)