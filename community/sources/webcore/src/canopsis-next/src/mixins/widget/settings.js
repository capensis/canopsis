import { setField } from '@/helpers/immutable';

import { prepareQuery } from '@/helpers/query';

import queryMixin from '@/mixins/query';
import { entitiesViewMixin } from '@/mixins/entities/view';

import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetSettingsMixin = {
  mixins: [
    queryMixin,
    entitiesViewMixin,
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
        const userPreference = this.getPreparedUserPreference();
        const widget = {
          ...this.settings.widget,
          ...this.prepareWidgetSettings(),
        };

        await Promise.all([
          this.updateUserPreference({ data: userPreference }),
          /**
           * TODO: update widget request
           */
        ]);

        const oldQuery = this.getQueryById(widget._id);
        const newQuery = prepareQuery(widget, userPreference);

        this.updateQuery({
          id: widget._id,
          query: this.prepareWidgetQuery(newQuery, oldQuery),
        });

        this.$sidebar.hide();
      }
    },
  },
};
