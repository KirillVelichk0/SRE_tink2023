import requests, json
from help_classes import ConstructJsonFromDict, OncallInfoContainer, ConstructParamsFromDict, TimezoneRepr

def CreateEvent(oncall_container: OncallInfoContainer, team: str, start_time_unix: int,
                end_time_unix, user: str, role: str):
    host_port_url = oncall_container.GetURL()
    uri = host_port_url + '/api/v0/events'
    json_data = json.dumps({"start": start_time_unix,
    "end": end_time_unix,
    "user": user,
    "team": team,
    "role": role})
    try:
        responce = oncall_container.session.post(url = uri, data=json_data)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce