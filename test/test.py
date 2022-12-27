#!/usr/bin/python3

import requests, pytest, uuid, json, sys

headers = {"Content-Type": "application/json"}

def test_todo_set(request, capsys):
    todo = {
        "id": str(uuid.uuid4()),
        "title": "Title",
        "description": "Description"
    }
    response = requests.put(request.config.option.url + "/v1/todo", json=todo, headers=headers)
    pytest.shared = todo
    request.config.cache.set("shared_todo", todo)
    assert response.status_code == 200

def test_todos_get(request, capsys):
    response = requests.get(request.config.option.url + "/v1/todo")
    assert response.status_code == 200 and response.json() != []

def test_todo_get(request, capsys):
    todo = request.config.cache.get("shared_todo", "")
    response = requests.get(request.config.option.url + "/v1/todo/" + todo["id"], headers=headers)
    assert response.status_code == 200 and response.json() == todo


def test_todo_update(request, capsys):
    old_todo_id = request.config.cache.get("shared_todo", "")["id"]
    new_todo = {
        "id": str(uuid.uuid4()),
        "title": "Title",
        "description": "Description"
    }

    #with capsys.disabled():
    #    print("PUT [%s] [%s]" % (old_todo_id, new_todo["id"]))
    response = requests.put(request.config.option.url + "/v1/todo/" + old_todo_id, json=new_todo, headers=headers)
    assert response.status_code == 200


    response = requests.get(request.config.option.url + "/v1/todo/" + old_todo_id, headers=headers)
    assert response.status_code == 400
    
    response = requests.get(request.config.option.url + "/v1/todo/" + new_todo["id"], headers=headers)
    request.config.cache.set("shared_todo", response.json())
    assert response.status_code == 200 and response.json() == new_todo


def test_todo_delete(request, capsys):
    todo_id = request.config.cache.get("shared_todo", "")["id"]

    response = requests.get(request.config.option.url + "/v1/todo/" + todo_id, headers=headers)
    assert response.status_code == 200

    response = requests.delete(request.config.option.url + "/v1/todo/" + todo_id, headers=headers)
    assert response.status_code == 200
    
    response = requests.get(request.config.option.url + "/v1/todo/" + todo_id, headers=headers)
    assert response.status_code == 400
