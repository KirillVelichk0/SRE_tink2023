from help_classes import OncallInfoContainer, US_PacificTimezone_repr, CustomTimezone
from yaml_configurator import SafeParceYamlFromFile, ReactToAll
import users
import json

cont = OncallInfoContainer('localhost', 8080, 'http://')
cont.key = 'test_key'
cont.app = 'test_app'
cont.Login('root', '1234')


yaml_data = SafeParceYamlFromFile('./data.yaml')
ReactToAll(yaml_data, cont)
