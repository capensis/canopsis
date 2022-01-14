import { setField } from '@/helpers/immutable';

import { prepareQuery } from '@/helpers/query';
import { formToWidget } from '@/helpers/forms/widgets/common';

import { queryMixin } from '@/mixins/query';

import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetSettingsMixin = {
  mixins: [
    queryMixin,
    entitiesUserPreferenceMixin,
    confirmableModalMixinCreator({ field: 'settings', closeMethod: '$sidebar.hide' }),
  ],
  methods: {
    /**
     * Validate settings form
     *
     * @returns {boolean|Promise<boolean>}
     */
    isFormValid() {
      if (this.$validator) {
        return this.$validator.validateAll();
      }

      return true;
    },

    /**
     * We can customize widgets preparation by replacing the methods in the component
     *
     * @returns {Object}
     */
    prepareWidgetSettings() {
      return this.settings.widget;
    },

    /**
     * We can customize widget query for updating after saving by replacing the methods in the component
     *
     * @param {Object} newQuery
     * @returns {Object}
     */
    prepareWidgetQuery(newQuery) {
      return newQuery;
    },

    /**
     * Get prepared userPreferences for request sending
     *
     * @returns {Object}
     */
    getPreparedUserPreference() {
      return setField(this.userPreference, 'content', value => ({
        ...value,
        ...this.settings.userPreferenceContent,
      }));
    },

    /**
     * Submit settings form
     *
     * @returns {Promise<void>}
     */
    async submit() {
      const isFormValid = await this.isFormValid();

      if (isFormValid) {
        const data = formToWidget(this.settings.widget);

        /**
         * TODO: update widget request
         */

        if (data._id) {
          const userPreference = this.getPreparedUserPreference();

          await this.updateUserPreference({ data: userPreference });

          const oldQuery = this.getQueryById(data._id);
          const newQuery = prepareQuery(data, userPreference);

          this.updateQuery({
            id: data._id,
            query: this.prepareWidgetQuery(newQuery, oldQuery),
          });
        }

        this.$sidebar.hide();
      }
    },
  },
};
