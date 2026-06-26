import { defineStore } from 'pinia';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const usePbehaviorStore = defineStore('pbehavior', () => {
  const fetchPbehaviorsListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.pbehavior.pbehaviors, { params });
  const fetchPbehaviorReasonsListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.pbehavior.reasons, { params });
  const fetchPbehaviorTypesListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.pbehavior.types, { params });
  const fetchPbehaviorTypesFieldList = ({ params = { paginate: false, with_hidden: true } } = {}) => (
    request.get(API_ROUTES.pbehavior.types, { params })
  );
  const fetchNextPbehaviorTypePriority = () => request.get(API_ROUTES.pbehavior.nextTypesPriority);

  return {
    fetchPbehaviorsListWithoutStore,
    fetchPbehaviorReasonsListWithoutStore,
    fetchPbehaviorTypesListWithoutStore,
    fetchPbehaviorTypesFieldList,
    fetchNextPbehaviorTypePriority,
  };
});
