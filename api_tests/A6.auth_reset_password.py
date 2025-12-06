import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils

# REPLACE THIS with a valid token from your database/logs
dummy_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy"

payload = {
    "password": "newpassword123"
}

URL = f"{utils.BASE_URL}/auth/reset-password?token={dummy_token}"

response = utils.send_and_print(
    url=URL,
    method="POST",
    body=payload,
    headers={"Content-Type": "application/json"},
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)