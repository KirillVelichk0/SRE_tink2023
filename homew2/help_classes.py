from abc import ABC, abstractmethod
import requests
import json
from datetime import timezone, timedelta

class TimezoneRepr(ABC):
    @abstractmethod
    def GetStringRepr(self):
        pass

class US_PacificTimezone_repr(TimezoneRepr):
    def __init__(self):
        ...
    def GetStringRepr(self):
        return 'US/Pacific'
    
class CustomTimezone(TimezoneRepr):
    def __init__(self, text_repr: str):
        self.text_repr = text_repr
    
    def GetStringRepr(self):
        return self.text_repr


class OncallInfoContainer:
    host : str
    port : int
    adapter : str
    app: str
    key: str

    def __init__(self, host: str, port: int, adapter: str):
        self.host = host
        self.port = port
        self.adapter = adapter
        self.session = requests.session()

    def GetURL(self):
        return self.adapter + self.host + ':' + str(self.port)

    def Login(self, login: str, password: str):
        data = f'username={login}&password={password}'
        res = self.session.post(self.GetURL() + '/login', data=data)
        if res.ok:
            self.token = res.json()['csrf_token']
        else:
            self.token = None

def ConstructParamsFromDict(dict_params):
    result = '?'
    for key, val in dict_params.items():
        if val is not None:
            result = result + str(key) + '=' + str(val) +'&'
    '''Если параметры есть, то удаляется лишний последний &. Если нет, то удалится ?'''
    result = result[:-1]
    return result

def ConstructJsonFromDict(dict_params):
    return json.dumps({k: v for k, v in dict_params.items() if v is not None})

def ConvertStringTimezoneToPythonTimezone(timezone_str: str):
    offset = None
    print(timezone_str == 'Asia/Novosibirsk')
    print(timezone_str)
    if timezone_str == "Europe/Moscow":
        offset = timedelta(hours=3)
    if timezone_str == 'Asia/Novosibirsk':
        offset = timedelta(hours=7)
        'Asia/Novosibirsk'
    if offset is not None:
        return timezone(offset, timezone_str)
    