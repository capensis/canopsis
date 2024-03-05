<template>
  <v-menu offset-y>
    <template #activator="{ on }">
      <v-btn
        icon
        small
        v-on="on"
      >
        <v-icon small>
          more_horiz
        </v-icon>
      </v-btn>
    </template>
    <v-list dense>
      <v-list-item @click="showSettings">
        <div>{{ $t('common.edit') }}</div>
      </v-list-item>
      <v-list-item @click="showSelectViewTabModal">
        <div>{{ $t('common.duplicate') }}</div>
      </v-list-item>
      <v-list-item
        v-clipboard:copy="widget._id"
        v-clipboard:success="addWidgetIdCopiedSuccessPopup"
        v-clipboard:error="addWidgetIdCopiedErrorPopup"
        @click.stop=""
      >
        <div>{{ $t('view.copyWidgetId') }}</div>
      </v-list-item>
      <v-list-item @click="showDeleteWidgetModal">
        <v-list-item-title class="error--text">
          {{ $t('common.delete') }}
        </v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { MODALS, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import { calculateNewWidgetGridParametersY } from '@/helpers/entities/widget/grid';
import { setSeveralFields } from '@/helpers/immutable';

import { activeViewMixin } from '@/mixins/active-view';
import { viewRouterMixin } from '@/mixins/view/router';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';
import { entitiesViewTabMixin } from '@/mixins/entities/view/tab';

export default {
  mixins: [
    activeViewMixin,
    viewRouterMixin,
    entitiesWidgetMixin,
    entitiesViewTabMixin,
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
     * Copy a widget's parameters, and open corresponding settings panel
     */
    async cloneWidget({ viewId, tabId }) {
      const tab = await this.fetchViewTab({ id: tabId });

      const { mobile, tablet, desktop } = calculateNewWidgetGridParametersY(tab.widgets);
      const newWidget = setSeveralFields(this.widget, {
        tab: tabId,
        'grid_parameters.mobile.y': mobile,
        'grid_parameters.tablet.y': tablet,
        'grid_parameters.desktop.y': desktop,
      });

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
