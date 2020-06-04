import { addTo, setField } from '@/helpers/immutable';

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
        const newWidget = {
          ...this.settings.widget,
          ...this.prepareWidgetSettings(),
        };

        const { tabs } = this.activeView;
        const tabIndex = tabs.findIndex(tab => tab._id === this.config.tabId);
        const { widgets } = tabs[tabIndex];
        const widgetIndex = widgets.findIndex(widget => widget._id === newWidget._id);

        if (widgetIndex === -1) {
          const newGridParameters = tabs[tabIndex].widgets.reduce((acc, { gridParameters }) => {
            if (gridParameters.mobile.y >= acc.mobile) {
              acc.mobile = gridParameters.mobile.y + gridParameters.mobile.h + 1;
            }

            if (gridParameters.tablet.y >= acc.tablet) {
              acc.tablet = gridParameters.tablet.y + gridParameters.mobile.h + 1;
            }

            if (gridParameters.desktop.y >= acc.desktop) {
              acc.desktop = gridParameters.desktop.y + gridParameters.mobile.h + 1;
            }

            return acc;
          }, { mobile: 0, tablet: 0, desktop: 0 });

          newWidget.gridParameters.mobile.y = newGridParameters.mobile;
          newWidget.gridParameters.tablet.y = newGridParameters.tablet;
          newWidget.gridParameters.desktop.y = newGridParameters.desktop;
        }

        const userPreference = {
          ...this.userPreference,

          widget_preferences: {
            ...this.userPreference.widget_preferences,
            ...this.settings.widget_preferences,
          },
        };

        const viewData = widgetIndex === -1
          ? addTo(this.activeView, ['tabs', tabIndex, 'widgets'], newWidget)
          : setField(this.activeView, ['tabs', tabIndex, 'widgets', widgetIndex], newWidget);

        await Promise.all([
          this.createUserPreference({ userPreference }),
          this.updateView({ id: this.activeView._id, data: viewData }),
        ]);

        const oldQuery = this.getQueryById(this.widget._id);
        const newQuery = prepareQuery(newWidget, userPreference);

        this.updateQuery({
          id: newWidget._id,
          query: this.prepareWidgetQuery(newQuery, oldQuery),
        });

        this.hideSideBar();
      }
    },
  },
};
