import { API_ROUTES } from '@/config';
import request from '@/services/request';
import { createCRUDStore } from '@/stores/plugins/crud';

export const usePbehaviorStore = createCRUDStore({
  id: 'pbehavior',
  route: API_ROUTES.pbehavior.pbehaviors,
  dataPreparer: response => response?.data,
  withFetchingParams: true,
  withMeta: true,
  withWithoutStore: true,
  actions: {
    fetchPbehaviorReasonsListWithoutStore({ params } = {}) {
      return request.get(API_ROUTES.pbehavior.reasons, { params });
    },

    fetchPbehaviorTypesListWithoutStore({ params } = {}) {
      return request.get(API_ROUTES.pbehavior.types, { params });
    },

    fetchPbehaviorTypesFieldList({ params = { paginate: false, with_hidden: true } } = {}) {
      return request.get(API_ROUTES.pbehavior.types, { params });
    },

    fetchNextPbehaviorTypePriority() {
      return request.get(API_ROUTES.pbehavior.nextTypesPriority);
    },

    bulkCreate({ data } = {}) {
      return request.post(API_ROUTES.pbehavior.bulkPbehaviors, data);
    },

    fetchItemWithoutStore({ id } = {}) {
      return request.get(`${API_ROUTES.pbehavior.pbehaviors}/${id}`);
    },

    bulkUpdate({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkPbehaviors, data);
    },

    bulkEnable({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkPbehaviorsEnable, data);
    },

    bulkDisable({ data } = {}) {
      return request.put(API_ROUTES.pbehavior.bulkPbehaviorsDisable, data);
    },

    bulkRemove({ data } = {}) {
      return request.delete(API_ROUTES.pbehavior.bulkPbehaviors, { data });
    },

    fetchPbehaviorsByEntityIdWithoutStore({ id, params = {} } = {}) {
      return request.get(API_ROUTES.pbehavior.entities, { params: { _id: id, ...params } });
    },

    fetchPbehaviorsCalendarWithoutStore({ params } = {}) {
      return request.get(API_ROUTES.pbehavior.calendar, { params });
    },

    fetchEntitiesPbehaviorsCalendarWithoutStore({ params } = {}) {
      return request.get(API_ROUTES.pbehavior.entitiesCalendar, { params });
    },

    bulkCreateEntityPbehaviors({ data } = {}) {
      return request.post(API_ROUTES.pbehavior.bulkEntityPbehaviors, data);
    },

    bulkRemoveEntityPbehaviors({ data } = {}) {
      return request.delete(API_ROUTES.pbehavior.bulkEntityPbehaviors, { data });
    },
  },
});
