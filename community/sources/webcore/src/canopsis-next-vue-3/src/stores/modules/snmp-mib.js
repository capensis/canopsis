import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';
import request from '@/services/request';

export const useSnmpMibStore = defineStore('snmpMib', () => {
  const fetchList = ({ params = {} } = {}) => request.get(API_ROUTES.snmpMib, { params });
  const upload = ({ data } = {}) => request.post(API_ROUTES.snmpMib, data);

  return {
    fetchList,
    upload,
  };
});
