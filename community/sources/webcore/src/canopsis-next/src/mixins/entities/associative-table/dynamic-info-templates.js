import { get } from 'lodash';

import { ASSOCIATIVE_TABLES_NAMES } from '@/constants';

import { entitiesAssociativeTableMixin } from './index';

export const entitiesDynamicInfoTemplatesMixin = {
  mixins: [entitiesAssociativeTableMixin],
  methods: {
    /**
     * Fetch dynamic info templates
     *
     * @returns {Promise<[]>}
     */
    async fetchDynamicInfoTemplatesList() {
      const content = await this.fetchAssociativeTable({
        name: ASSOCIATIVE_TABLES_NAMES.dynamicInfoTemplates,
      });

      return get(content, 'templates', []);
    },

    /**
     * Main update dynamic info templates method
     *
     * @param {Object} data
     * @returns {Promise<[]>}
     */
    async updateDynamicInfoTemplateMethod(data) {
      const content = await this.updateAssociativeTable({
        name: ASSOCIATIVE_TABLES_NAMES.dynamicInfoTemplates,
        data,
      });

      return get(content, 'templates', []);
    },

    /**
     * Add new dynamic info template into templates
     *
     * @param {DynamicInfoTemplate} data
     * @returns {Promise<[]>}
     */
    async createDynamicInfoTemplate({ data }) {
      const templates = await this.fetchDynamicInfoTemplatesList();

      return this.updateDynamicInfoTemplateMethod({
        templates: [...templates, data],
      });
    },

    /**
     * Update existing dynamic info template in templates
     *
     * @param {string} id
     * @param {DynamicInfoTemplate} data
     * @returns {Promise<[]>}
     */
    async updateDynamicInfoTemplate({ id, data }) {
      const templates = await this.fetchDynamicInfoTemplatesList();

      return this.updateDynamicInfoTemplateMethod({
        templates: templates.map(item => (item._id === id ? data : item)),
      });
    },

    /**
     * Remove dynamic info template from templates
     *
     * @param {string} id
     * @returns {Promise<[]>}
     */
    async removeDynamicInfoTemplate({ id }) {
      const templates = await this.fetchDynamicInfoTemplatesList();

      return this.updateDynamicInfoTemplateMethod({
        templates: templates.filter(item => item._id !== id),
      });
    },
  },
};
