#!/bin/bash

d=$(cd "$(dirname "$0")"; pwd)

bash deploy_x.sh notifier notifier

bash deploy_x.sh notifier nm/notifier notifier nm