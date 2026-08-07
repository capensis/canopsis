import { isEqual } from 'lodash';
import { mapActions } from 'pinia';

import { MODALS } from '@/constants';

import { formToWidget, widgetToForm } from '@/helpers/entities/widget/form';

import { activeViewMixin } from '@/mixins/active-view';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

import { useWidgetStore } from '@/stores/modules/widget';

// TODO: remove after every widget settings sidebar is rewritten to Composition API with useWidgetSettings().
export const widgetSettingsMixin = {
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },

  mixins: [activeViewMixin, entitiesUserPreferenceMixin],

  data() {
    const form = widgetToForm(this.sidebar.config?.widget);

    return {
      form,
      originalForm: widgetToForm(this.sidebar.config?.widget),
      hasChanges: false,
      submitting: false,
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

  created() {
    this.registerWatchOnceForForm();
  },

  methods: {
    ...mapActions(useWidgetStore, {
      createWidget: 'create',
      updateWidget: 'update',
    }),

    hideSettingsSidebar() {
      this.$sidebar.hide({ id: this.sidebar.id });
    },

    close() {
      if (!this.hasChanges || isEqual(this.form, this.originalForm)) {
        this.hideSettingsSidebar();

        return;
      }

      this.$modals.show({
        name: MODALS.confirmation,
        dialogProps: {
          persistent: true,
        },
        config: {
          action: this.hideSettingsSidebar,
        },
      });
    },

    registerWatchOnceForForm() {
      const unwatch = this.$watch(
        () => this.form,
        () => {
          this.hasChanges = true;

          this.$nextTick(() => unwatch());
        },
        { deep: true },
      );
    },

    scrollToFirstError() {
      this.$nextTick(() => {
        this.$el.querySelector('.v-input--error')?.scrollIntoView({
          block: 'center',
          behavior: 'smooth',
        });
      });
    },

    /**
     * Submit settings form.
     */
    async submit() {
      const isFormValid = await (this.$refs.settings?.validate?.() ?? true);

      if (!isFormValid) {
        this.scrollToFirstError();

        return;
      }

      try {
        this.submitting = true;

        const { _id: widgetId, tab: tabId } = this.widget;
        const data = formToWidget(this.form);

        data.tab = tabId;

        if (this.duplicate || !widgetId) {
          await this.createWidget({ data });
        } else {
          await this.updateWidget({ id: widgetId, data });
        }

        if (
          data.parameters.mainFilter &&
          this.userPreference.content.mainFilter === data.parameters.mainFilter
        ) {
          await this.updateContentInUserPreference({ mainFilter: null });
        }

        if (widgetId) {
          await this.fetchUserPreference({ id: widgetId });
        }

        await this.fetchActiveView();

        this.hideSettingsSidebar();
      } finally {
        this.submitting = false;
      }
    },
  },
};
