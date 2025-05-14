import requests
import pprint
import os

class Telegram:
    __HOST = 'https://api.telegram.org/bot'
    
    @classmethod
    def send(cls, msg: str):

        token = os.getenv('TELEGRAM_TOKEN')
        chat_id = os.getenv('TELEGRAM_CHAT_ID')
        if token is None or chat_id is None : 
            print('TELEGRAM_TOKEN or TELEGRAM_CHAT_ID not found')
            exit(1)

        try:
            response = requests.post(cls.__HOST + token + "/sendMessage", data=[
                ('chat_id', chat_id),
                ('text', os.path.dirname(__file__) + "\n" + msg),
            ])
            if (response.status_code != 200):
                print(response.status_code)
                pprint.pprint(response.json())
        except Exception as e:
            print(str(e))