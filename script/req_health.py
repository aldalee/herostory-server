import time

import requests

URL = "http://localhost:12345/health"


def health_check():
    count = 0
    while True:
        try:
            resp = requests.get(URL, timeout=5)
            print(f"[{count}] status: {resp.status_code}, body: {resp.text}")
        except Exception as e:
            print(f"[{count}] request failed: {e}")
        count += 1
        time.sleep(2)


if __name__ == "__main__":
    health_check()
