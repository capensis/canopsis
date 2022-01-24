import { setField } from '@/helpers/immutable';

import { prepareQuery } from '@/helpers/query';
import { widgetToForm, formToWidget } from '@/helpers/forms/widgets/common';

import { queryMixin } from '@/mixins/query';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

export const widgetSettingsMixin = {
  mixins: [
    queryMixin,
    entitiesWidgetMixin,
    entitiesUserPreferenceMixin, // TODO: remove it
    confirmableModalMixinCreator({ field: 'settings', closeMethod: '$sidebar.hide' }), // TODO: change field to form
  ],
  computed: {
    config() {
      return this.sidebar.config ?? {};
    },

    widget() {
      return this.config.widget;
    },
  },
  data() {
    return {
      form: widgetToForm(this.sidebar.config?.widget),
    };
  },
  methods: {
    /**
     * Validate settings form
     *
     * @returns {boolean|Promise<boolean>}
     */
    isFormValid() {
      return this.$validator?.validateAll() ?? true;
    },

    /**
     * We can customize widgets preparation by replacing the methods in the component
     *
     * @returns {Object}
     */
    prepareWidgetSettings() { // TODO: remove it
      return this.settings.widget;
    },

    /**
     * We can customize widget query for updating after saving by replacing the methods in the component
     *
     * @param {Object} newQuery
     * @returns {Object}
     */
    prepareWidgetQuery(newQuery) { // TODO: refactor it
      return newQuery;
    },

    /**
     * Get prepared userPreferences for request sending
     *
     * @returns {Object}
     */
    getPreparedUserPreference() { // TODO: remove it
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
        let method = this.createWidget;

        if (this.widget) {
          if (this.widget._id) {
            method = this.updateWidget;
          } else {
            method = this.copyWidget;
          }
        }

        const newWidget = await method({ data });

        /**
         * TODO: update widget request
         */

        if (newWidget._id) {
          const userPreference = this.getPreparedUserPreference();

          await this.updateUserPreference({ data: userPreference });

          const oldQuery = this.getQueryById(newWidget._id);
          const newQuery = prepareQuery(newWidget, userPreference);

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
