import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils

# REPLACE THIS with a valid token from your database/logs if you want a 204 Success
dummy_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.dummy"

URL = f"{utils.BASE_URL}/auth/verify-email?token={dummy_token}"

response = utils.send_and_print(
    url=URL,
    method="POST",
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)