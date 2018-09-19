import queryMixin from '@/mixins/query';
import sideBarMixins from '@/mixins/side-bar/side-bar';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { convertUserPreferenceToQuery, convertWidgetToQuery } from '@/helpers/query';

export default {
  props: {
    config: {
      type: Object,
      required: true,
    },
  },
  mixins: [
    queryMixin,
    sideBarMixins,
    entitiesUserPreferenceMixin,
  ],
  computed: {
    widget() {
      return this.config.widget;
    },
  },
  methods: {
    isFormValid() {
      if (this.$validator) {
        return this.$validator.validateAll();
      }

      return true;
    },

    prepareSettingsWidget() {
      return this.settings.widget;
    },

    async submit() {
      const isFormValid = await this.isFormValid();

      if (isFormValid) {
        const widget = {
          ...this.widget,
          ...this.prepareSettingsWidget(),
        };

        const userPreference = {
          ...this.userPreference,
          widget_preferences: {
            ...this.userPreference.widget_preferences,
            ...this.settings.widget_preferences,
          },
        };
        const actions = [this.createUserPreference({ userPreference })];
        if (this.config.isNew) {
          actions.push(this.createWidget({ widget }));
        } else {
          actions.push(this.updateWidget({ widget }));
        }

        await Promise.all(actions);

        this.mergeQuery({
          id: widget._id,
          query: {
            ...convertWidgetToQuery(widget),
            ...convertUserPreferenceToQuery(userPreference),
          },
        });

        this.hideSideBar();
      }
    },
  },
};
