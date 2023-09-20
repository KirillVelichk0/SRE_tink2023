import requests
import events, teams, users
import yaml, json
from copy import deepcopy
from collections import deque
from help_classes import CustomTimezone, OncallInfoContainer, ConvertStringTimezoneToPythonTimezone
from datetime import datetime
from datetime import timezone

class YamlConfiguratorReactions:
    reactions_queue : deque
    oncall_container: OncallInfoContainer
    team_admin : str

    def __init__(self, oncall_container: OncallInfoContainer, team_admin: str = 'root'):
        self.reactions_queue = deque()
        self.oncall_container = oncall_container
        self.team_admin = team_admin

    def React(self, key, value, context: dict):
        delayed_reaction = None
        context_type = context['context_type']
        match key:
            case 'teams':
                delayed_reaction = self.ReactToTeams(value, context)
            case 'name':
                match context_type:
                    case 'teams':
                        delayed_reaction = self.ReactToTeamName(value, context)
                    case 'teams_users':
                        delayed_reaction = self.ReactToUserName(value, context)
            case 'email':
                match context_type:
                    case 'teams':
                        delayed_reaction = self.ReactToTeamEmail(value, context)
                    case 'teams_users':
                        delayed_reaction = self.ReactToUserEmail(value, context)
            case 'scheduling_timezone':
                delayed_reaction = self.ReactToTeamTimezone(value, context)
            case 'slack_channel':
                delayed_reaction = self.ReactToTeamSlackChannel(value, context)
            case 'users':
                delayed_reaction = self.ReactToUsersInTeam(value, context)
            case 'full_name':
                delayed_reaction = self.ReactToUserFullName(value, context)
            case 'phone_number':
                delayed_reaction = self.ReactToUserPhoneNumber(value, context)
            case 'duty':
                delayed_reaction = self.ReactToTeamUserDuty(value, context)
            case 'role':
                delayed_reaction = self.ReactToDutyRole(value, context)
            case 'date':
                delayed_reaction = self.ReactToDutyDate(value, context)

        if delayed_reaction is not None:
            self.reactions_queue.append(delayed_reaction)

    def ReactToTeamName(self, value, context: dict):
        context['name'] = value

    def ReactToTeamTimezone(self, value, context: dict):
        context['scheduling_timezone'] = CustomTimezone(value)

    def ReactToTeamEmail(self, value, context: dict):
        context['email'] = value

    def ReactToTeamSlackChannel(self, value, context: dict):
        context['slack_channel'] = value


    def ReactToTeams(self, value, context: dict):
        reactions = []
        for team in value:
            current_context = dict()
            current_context['context_type'] = 'teams'
            def SpawnReact():
                local_current_context = deepcopy(current_context)
                def TeamsDelayedReaction():
                    responce= teams.TryCreateTeam(self.oncall_container, local_current_context['name'], 
                                    local_current_context['scheduling_timezone'], local_current_context['email'],
                                    local_current_context['slack_channel'], admin=self.team_admin)
                    if responce is not None:
                        return responce
                    else:
                        raise ValueError("API connection error")
                return TeamsDelayedReaction
            
            for team_key, team_value in team.items():
                self.React(team_key, team_value, current_context)
            reactions.append(SpawnReact())
        return reactions
    
    def ReactToUserName(self, value, context: dict):
        context['name'] = value

    def ReactToUserFullName(self, value, context: dict):
        context['full_name'] = value

    def ReactToUserEmail(self, value, context: dict):
        context['email'] = value

    def ReactToUserPhoneNumber(self, value, context: dict):
        context['phone_number'] = value

    def ReactToUsersInTeam(self, value, context: dict):
        reactions = []
        for user in value:
            current_context = dict()
            current_context['context_type'] = 'teams_users'
            current_context['team_name'] = context['name']
            current_context['scheduling_timezone'] = context['scheduling_timezone']
            def SpawnReact():
                local_current_context = deepcopy(current_context)
                def UsersCreateDelayedReaction():
                    contacts_data = {'call': local_current_context['phone_number'],
                        'email': local_current_context['email']}
                    responce = users.CreateFullWritedUser(self.oncall_container, 
                                                        local_current_context['name'],
                                        local_current_context['full_name'],
                                        local_current_context['scheduling_timezone'], None, contacts_data)
                    if responce is not None:
                        responce = teams.AddUserToTeam(self.oncall_container, local_current_context['team_name'],
                                                local_current_context['name'])
                        if responce is not None:
                            return responce
                        else:
                            raise ValueError('API connection error')
                    else:
                        raise ValueError('API connection error')
                return UsersCreateDelayedReaction
            for user_key, user_value in user.items():
                self.React(user_key, user_value, current_context)
            reactions.append(SpawnReact())
        return reactions
    
    def ReactToDutyDate(self, value, context: dict):
        context['date'] = value

    def ReactToDutyRole(self, value, context: dict):
        context['role'] = value
    
    def ReactToTeamUserDuty(self, value, context: dict):
        reactions = []
        for event in value:
            current_context = dict()
            current_context['context_type'] = 'duty'
            current_context['team_name'] = context['team_name']
            current_context['user'] = context['name']
            current_context['scheduling_timezone'] = context['scheduling_timezone']
            def SpawnReact():
                local_current_context = deepcopy(current_context)
                def ReactToUserDuty():
                    date = local_current_context['date']
                    tz = ConvertStringTimezoneToPythonTimezone(local_current_context['scheduling_timezone'].GetStringRepr())
                    if tz is None:
                        raise ValueError('Uncorrect Timezone')
                    def GetUnixFromDateTime(datetime_str: str):
                        datetime_object = datetime.strptime(datetime_str, '%d/%m/%Y %H:%M:%S')
                        return datetime_object.replace(tzinfo=tz).timestamp()
                    responce = events.CreateEvent(self.oncall_container, local_current_context['team_name'], 
                                    GetUnixFromDateTime(date + ' 00:00:00'), 
                                    GetUnixFromDateTime(date + ' 23:59:59'),
                                        local_current_context['user'], local_current_context['role'])
                    if responce is not None:
                        return responce
                    else:
                        raise ValueError('API connection error')
                return ReactToUserDuty
            for event_key, event_val in event.items():
                self.React(event_key, event_val, current_context)
            reactions.append(SpawnReact())
        return reactions

    def ProcessReactions(self):
        try: 
            while len(self.reactions_queue) != 0:
                reactions_list = self.reactions_queue.pop()
                for reaction in reactions_list:
                    print(reaction())
        except ValueError as e:
            print(e)
        



def SafeParceYamlFromString(yaml_data: str):
    try:
        return yaml.safe_load(yaml_data)
    except yaml.error.YAMLError as e:
        raise e

def SafeParceYamlFromFile(path: str):
    try:
        with open(path, "r") as stream:
            return yaml.safe_load(stream)
    except (yaml.error.YAMLError, OSError) as e:
        raise e
    
def ReactToAll(parsed_yaml, oncall_container: OncallInfoContainer, team_admin: str = 'root'):
    configurator = YamlConfiguratorReactions(oncall_container, team_admin)
    context = dict()
    context['context_type'] = 'NonContext'
    for key, val in parsed_yaml.items():
        configurator.React(key, val, context)
    configurator.ProcessReactions()


