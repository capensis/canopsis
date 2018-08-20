import queryMixin from '@/mixins/query';
import sideBarMixins from '@/mixins/side-bar/side-bar';
import entitiesWidgetMixin from '@/mixins/entities/widget';
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
    entitiesWidgetMixin,
    entitiesUserPreferenceMixin,
  ],
  computed: {
    widget() {
      return this.config.widget;
    },
  },
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

      if (this.config.isNew) {
        actions.push(this.createWidget({ widget }));
      } else {
        actions.push(this.updateWidget({ widget }));
      }

      await Promise.all(actions);

      this.mergeQuery({
        id: widget.id,
        query: {
          ...convertWidgetToQuery(widget),
          ...convertUserPreferenceToQuery(userPreference),
        },
      });

      this.hideSideBar();
    },
  },
};
