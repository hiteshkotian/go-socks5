#!/usr/bin/env python3
import requests

response = requests.get(
                'https://www.facebook.com', 
                proxies={"https": "socks5://127.0.0.1:1080"})

print(f'Status code {response.status_code}')
print(f'Response : {response.content}')