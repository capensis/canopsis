import { addTo, setField } from '@/helpers/immutable';

import { prepareQuery } from '@/helpers/query';
import { getNewWidgetGridParametersY } from '@/helpers/grid-layout';
import { viewToRequest } from '@/helpers/forms/view';

import queryMixin from '@/mixins/query';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import entitiesViewMixin from '@/mixins/entities/view';

import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { entitiesUserPreferenceMixin } from '@/mixins/entities/user-preference';

export const widgetSettingsMixin = {
  props: {
    config: {
      type: Object,
      required: true,
    },
  },
  mixins: [
    queryMixin,
    sideBarMixin,
    entitiesViewMixin,
    entitiesUserPreferenceMixin,
    confirmableModalMixinCreator({ field: 'settings', closeMethod: 'hideSideBar' }),
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
    /**
     * Validate settings form
     *
     * @returns {boolean|Promise<boolean>}
     */
    isFormValid() {
      if (this.$validator) {
        return this.$validator.validateAll();
      }

      return true;
    },

    /**
     * We can customize widgets preparation by replacing the methods in the component
     *
     * @returns {Object}
     */
    prepareWidgetSettings() {
      return this.settings.widget;
    },

    /**
     * We can customize widget query for updating after saving by replacing the methods in the component
     *
     * @param {Object} newQuery
     * @returns {Object}
     */
    prepareWidgetQuery(newQuery) {
      return newQuery;
    },

    /**
     * Get prepared userPreferences for request sending
     *
     * @returns {Object}
     */
    getPreparedUserPreference() {
      return setField(this.userPreference, 'content', value => ({
        ...value,
        ...this.settings.userPreferenceContent,
      }));
    },

    /**
     * Get prepared view and widgets objects for request sending
     *
     * @returns {{view: Object, widget: Object}}
     */
    getPreparedViewAndWidget() {
      const preparedWidget = {
        ...this.settings.widget,
        ...this.prepareWidgetSettings(),
      };

      const { tabs = [] } = this.activeView;
      const tabIndex = tabs.findIndex(tab => tab._id === this.config.tabId);
      const { widgets = [] } = tabs[tabIndex];
      const widgetIndex = widgets.findIndex(widget => widget._id === preparedWidget._id);

      if (widgetIndex === -1) {
        const newGridParametersY = getNewWidgetGridParametersY(tabs[tabIndex].widgets);

        preparedWidget.grid_parameters.mobile.y = newGridParametersY.mobile;
        preparedWidget.grid_parameters.tablet.y = newGridParametersY.tablet;
        preparedWidget.grid_parameters.desktop.y = newGridParametersY.desktop;
      }

      const preparedView = widgetIndex === -1
        ? addTo(this.activeView, ['tabs', tabIndex, 'widgets'], preparedWidget)
        : setField(this.activeView, ['tabs', tabIndex, 'widgets', widgetIndex], preparedWidget);

      return {
        view: preparedView,
        widget: preparedWidget,
      };
    },

    /**
     * Submit settings form
     *
     * @returns {Promise<void>}
     */
    async submit() {
      const isFormValid = await this.isFormValid();

      if (isFormValid) {
        const userPreference = this.getPreparedUserPreference();
        const { view, widget } = this.getPreparedViewAndWidget();

        await Promise.all([
          this.updateUserPreference({ data: userPreference }),
          this.updateView({ id: this.activeView._id, data: viewToRequest(view) }),
        ]);

        const oldQuery = this.getQueryById(widget._id);
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
