import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const types = {
  FETCH_DYNAMIC_INFO_TEMPLATES: 'FETCH_DYNAMIC_INFO_TEMPLATES',
  FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED: 'FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED',
};

export default {
  namespaced: true,
  state: {
    pending: false,
    dynamicInfoTemplates: [],
  },
  getters: {
    pending: state => state.pending,
    dynamicInfoTemplates: state => state.dynamicInfoTemplates,
  },
  mutations: {
    [types.FETCH_DYNAMIC_INFO_TEMPLATES](state) {
      state.pending = true;
    },
    [types.FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED](state, { templates = [] }) {
      state.pending = false;
      state.dynamicInfoTemplates = templates;
    },
  },
  actions: {
    async fetchList({ commit }) {
      try {
        commit(types.FETCH_PBEHAVIOR_REASONS);

        const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

        commit(types.FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED, { templates });
      } catch (err) {
        console.error(err);
      }
    },

    async create({ dispatch }, { template }) {
      const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

      if (templates) {
        return dispatch('updateTemplate', { template });
      }

      return request.post(API_ROUTES.dynamicInfoTemplates, { templates: [template] });
    },

    async update({ dispatch }, { template }) {
      const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

      if (!templates) {
        return dispatch('createTemplate', { template });
      }

      return request.put(API_ROUTES.dynamicInfoTemplates, {
        templates: templates.map((v) => {
          if (v._id === template._id) {
            return template;
          }

          return v;
        }),
      });
    },

    async delete(context, { id }) {
      const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

      if (templates) {
        await request.put(API_ROUTES.dynamicInfoTemplates, {
          templates: templates.filter(v => v._id !== id),
        });
      }
    },
  },
};
