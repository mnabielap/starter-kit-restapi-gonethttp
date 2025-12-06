import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils

token = utils.load_config("access_token")

URL = f"{utils.BASE_URL}/auth/send-verification-email"

response = utils.send_and_print(
    url=URL,
    method="POST",
    headers={
        "Authorization": f"Bearer {token}"
    },
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)