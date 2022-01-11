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
import { omit } from 'lodash';

import { MODALS, ROUTES_NAMES, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';
import { entitiesViewTabMixin } from '@/mixins/entities/view/tab';

export default {
  mixins: [entitiesWidgetMixin, entitiesViewTabMixin],
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
    /**
     * Redirect to selected view and tab, if it's different then the view/tab we're actually on
     */
    redirectToSelectedViewAndTab({ tabId, viewId }) {
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
    cloneWidget({ viewId, tabId }) {
      const newWidget = omit(this.widget, ['_id']);

      this.redirectToSelectedViewAndTab({ tabId, viewId });

      this.$sidebar.show({
        name: SIDE_BARS_BY_WIDGET_TYPES[newWidget.type],
        config: {
          tabId,
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

            return this.fetchViewTab({ id: this.tab._id });
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
