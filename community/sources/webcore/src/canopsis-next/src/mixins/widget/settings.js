import { VALIDATION_DELAY } from '@/constants';

import { widgetToForm, formToWidget } from '@/helpers/entities/widget/form';

import { queryMixin } from '@/mixins/query';
import { activeViewMixin } from '@/mixins/active-view';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { submittableMixinCreator } from '@/mixins/submittable';

export const widgetSettingsMixin = {
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
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
    submittableMixinCreator(),
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
          await this.createWidget({ data });
        } else if (widgetId) {
          await this.updateWidget({ id: widgetId, data });
        } else {
          await this.createWidget({ data });
        }

        if (data.parameters.mainFilter && this.userPreference.content.mainFilter === data.parameters.mainFilter) {
          await this.updateContentInUserPreference({ mainFilter: null });
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
