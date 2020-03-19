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
     * @param context
     * @param template
     * @returns {Promise<void>}
     */
    async create({ commit }, { data }) {
      const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);
      const method = templates ? 'put' : 'post';

      const newTemplates = [...templates, data];
      const success = await request[method](API_ROUTES.dynamicInfoTemplates, { templates: newTemplates });

      if (success) {
        commit(types.FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED, { templates: newTemplates });
      }
    },

    /**
     *
     * @param context
     * @param template
     * @returns {Promise<void>}
     */
    async update({ commit }, { data }) {
      const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

      if (templates) {
        const newTemplates = templates.map((v) => {
          if (v._id === data._id) {
            return data;
          }

          return v;
        });

        const success = await request.put(API_ROUTES.dynamicInfoTemplates, {
          templates: newTemplates,
        });

        if (success) {
          commit(types.FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED, { templates: newTemplates });
        }
      }
    },

    /**
     *
     * @param context
     * @param id
     * @returns {Promise<void>}
     */
    async delete({ commit }, { id }) {
      const { templates } = await request.get(API_ROUTES.dynamicInfoTemplates);

      if (templates) {
        const newTemplates = templates.filter(v => v._id !== id);

        const success = await request.put(API_ROUTES.dynamicInfoTemplates, {
          templates: newTemplates,
        });

        if (success) {
          commit(types.FETCH_DYNAMIC_INFO_TEMPLATES_COMPLETED, { templates: newTemplates });
        }
      }
    },
  },
};
