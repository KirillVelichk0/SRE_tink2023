import os
import yaml
import requests
import subprocess
import time
import json
dirname = os.path.dirname(__file__)
config_path = os.path.join(dirname, 'config.yaml')
green_path = os.path.join(dirname, os.pardir, os.pardir, 'docker-compose.green.yml')
blue_path = os.path.join(dirname, os.pardir, os.pardir, 'docker-compose.blue.yml')

print("Paths", config_path, blue_path, green_path)

with open("config.yaml", "r") as file:
    config = yaml.safe_load(file)

print('config', config)

def CheckGreenToStart() -> bool:
    result = subprocess.run(['sudo', 'bash', os.path.join(dirname, 'get_green.sh')], stdout=subprocess.PIPE)
    pr_text = result.stdout
    print("text", pr_text)
    return len(pr_text) != 0

def CheckBlueToStart() -> bool:
    result = subprocess.run(['sudo', 'bash', os.path.join(dirname, 'get_blue.sh')], stdout=subprocess.PIPE)
    pr_text = result.stdout
    print("text", pr_text)
    return len(pr_text) != 0

def CheckToOk(url) -> bool:
    print("url", url + '/')
    for i in range(0, config['seconds_to_try']):
        try:
            r = requests.get(url + '/')
            if r.status_code == 200:
                return True
        except:
            time.sleep(1)
    return False

def TryToStart(path):
    os.system("sudo docker stack deploy --with-registry-auth -c" + path + " prod")

def TryToRm(name):
    os.system("sudo docker service rm "+ name)

def TryToRmGreen():
    TryToRm(config["green_name"])

def TryToRmBlue():
    TryToRm(config["blue_name"])

def TryToStartBlue():
    TryToStart(blue_path)

def TryToStartGreen():
    TryToStart(green_path)

def CheckGreenToOk():
    return CheckToOk(config['green_url'])

def CheckBlueToOk():
    return CheckToOk(config['blue_url'])


def ChangeColor(isGreen):
    isGreenq = {"IsGreen":isGreen}
    js = json.dumps(isGreenq)
    r = requests.get(config["redirecter_url"] + "/set_state", json=isGreenq)
    print(r.status_code)

def ChangeToGreen():
    ChangeColor(True)

def ChangeToBlue():
    ChangeColor(False)


isGreen = CheckGreenToStart()
isBlue = CheckBlueToStart()

if isGreen and isBlue:
    print("Blue and green in one moment")
elif not isGreen and not isBlue:
    print("All colors are bad")
elif isGreen:
    print("Try change to blue")
    TryToStartBlue()
    if CheckBlueToOk():
        print("Blue is ok")
        ChangeToBlue()
        TryToRmGreen()
    else:
        print("Blue is bad")
        TryToRmBlue()
else:
    print("Try change to green")
    TryToStartGreen()
    if CheckGreenToOk():
        print("Green is ok")
        ChangeToGreen()
        TryToRmBlue()
    else:
        print("Green is bad")
        TryToRmGreen()
