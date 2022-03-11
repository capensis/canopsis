#Feature: Update a metaalarmrule
#  I need to be able to update a metaalarmrule
#  Only admin should be able to update a metaalarmrule
#
#  Scenario: given update request should update metaalarmrule
#    When I am admin
#    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-1:
#    """
#    {
#      "auto_resolve": false,
#      "name": "test-metaalarm-to-update-1-updated",
#      "patterns": null,
#      "type": "complex"
#    }
#    """
#    Then the response code should be 200
#    Then the response body should contain:
#    """
#    {
#      "_id": "test-metaalarm-to-update-1",
#      "auto_resolve": false,
#      "name": "test-metaalarm-to-update-1-updated",
#      "author": "root",
#      "type": "complex"
#    }
#    """
#
#  Scenario: given get request ad no auth user should not allow access
#    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-1
#    Then the response code should be 401
#
#  Scenario: given get request ad auth user by api key without permissions should not allow access
#    When I am noperms
#    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-1
#    Then the response code should be 403
#
#  Scenario: given update request with not exist id should return not found error
#    When I am admin
#    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarmrule-not-found:
#    """
#    {
#      "_id": "test-alarm-get-metaalarm-rule-1",
#      "auto_resolve": false,
#      "config": {
#        "time_interval": {
#          "value": 45,
#          "unit": "s"
#        }
#      },
#      "name": "Test alarm get",
#      "patterns": null,
#      "type": "timebased"
#    }
#    """
#    Then the response code should be 404
#    Then the response body should be:
#    """
#    {
#      "error": "Not found"
#    }
#    """
#
#  Scenario: given update request with wrong entity patterns should return bad request
#    When I am admin
#    Then I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-1:
#    """
#    {
#      "auto_resolve": false,
#      "config": {
#        "time_interval": {
#          "value": 20,
#          "unit": "s"
#        },
#        "entity_patterns": [
#          {
#            "address": "current"
#          }
#        ]
#      },
#      "name": "Test alarm get",
#      "patterns": null,
#      "type": "timebased"
#    }
#    """
#    Then the response code should be 400
#    Then the response body should be:
#    """
#    {
#      "errors": {
#        "config": "entity_patterns config can not be in type timebased.",
#        "config.entity_patterns": "Invalid entity patterns."
#      }
#    }
#    """
#
#  Scenario: given update request with wrong evevnt patterns should return bad request
#    When I am admin
#    When I do PUT /api/v4/cat/metaalarmrules/test-metaalarm-to-update-1:
#    """
#    {
#      "auto_resolve": true,
#      "name": "complex-test-1",
#      "type": "complex",
#      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
#      "config": {
#        "time_interval": {
#          "value": 1,
#          "unit": "m"
#        },
#        "threshold_rate": 1,
#        "event_patterns": [
#          {
#            "address": {
#              "regex_match1": "matched"
#            }
#          }
#        ]
#      }
#    }
#    """
#    Then the response code should be 400
#    Then the response body should be:
#    """
#    {
#      "errors": {
#        "config.event_patterns":"Invalid event patterns."
#      }
#    }
#    """
