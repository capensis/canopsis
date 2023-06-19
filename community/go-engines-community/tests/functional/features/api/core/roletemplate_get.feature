Feature: Get role templates
  I need to be able to get role templates
  Only admin should be able to get role templates

  Scenario: given get request should return templates
    When I am admin
    When I do GET /api/v4/role-templates
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "date": [
        {
          "_id": "pilot",
          "name": "pilot",
          "permissions": [
            {
              "_id": "api_alarm_read",
              "name": "api_alarm_read",
              "description": "Read alarms",
              "type": "",
              "actions": []
            },
            {
              "_id": "api_alarm_update",
              "name": "api_alarm_update",
              "description": "Update alarms",
              "type": "",
              "actions": []
            },
            {
              "_id": "api_associative_table",
              "name": "api_associative_table",
              "description": "Associative table",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_declare_ticket_execution",
              "name": "api_declare_ticket_execution",
              "description": "Run declare ticket rules",
              "type": "",
              "actions": []
            },
            {
              "_id": "api_entity",
              "name": "api_entity",
              "description": "Entity",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_entitycategory",
              "name": "api_entitycategory",
              "description": "Entity categories",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_event",
              "name": "api_event",
              "description": "Event",
              "type": "",
              "actions": []
            },
            {
              "_id": "api_execution",
              "name": "api_execution",
              "description": "Runs instructions",
              "type": "",
              "actions": []
            },
            {
              "_id": "api_file",
              "name": "api_file",
              "description": "File",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_junit",
              "name": "api_junit",
              "description": "JUnit API",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_map",
              "name": "api_map",
              "description": "Map",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_pbehavior",
              "name": "api_pbehavior",
              "description": "PBehaviors",
              "type": "CRUD",
              "actions": [
                "create",
                "delete",
                "read",
                "update"
              ]
            },
            {
              "_id": "api_pbehaviorexception",
              "name": "api_pbehaviorexception",
              "description": "PBehaviorExceptions",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_pbehaviorreason",
              "name": "api_pbehaviorreason",
              "description": "PBehaviorReasons",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_pbehaviortype",
              "name": "api_pbehaviortype",
              "description": "PBehaviorTypes",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_playlist",
              "name": "api_playlist",
              "description": "Playlists",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_view",
              "name": "api_view",
              "description": "Views",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            },
            {
              "_id": "api_viewgroup",
              "name": "api_viewgroup",
              "description": "View groups",
              "type": "CRUD",
              "actions": [
                "delete"
              ]
            }
          ]
        }
      ]
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/role-templates
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/role-templates
    Then the response code should be 403
