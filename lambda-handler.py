from flask import Flask, request, abort, make_response
import requests
from urllib.parse import urlparse
from urllib.parse import parse_qs

app = Flask(__name__)


@app.route('/', methods=['POST'])
def get_webhook():
    if request.method == 'POST':
        # Get the request event from the 'POST' call
        event = request.json

        # Get the object context
        object_context = event["getObjectContext"]

        # print("Custom Rage Param in Request:", rangeValue)

        # Get the presigned URL
        # Used to fetch the original object from MinIO
        s3_url = object_context["inputS3Url"]

        # Extract the route and request tokens from the input context
        request_route = object_context["outputRoute"]
        request_token = object_context["outputToken"]

        captured_value = ""
        # Work Around to make the range request.
        try:
            req_obj_url = event["userRequest"]['url']
            parsed_url = urlparse(req_obj_url)
            captured_value = parse_qs(parsed_url.query)['my-range'][0]
        except:
            print("An exception occurred")

        headers = {}
        if captured_value:
            headers = {'Range': "bytes=" + captured_value}

        # Get the original S3 object using the presigned URL
        r = requests.get(url=s3_url, headers=headers)
        # r.headers["Range"] = "bytes=" + captured_value
        # Range: bytes=0-1023

        original_object = r.content.decode('utf-8')

        print("Original:", original_object)
        # Transform the text in the object by swapping the case of each char
        transformed_object = original_object.swapcase()

        print("Transformed:", transformed_object)
        # Return the object back to Object Lambda, with required headers
        # This sends the transformed data to MinIO
        # and then to the user
        resp = make_response(transformed_object, 200)
        resp.headers['x-amz-request-route'] = request_route
        resp.headers['x-amz-request-token'] = request_token
        resp.headers['x-max-transformed-by'] = 'prakash'
        return resp

    else:
        abort(400)


if __name__ == '__main__':
    app.run()
