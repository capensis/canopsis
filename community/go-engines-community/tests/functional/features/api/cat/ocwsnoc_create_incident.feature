Feature: post create-incident request
  I need to be able to post create-incident data
  Only admin should be able to post create-incident

  Scenario: POST unauthorized
    When I do POST /api/v4/cat/ocws-noc/create-incident
    Then the response code should be 401

  Scenario: POST without permissions
    When I am noperms
    When I do POST /api/v4/cat/ocws-noc/create-incident
    Then the response code should be 403

  Scenario: POST create-incident success
    When I am admin
    Then I do POST /api/v4/cat/ocws-noc/create-incident:
    ```
    {
        "ci_provenance": "company-create-incident-1-provenance-1",
        "company_sys_id_snow": "company-create-incident-1",
        "machine_snow": "",
        "service_contract_snow": "contract-create-incident-1",
        "impact": "1",
        "urgency": "2",
        "branch_snow": "location-create-incident-1",
        "description": "[KEYWORD] - #9803 TICKET TITLÈ",
        "resource": "resource-create-incident-1",
        "component": "component-create-incident-1",
        "connector": "snow",
        "connector_name": "SNow",
        "site": "site-create-incident-1",
        "opened_by": "user-create-incident-1-opened_by-1",
        "assignment_group_snow": "group-create-incident-1-assgnment-1",
        "ci_alert_list_sms_snow": "alert_list-sms"
    }
    ```
    Then the response code should be 200
    Then the response body should contain:
    ```
    {
        "event_type": "assocticket",
        "author": "user-create-incident-1-opened_by-1",
        "resource": "resource-create-incident-1",
        "component": "component-create-incident-1",
        "connector": "snow",
        "connector_name": "SNow",
        "source_type": "resource",
        "ticket": "OCD_INC0411241"
    }
    ```

  Scenario: POST create-incident contract not found 
    When I am admin
    Then I do POST /api/v4/cat/ocws-noc/create-incident:
    ```
    {
        "ci_provenance": "company-create-incident-2-provenance-1",
        "company_sys_id_snow": "company-create-incident-2",
        "machine_snow": "",
        "service_contract_snow": "contract-create-incident-2",
        "impact": "1",
        "urgency": "2"
    }
    ```
    Then the response code should be 404
    Then the response body should be:
    ```
    {
        "errors": {
            "service_contract_snow": "u_contract_owner not found"
        }
    }
    ```
