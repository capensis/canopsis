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
    templates: [],
  },
  getters: {
    pending: state => state.pending,
    templates: state => state.templates,
  },
  mutations: {
    [types.FETCH_DYNAMIC_INFO_TEMPLATES](state) {
      state.pending = true;
    },
    [types.FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED](state, { templates = [] }) {
      state.pending = false;
      state.templates = templates;
    },
  },
  actions: {
    /**
     *
     * @param commit
     * @returns {Promise<void>}
     */
    async fetchList({ commit }) {
      try {
        commit(types.FETCH_DYNAMIC_INFO_TEMPLATES);

        const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

        commit(types.FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED, { templates });
      } catch (err) {
        console.error(err);
      }
    },

    /**
     *
     * @param dispatch
     * @param template
     * @returns {Promise<AxiosResponse<any>|*>}
     */
    async create({ dispatch }, { data }) {
      const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

      if (templates) {
        return dispatch('updateTemplate', { template: data });
      }

      return request.post(API_ROUTES.dynamicInfoTemplates, { templates: [data] });
    },

    /**
     *
     * @param dispatch
     * @param template
     * @returns {Promise<AxiosResponse<any>|*>}
     */
    async update({ dispatch }, { data }) {
      const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

      if (!templates) {
        return dispatch('createTemplate', { template: data });
      }

      return request.put(API_ROUTES.dynamicInfoTemplates, {
        templates: templates.map((v) => {
          if (v._id === data._id) {
            return data;
          }

          return v;
        }),
      });
    },

    /**
     *
     * @param context
     * @param id
     * @returns {Promise<void>}
     */
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
