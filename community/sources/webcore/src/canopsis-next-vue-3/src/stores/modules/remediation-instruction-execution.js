import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useRemediationInstructionExecutionStore = defineStore(
  'remediationInstructionExecution',
  () => {
    const fetchItemWithoutStore = ({ id, params } = {}) => request.get(`${API_ROUTES.remediation.executions}/${id}`, {
        params,
      });
    const fetchPausedListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.remediation.pausedExecutions, {
        params,
      });
    const fetchStatusesListWithoutStore = () => request.get(API_ROUTES.remediation.executionStatuses);
    const create = ({ data } = {}) => request.post(API_ROUTES.remediation.executions, data);
    const update = ({ path, id, data } = {}) => request.put(`${API_ROUTES.remediation.executions}/${id}/${path}`, data, {
        fullResponse: true,
      });
    const cancel = ({ id } = {}) => update({ path: 'cancel', id });
    const nextOperation = ({ id } = {}) => update({ path: 'next', id });
    const nextStep = ({ id, data } = {}) => update({ path: 'next-step', id, data });
    const pause = ({ id } = {}) => update({ path: 'pause', id });
    const previousOperation = ({ id } = {}) => update({ path: 'previous', id });
    const resume = ({ id } = {}) => update({ path: 'resume', id });
    const fetchAlarmExecutionsWithoutStore = ({ alarmId, params } = {}) => request.get(`${API_ROUTES.remediation.alarmExecutions}/${alarmId}`, { params });
    const read = ({ id } = {}) => request.put(`${API_ROUTES.remediation.executions}/${id}/read`);

    return {
      fetchItemWithoutStore,
      fetchPausedListWithoutStore,
      fetchStatusesListWithoutStore,
      create,
      update,
      cancel,
      nextOperation,
      nextStep,
      pause,
      previousOperation,
      resume,
      fetchAlarmExecutionsWithoutStore,
      read,
    };
  },
);
