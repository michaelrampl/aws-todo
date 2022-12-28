# Todo API

This repository contains 2 different binaries.

## Standalone

This binary implements the fiber framework and can be used as a standalone version for servers.

## Lambda

This binary is intended to be used as AWS Lambda function behind an API Gateway.

## Structure

Both binaries access the same handler but inject a different database (all based on a common interface). At the time being, either mongodb or dynamodb are possible backends.

## Public AWS API

The AWS Lambda is publicly available on `https://d4rarz3yu2.execute-api.us-east-1.amazonaws.com/production`


## Tests

Inside the testfolder there are 2 shell scripts utilizing test.py (uses pytest) to execute 15 tests on both a local or the aws instance
This may look like the following snipped

```
====================================== test session starts ======================================
platform linux -- Python 3.10.9, pytest-7.1.2, pluggy-1.0.0 -- /usr/bin/python3
cachedir: .pytest_cache
rootdir: /home/michael/Documents/Resume/sda_2022/aws-todo/test
plugins: dependency-0.5.1, anyio-3.6.2
collected 15 items                                                                              

test.py::test_todo_set_ok PASSED                                                          [  6%]
test.py::test_todo_set_existing_id PASSED                                                 [ 13%]
test.py::test_todo_set_no_id PASSED                                                       [ 20%]
test.py::test_todo_set_invalid_id PASSED                                                  [ 26%]
test.py::test_todo_set_no_title PASSED                                                    [ 33%]
test.py::test_todo_set_title_min PASSED                                                   [ 40%]
test.py::test_todo_set_title_max PASSED                                                   [ 46%]
test.py::test_todo_set_no_description PASSED                                              [ 53%]
test.py::test_todo_get_all PASSED                                                         [ 60%]
test.py::test_todo_get_invalid_id PASSED                                                  [ 66%]
test.py::test_todo_get_ok PASSED                                                          [ 73%]
test.py::test_todo_update_ok PASSED                                                       [ 80%]
test.py::test_todo_update_invalid_id PASSED                                               [ 86%]
test.py::test_todo_delete_ok PASSED                                                       [ 93%]
test.py::test_todo_delete_invalid_id PASSED                                               [100%]

====================================== 15 passed in 10.01s ======================================
```