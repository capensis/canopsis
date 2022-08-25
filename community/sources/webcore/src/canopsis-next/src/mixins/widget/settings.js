import { widgetToForm, formToWidget } from '@/helpers/forms/widgets/common';

import { queryMixin } from '@/mixins/query';
import { activeViewMixin } from '@/mixins/active-view';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';
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
    activeViewMixin,
    entitiesWidgetMixin,
    entitiesUserPreferenceMixin,
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
      this.form.parameters.mainFilterUpdatedAt = Date.now();
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

        if (widgetId) {
          await this.fetchUserPreference({ id: widgetId });
        }

        await this.fetchActiveView();

        this.$sidebar.hide();
      }
    },
  },
};
