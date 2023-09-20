import requests
import json
from help_classes import US_PacificTimezone_repr, OncallInfoContainer, ConstructParamsFromDict, TimezoneRepr


def TryGetTeams(oncall_container: OncallInfoContainer, 
             team_name: str = None, contained_param: str = None, name_start: str = None,
               name_end: str = None, id: int = None, active: bool = None):
    params = {'name': team_name, 'name__contains': contained_param, 'name__startswith':name_start,
              'name__endswith': name_end, 'id':id, 'active':active}
    params_str = ConstructParamsFromDict(params)
    if len(params_str) == 0:
        raise ValueError('Not correct args')
    host_port_str = oncall_container.GetURL()
    try:
        responce = oncall_container.session.get(host_port_str + '/api/v0/teams' + params_str)
    except (Exception):
        responce = None
    return responce

def TryGetTeamsFromName(oncall_container: OncallInfoContainer, team_name: str, active: bool = None):
    return TryGetTeams(oncall_container, team_name=team_name, active=active)

def TryGetTeamsFromId(oncall_container: OncallInfoContainer, team_id: int, active: bool = None):
    return TryGetTeams(oncall_container=oncall_container, id=team_id, active=active)



def TryCreateTeam(oncall_container: OncallInfoContainer, team_name: str, timezone, 
               email: str, slack_channel: str, admin: str = None):
    if not issubclass(type(timezone), TimezoneRepr):
        return None
    timezone_str : str = timezone.GetStringRepr()
    request_dict = {"name": team_name, "scheduling_timezone": timezone_str,
    "email": email, "slack_channel": slack_channel, 'csrf_token':oncall_container.token}
    if admin is not None:
        request_dict['admin'] = admin
    request_json = json.dumps(request_dict)
    host_port_str = oncall_container.GetURL()
    url = host_port_str + '/api/v0/teams'
    
    try:
        responce = oncall_container.session.request(url=url, data=request_json, method='POST')
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce
    
def EditTeamInfo(oncall_container: OncallInfoContainer, old_name: str,
                  timezone, email: str, slack_channel: str, new_name: str = None):
    if not issubclass(type(timezone), TimezoneRepr):
        return None
    if new_name is None:
        new_name = old_name
    json_data = json.dumps({
        "name": new_name,
        "slack_channel": slack_channel,
        "email": email,
        "scheduling_timezone": timezone.GetStringRepr()
    })
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + old_name
    try:
        oncall_container.session.put(uri, data=json_data)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

def DeleteTeam(oncall_container: OncallInfoContainer, team_name:str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name
    try:
        responce = oncall_container.session.delete(uri)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

def GetTeamSummary(oncall_container: OncallInfoContainer, team_name: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name + '/summary'
    try:
        responce = oncall_container.session.get(uri)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce


def GetOncallEventFromRole(oncall_container: OncallInfoContainer, role: str, team_name: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name + '/oncall/' + role
    try:
        responce = oncall_container.session.get(uri)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

def GetTeamAdmins(oncall_container: OncallInfoContainer, team_name: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name + '/admins'
    try:
        responce = oncall_container.session.get(uri)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

#NEED TO CHECK. DANGEROUS
def AddAdminToTeam(oncall_container: OncallInfoContainer, team_name: str, user: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name + '/admins'
    json_data = json.dumps({'name':user})
    try:
        responce = oncall_container.session.post(uri, data=json_data)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

def DeleteAdminFromTeam(oncall_container: OncallInfoContainer, team_name: str, user: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name + '/admins/' + user
    try:
        responce = oncall_container.session.delete(uri)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

def GetUsersInTeam(oncall_container: OncallInfoContainer, team_name: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name + '/users'
    try:
        responce = oncall_container.session.get(uri)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce


def AddUserToTeam(oncall_container: OncallInfoContainer, team_name: str, user: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name + '/users'
    json_data = json.dumps({'name': user})
    try:
        responce = oncall_container.session.post(url=uri, data=json_data)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce

def DeleteUserFromTeam(oncall_container: OncallInfoContainer, team_name: str, user: str):
    host_port_str = oncall_container.GetURL()
    uri = host_port_str + '/api/v0/teams/' + team_name + '/users/' + user
    try:
        responce = oncall_container.session.delete(uri)
    except requests.RequestException as e:
        print(str(e))
        responce = None
    return responce