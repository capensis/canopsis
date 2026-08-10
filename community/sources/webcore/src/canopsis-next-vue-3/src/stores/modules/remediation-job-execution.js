import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useRemediationJobExecutionStore = defineStore('remediationJobExecution', () => {
  const outputs = ref({});
  const getOutputById = computed(() => id => outputs.value[id] ?? '');

  const create = ({ data } = {}) => request.post(API_ROUTES.remediation.jobExecutions, data);
  const cancel = ({ id } = {}) =>
    request.put(`${API_ROUTES.remediation.jobExecutions}/${id}/cancel`);

  const fetchOutput = async ({ id } = {}) => {
    const output = await request.get(`${API_ROUTES.remediation.jobExecutions}/${id}/output`);

    outputs.value = { ...outputs.value, [id]: output };

    return output;
  };

  return {
    outputs,
    getOutputById,
    create,
    cancel,
    fetchOutput,
  };
});
