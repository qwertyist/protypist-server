import requests
from dotenv import dotenv_values

config = dotenv_values(".env")
api = config["PROTYPIST_SERVER"]

print("SERVER URL is: ", api)
print("Creating Session")
resp = requests.post(api + "/session")
print("Session ID: ", resp.text)

ID = resp.text
cancel = ""
while cancel != "x": 
    cancel = input("Update info")
    resp = requests.get(api + "/session/" + ID)
    sess = resp.text
    print("Session info: ", sess)
