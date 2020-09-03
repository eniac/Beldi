import re
import time
import json
import boto3
import numpy as np
from pprint import pprint
from argparse import ArgumentParser
from datetime import datetime, timedelta

log_client = boto3.client('logs')
cloudwatch = boto3.resource('cloudwatch')


def get_log_streams(lambda_id):
    group = '/aws/lambda/beldi-dev-{}'.format(lambda_id)
    r = log_client.describe_log_streams(logGroupName=group)
    r = r['logStreams']
    r = [x['logStreamName'] for x in r]
    return r


def delete_logs(lambda_id):
    group = '/aws/lambda/beldi-dev-{}'.format(lambda_id)
    try:
        log_client.delete_log_group(logGroupName=group)
    except:
        pass


def get_logs(lambda_id):
    group = '/aws/lambda/beldi-dev-{}'.format(lambda_id)
    streams = get_log_streams(lambda_id)
    res = []
    for stream in streams:
        r = log_client.get_log_events(logGroupName=group,
                                      logStreamName=stream)
        r = [e['message'].strip() for e in r['events']]
        res += r
    return res


def main():
    parser = ArgumentParser()
    parser.add_argument("--command", required=True)
    parser.add_argument("--config", required=False)
    parser.add_argument("--duration", required=False)
    args = parser.parse_args()
    duration = int(args.duration)
    if args.command == 'clean':
        delete_logs("gctest")
        return
    if args.command == 'run':
        end_time = datetime.utcnow()
        time.sleep(1* 60)  # Wait until metric generated
        metric = cloudwatch.Metric('AWS/Lambda', 'Duration')
        response = metric.get_statistics(
            Dimensions=[
                {
                    'Name': 'FunctionName',
                    'Value': 'beldi-dev-gctest'
                }
            ],
            ExtendedStatistics=['p50', 'p99'],
            StartTime=end_time - timedelta(minutes=duration + 1),
            EndTime=end_time + timedelta(minutes=1),
            Period=60,
            Unit='Milliseconds'
        )
        points = response['Datapoints']
        points.sort(key=lambda x: x['Timestamp'])
        res = []
        for point in points:
            point = point['ExtendedStatistics']
            res.append([point['p50'], point['p99']])
        res = res[1:-1]
        with open('result/gctest/{}.json'.format(args.config), "w") as f:
            json.dump(res, f)
        time.sleep(1 * 60)  # avoid conflicts


if __name__ == "__main__":
    main()
