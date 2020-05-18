import { setField } from '@/helpers/immutable';

import queryMixin from '@/mixins/query';
import sideBarMixins from '@/mixins/side-bar/side-bar';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import { prepareQuery } from '@/helpers/query';

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
    entitiesViewMixin,
    entitiesUserPreferenceMixin,
  ],
  computed: {
    activeView() {
      return this.config.viewId ? this.getViewById(this.config.viewId) : this.view;
    },

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

    prepareWidgetSettings() {
      return this.settings.widget;
    },

    prepareWidgetQuery(newQuery) {
      return newQuery;
    },

    async submit() {
      const isFormValid = await this.isFormValid();

      if (isFormValid) {
        const widget = {
          ...this.settings.widget,
          ...this.prepareWidgetSettings(),
        };

        const userPreference = {
          ...this.userPreference,

          widget_preferences: {
            ...this.userPreference.widget_preferences,
            ...this.settings.widget_preferences,
          },
        };

        const viewData = setField(this.activeView, ['tabs', this.config.tabId, widget._id], widget);

        await Promise.all([
          this.createUserPreference({ userPreference }),
          this.updateView({ id: this.activeView._id, data: viewData }),
        ]);

        const oldQuery = this.getQueryById(this.widget._id);
        const newQuery = prepareQuery(widget, userPreference);

        this.updateQuery({
          id: widget._id,
          query: this.prepareWidgetQuery(newQuery, oldQuery),
        });

        this.hideSideBar();
      }
    },
  },
};
