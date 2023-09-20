import requests, json
from help_classes import ConstructJsonFromDict, OncallInfoContainer, ConstructParamsFromDict, TimezoneRepr

def CreateUser(oncall_container: OncallInfoContainer, name: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/users'
    json_data = json.dumps({'name':name})
    try:
        responce = oncall_container.session.post(url= uri, data=json_data)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

def ModifyUser(oncall_container: OncallInfoContainer, name: str, full_name: str= None,
               tz = None, photo_url: str = None, contacts: dict = None, active: bool = None,
               new_name: str = None):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/users/' + name
    if tz is not None:
        tz = tz.GetStringRepr()
    json_data = ConstructJsonFromDict({'contacts': contacts, 'name': new_name,
                                       'full_name':full_name,
                                       'time_zone': tz,
                                       'photo_url': photo_url,
                                       'active': active})
    try:
        responce = oncall_container.session.put(uri, data=json_data)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

def CreateFullWritedUser(oncall_container: OncallInfoContainer, name: str, full_name: str= None,
               tz = None, photo_url: str = None, contacts: str = None):
    try:
        responce = CreateUser(oncall_container, name)
        if responce is not None:
            responce = ModifyUser(oncall_container, name=name, full_name=full_name, tz=tz,
                                  photo_url=photo_url, contacts=contacts, active=True, new_name=name)
    except requests.RequestException as e:
        print(str(e) + ' while creating full named user')
        responce = None
    return responce