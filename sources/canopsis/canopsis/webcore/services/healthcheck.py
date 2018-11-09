# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from bottle import request

from canopsis.healthcheck.manager import HealthcheckManager
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR


def exports(ws):

    healthcheckManager = HealthcheckManager(
        *HealthcheckManager.provide_default_basics()
    )

    @ws.application.get(
        '/api/v2/healthcheck/'
    )
    def get_healthcheck():
        """
        Get healthcheck status report.

        :returns: <Healthcheck>
        """
        criticals = request.query.criticals.split(',') or None
        health_obj = healthcheckManager.check(criticals=criticals)
        if health_obj is None:
            return gen_json_error({'description': 'Healthcheck is empty !'},
                                  HTTP_ERROR)

        return gen_json(health_obj)
