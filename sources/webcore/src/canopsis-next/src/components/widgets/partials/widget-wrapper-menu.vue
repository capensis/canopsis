<template lang="pug">
  v-menu(offset-y)
    v-btn.ma-0(icon, small, slot="activator")
      v-icon(small) more_horiz
    v-list(dense)
      v-list-tile(@click="showSettings({ tabId: tab._id, widget })")
        div {{ $t('common.edit') }}
      v-list-tile(@click="showSelectViewTabModal")
        div {{ $t('common.duplicate') }}
      v-list-tile(@click="showDeleteWidgetModal")
        v-list-tile-title.error--text {{ $t('common.delete') }}
</template>

<script>
import { MODALS, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import sideBarMixin from '@/mixins/side-bar/side-bar';

import { generateWidgetByType } from '@/helpers/entities';
import { removeFrom } from '@/helpers/immutable';

export default {
  mixins: [sideBarMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
    updateTabMethod: {
      type: Function,
      required: true,
    },
  },
  methods: {
    deleteWidgetFromTab(widgetId) {
      const widgetIndex = this.tab.widgets.findIndex(widget => widget._id === widgetId);

      return removeFrom(this.tab, 'widgets', widgetIndex);
    },

    /**
     * Redirect to selected view and tab, if it's different then the view/tab we're actually on
     */
    async redirectToSelectedViewAndTab({ tabId, viewId }) {
      await new Promise((resolve, reject) => {
        if (this.tab._id === tabId) {
          resolve();
        } else {
          this.$router.push({
            name: 'view',
            params: { id: viewId },
            query: { tabId },
          }, resolve, reject);
        }
      });
    },

    /**
     * Copy a widget's parameters, and open corresponding settings panel
     */
    cloneWidget({ viewId, tabId }) {
      const { _id: newWidgetId } = generateWidgetByType(this.widget.type);

      // Copy widget parameters,
      const newWidget = { ...this.widget, _id: newWidgetId };

      this.redirectToSelectedViewAndTab({ tabId, viewId });

      this.showSettings({ viewId, tabId, widget: newWidget });
    },

    showSettings({
      viewId,
      widget,
      tabId,
    }) {
      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[widget.type],
        config: {
          viewId,
          tabId,
          widget,
        },
      });
    },

    showSelectViewTabModal() {
      this.$modals.show({
        name: MODALS.selectViewTab,
        config: {
          action: ({ tabId, viewId }) => this.cloneWidget({ tabId, viewId }),
        },
      });
    },

    showDeleteWidgetModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            const updatedTab = this.deleteWidgetFromTab(this.widget._id);

            return this.updateTabMethod(updatedTab);
          },
        },
      });
    },
  },
};
</script>
