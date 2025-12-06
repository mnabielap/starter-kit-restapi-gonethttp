import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

import utils

token = utils.load_config("access_token")

URL = f"{utils.BASE_URL}/users?page=1&limit=5&sortBy=created_at%20desc"

response = utils.send_and_print(
    url=URL,
    method="GET",
    headers={
        "Authorization": f"Bearer {token}"
    },
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)