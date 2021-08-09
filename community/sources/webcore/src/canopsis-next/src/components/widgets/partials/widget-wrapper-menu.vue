<template lang="pug">
  v-menu(offset-y)
    v-btn.ma-0(slot="activator", icon, small)
      v-icon(small) more_horiz
    v-list(dense)
      v-list-tile(@click="showSettings({ tabId: tab._id, widget })")
        div {{ $t('common.edit') }}
      v-list-tile(@click="showSelectViewTabModal")
        div {{ $t('common.duplicate') }}
      v-list-tile(
        v-clipboard:copy="widget._id",
        v-clipboard:success="addWidgetIdCopiedSuccessPopup",
        v-clipboard:error="addWidgetIdCopiedErrorPopup",
        @click=""
      )
        div {{ $t('view.copyWidgetId') }}
      v-list-tile(@click="showDeleteWidgetModal")
        v-list-tile-title.error--text {{ $t('common.delete') }}
</template>

<script>
import { MODALS, ROUTE_NAMES, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import sideBarMixin from '@/mixins/side-bar/side-bar';

import { generateWidgetId } from '@/helpers/entities';
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
    /**
     * Delete widget from tab
     *
     * @param {string} widgetId
     */
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
            name: ROUTE_NAMES.view,
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
      const newWidget = { ...this.widget, _id: generateWidgetId(this.widget.type) };

      this.redirectToSelectedViewAndTab({ tabId, viewId });

      this.showSettings({ viewId, tabId, widget: newWidget });
    },

    /**
     * Show widget settings side bar
     *
     * @param {string} [viewId]
     * @param {string} tabId
     * @param {Object} widget
     */
    showSettings({
      viewId,
      tabId,
      widget,
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

    /**
     * Show select view tab modal window
     */
    showSelectViewTabModal() {
      this.$modals.show({
        name: MODALS.selectViewTab,
        config: {
          action: ({ tabId, viewId }) => this.cloneWidget({ tabId, viewId }),
        },
      });
    },

    /**
     * Show delete widget modal window
     */
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

    /**
     * Add success popup for widgetId copying
     */
    addWidgetIdCopiedSuccessPopup() {
      this.$popups.success({ text: this.$t('success.widgetIdCopied') });
    },

    /**
     * Add error popup for widgetId copying
     */
    addWidgetIdCopiedErrorPopup() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>
