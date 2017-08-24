import boto3
import json

""" sample event
{
  "session": {
    "sessionId": "SessionId.6ab325dd-xxxx-xxxx-aee5-456cd330932a",
    "application": {
      "applicationId": "amzn1.echo-sdk-ams.app.bd304b90-xxxx-xxxx-86ae-1e4fd4772bab"
    },
    "attributes": {},
    "user": {
      "userId": "amzn1.ask.account.XXXXXX"
    },
    "new": true
  },
  "request": {
    "type": "IntentRequest",
    "requestId": "EdwRequestId.b851ed18-2ca8-xxxx-xxxx-cca3f2b521e4",
    "timestamp": "2016-07-05T15:27:34Z",
    "intent": {
      "name": "GetTrainTimes",
      "slots": {
        "Station": {
          "name": "Station",
          "value": "Balboa Park"
        }
      }
    },
    "locale": "en-US"
  },
  "version": "1.0"
}

"""


APPLICATION_ID = "amzn1.ask.skill.28f38919-c3b3-4d7b-85b3-0c48156f342a"


def lambda_handler(event, context):

    # ensure this was called from the correct application latform
    if event['session']['application']['applicationId'] != APPLICATION_ID:
        raise ValueError("Invalid Application ID")

    if event["session"]["new"]:
        on_session_started({"requestId": event["request"]["requestId"]}, event["session"])

    if event["request"]["type"] == "LaunchRequest":
        return on_launch(event["request"], event["session"])
    elif event["request"]["type"] == "IntentRequest":
        return on_intent(event["request"], event["session"])
    elif event["request"]["type"] == "SessionEndedRequest":
        return on_session_ended(event["request"], event["session"])


def on_session_started(session_started_request, session):
    print "Starting new session."


def on_launch(launch_request, session):
    return get_welcome_response()


def on_intent(intent_request, session):
    intent = intent_request["intent"]
    intent_name = intent_request["intent"]["name"]

    if intent_name == "LaunchApp":
        return launch_application(intent)
    elif intent_name == "PowerTV":
        return power_tv()
    elif intent_name == "AMAZON.HelpIntent":
        return get_welcome_response()
    else:
        raise ValueError("Invalid intent")


def on_session_ended(session_ended_request, session):
    print "Ending session."
    # Cleanup goes here...


def get_welcome_response():
    session_attributes = {}
    card_title = "ROKU"
    speech_output = "Welcome to the Alexa Roku skill. " \
                    "You can ask me to launch an application, or " \
                    "ask me to power on or off the television." \
                    "Supported applications include Netflix, Amazon, Pandora, and DS Video"
    reprompt_text = "Please ask me to launch an application or power the television."
    should_end_session = False
    return build_response(session_attributes, build_speechlet_response(
        card_title, speech_output, reprompt_text, should_end_session))


def launch_application(intent):
    session_attributes = {}
    card_title = "Roku Application Launch"
    reprompt_text = ""
    should_end_session = False

    # get application name
    if "App" in intent["slots"]:
        app = intent["slots"]["App"]["value"]
        app_name = map_application_name(app.lower())
    else:
        app_name, app = "unknown", "unknown"

    if app_name != "unknown":

        # get queue
        queue = get_queue('roku-control')

        # send message
        queue.send_message(MessageBody=json.dumps({"command": "launch_app",
                                                   "data": {"app": app_name}}))

        speech_output = "Sent request to launch application {}".format(app_name)
    else:
        speech_output = "Unknown application {}".format(app)

    return build_response(session_attributes, build_speechlet_response(
        card_title, speech_output, reprompt_text, should_end_session))


def power_tv():
    session_attributes = {}
    card_title = "Roku System Power"
    reprompt_text = ""
    should_end_session = False

    # get queue
    queue = get_queue('roku-control')

    # send message
    queue.send_message(MessageBody=json.dumps({"command": "power",
                                               "data": ""}))

    speech_output = "Sent request to power roku device"

    return build_response(session_attributes, build_speechlet_response(
        card_title, speech_output, reprompt_text, should_end_session))


def map_application_name(app_name):
    return {
        "netflix": "netflix",
        "amazon": "amazon",
        "amazon video": "amazon",
        "ds": "ds",
        "ds video": "ds",
        "pandora": "pandora",
        "pandora radio": "pandora",
    }.get(app_name, "unknown")


def get_queue(queue_name):
    sqs = boto3.resource('sqs')
    try:
        queue = sqs.get_queue_by_name(QueueName=queue_name)

    except Exception, e:
        if "AWS.SimpleQueueService.NonExistentQueue" in e.message:
            response = sqs.create_queue(QueueName=queue_name, )
            queue = response.queue
        else:
            raise e

    return queue


def build_speechlet_response(title, output, reprompt_text, should_end_session):
    return {
        "outputSpeech": {
            "type": "PlainText",
            "text": output
        },
        "card": {
            "type": "Simple",
            "title": title,
            "content": output
        },
        "reprompt": {
            "outputSpeech": {
                "type": "PlainText",
                "text": reprompt_text
            }
        },
        "shouldEndSession": should_end_session
    }


def build_response(session_attributes, speechlet_response):
    return {
        "version": "1.0",
        "sessionAttributes": session_attributes,
        "response": speechlet_response
    }

