import os
import requests

from flask import Flask, request, jsonify

app = Flask(__name__)

TELEGRAM_TOKEN = os.environ["TELEGRAM_TOKEN"]
CHAT_WARN_ID = os.environ["TELEGRAM_WARN_CHAT_ID"]
CHAT_CRIT_ID = os.environ["TELEGRAM_CRITICAL_CHAT_ID"]
TELEGRAM_URL = f"https://api.telegram.org/bot{TELEGRAM_TOKEN}/sendMessage"


@app.route("/alert/warn", methods=["POST"])
def webhook_warn():
    data = request.json
    for alert in data["alerts"]:
        alertname = alert["labels"].get("alertname", "Uknown")
        instance = alert["labels"].get("instance", "Uknown")
        summary = alert["annotations"]["summary"]
        CHAT_ID = CHAT_WARN_ID
        text = f"⚠️ Alarm! {alertname} for {instance}. \n{summary}"
        payload = {"chat_id": CHAT_ID, "text": text}
        response = requests.post(TELEGRAM_URL, data=payload)
    return jsonify(status="success"), 200

@app.route("/alert/crit", methods=["POST"])
def webhook_crit():
    data = request.json
    for alert in data["alerts"]:
        alertname = alert["labels"].get("alertname", "Uknown")
        instance = alert["labels"].get("instance", "Uknown")
        summary = alert["annotations"]["summary"]
        CHAT_ID = CHAT_CRIT_ID
        text = f"⚠️ Alarm! {alertname} for {instance}. \n{summary}"
        payload = {"chat_id": CHAT_ID, "text": text}
        response = requests.post(TELEGRAM_URL, data=payload)
    return jsonify(status="success"), 200


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000)