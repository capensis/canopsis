<template lang="pug">
  v-menu(offset-y)
    template(#activator="{ on }")
      v-btn.ma-0(v-on="on", icon, small)
        v-icon(small) more_horiz
    v-list(dense)
      v-list-tile(@click="showSettings")
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
import { MODALS, ROUTES_NAMES, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import { activeViewMixin } from '@/mixins/active-view';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

export default {
  mixins: [
    activeViewMixin,
    entitiesWidgetMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
  },
  methods: {
    showSettings() {
      this.$sidebar.show({
        name: SIDE_BARS_BY_WIDGET_TYPES[this.widget.type],
        config: {
          widget: { ...this.widget, tab: this.tab._id },
        },
      });
    },

    /**
     * Redirect to selected view and tab, if it's different then the view/tab we're actually on
     */
    redirectToSelectedViewAndTab({ viewId, tabId }) {
      return new Promise((resolve, reject) => {
        if (this.tab._id === tabId) {
          return resolve();
        }

        return this.$router.push({
          name: ROUTES_NAMES.view,
          params: { id: viewId },
          query: { tabId },
        }, resolve, reject);
      });
    },

    /**
     * Copy a widget's parameters, and open corresponding settings panel
     */
    async cloneWidget({ viewId, tabId }) {
      const newWidget = {
        ...this.widget,

        tab: tabId,
      };

      await this.redirectToSelectedViewAndTab({ viewId, tabId });

      this.$sidebar.show({
        name: SIDE_BARS_BY_WIDGET_TYPES[newWidget.type],
        config: {
          duplicate: true,
          widget: newWidget,
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
          action: async () => {
            await this.removeWidget({ id: this.widget._id });

            return this.fetchActiveView();
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
