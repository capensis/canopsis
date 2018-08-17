import queryMixin from '@/mixins/query';
import entitiesWidgetMixin from '@/mixins/entities/widget';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { convertUserPreferenceToQuery, convertWidgetToQuery } from '@/helpers/query';

export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  mixins: [
    queryMixin,
    entitiesWidgetMixin,
    entitiesUserPreferenceMixin,
  ],
  methods: {
    async submit() {
      const widget = {
        ...this.widget,
        ...this.settings.widget,
      };

      const userPreference = {
        ...this.userPreference,
        widget_preferences: {
          ...this.userPreference.widget_preferences,
          ...this.settings.widget_preferences,
        },
      };

      const actions = [this.createUserPreference({ userPreference })];

      if (this.isNew) {
        actions.push(this.createWidget({ widget }));
      } else {
        actions.push(this.updateWidget({ widget }));
      }

      await Promise.all(actions);

      await this.mergeQuery({
        id: widget.id,
        query: {
          ...convertWidgetToQuery(widget),
          ...convertUserPreferenceToQuery(userPreference),
        },
      });

      this.$emit('closeSettings');
    },
  },
};
