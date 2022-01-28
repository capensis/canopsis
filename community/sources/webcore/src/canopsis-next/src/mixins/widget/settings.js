import { widgetToForm, formToWidget } from '@/helpers/forms/widgets/common';

import { queryMixin } from '@/mixins/query';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';
import { entitiesViewTabMixin } from '@/mixins/entities/view/tab';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

export const widgetSettingsMixin = {
  $_veeValidate: {
    validator: 'new',
  },
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },
  mixins: [
    queryMixin,
    entitiesWidgetMixin,
    entitiesViewTabMixin,
    confirmableModalMixinCreator({ field: 'form', closeMethod: '$sidebar.hide' }),
  ],
  data() {
    return {
      form: widgetToForm(this.sidebar.config?.widget),
    };
  },
  computed: {
    config() {
      return this.sidebar.config ?? {};
    },

    widget() {
      return this.config.widget ?? {};
    },

    duplicate() {
      return this.config.duplicate;
    },
  },
  methods: {
    /**
     * Update main filter updated at value. We are using this value for checking which filter was changed later
     */
    updateMainFilterUpdatedAt() {
      this.form.parameters.main_filter_updated_at = Date.now();
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
     * Submit settings form
     *
     * @returns {Promise<void>}
     */
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const { _id: widgetId, tab: tabId } = this.widget;
        const data = formToWidget(this.form);

        data.tab = tabId;

        if (this.duplicate) {
          await this.copyWidget({ id: widgetId, data });
        } else if (widgetId) {
          await this.updateWidget({ id: widgetId, data });
        } else {
          await this.createWidget({ data });
        }

        await this.fetchViewTab({ id: tabId });

        /**
         * TODO: update widget request
         */

        /*        if (newWidget._id) {
          const oldQuery = this.getQueryById(newWidget._id);
          const newQuery = prepareQuery(newWidget, userPreference);

          this.updateQuery({
            id: data._id,
            query: this.prepareWidgetQuery(newQuery, oldQuery),
          });
        } */

        this.$sidebar.hide();
      }
    },
  },
};
