Feature: get check-incident-exists request
  I need to be able to get check-incident-exists data
  Only admin should be able to get check-incident-exists

  Scenario: POST unauthorized
    When I do POST /api/v4/cat/ocws-noc/check-incident-exists
    Then the response code should be 401

  Scenario: POST without permissions
    When I am noperms
    When I do POST /api/v4/cat/ocws-noc/check-incident-exists
    Then the response code should be 403

  Scenario: POST check-incident-exists success
    When I am admin
    Then I do POST /api/v4/cat/ocws-noc/check-incident-exists:
    ```
    {
        "ci_provenance": "company-incident-exists-1-provenance-1",
        "company_sys_id_snow": "company-incident-exists-1",
        "machine_snow": "",
        "service_contract_snow": "contract-incident-exists-1"
    }
    ```
    Then the response code should be 200
    Then the response body should contain:
    ```
    {
        "found": 1,
        "tickets": [
            {
            "number": "OCWS_INC01976",
            "state": "1",
            "sys_created_on": "2029-03-29 11:11:04",
            "sys_id": "2000-PD0809",
            "u_incident_start_time": "2029-07-30 06:13:37"
            }
        ]
    }
    ```

  Scenario: POST check-incident-exists contract not found 
    When I am admin
    Then I do POST /api/v4/cat/ocws-noc/check-incident-exists:
    ```
    {
        "ci_provenance": "company-incident-exists-2-provenance-1",
        "company_sys_id_snow": "company-incident-exists-2",
        "machine_snow": "",
        "service_contract_snow": "contract-incident-exists-2"
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
