<template>
  <v-menu>
    <template #activator="{ props }">
      <v-btn
        size="small"
        variant="text"
        icon
        v-bind="props"
      >
        <v-icon size="small">
          more_horiz
        </v-icon>
      </v-btn>
    </template>

    <v-list density="compact">
      <v-list-item v-if="hasSettingsSidebar" @click="showSettings">
        <v-list-item-title>
          {{ $t('common.edit') }}
        </v-list-item-title>
      </v-list-item>

      <v-list-item v-if="hasSettingsSidebar" @click="showSelectViewTabModal">
        <v-list-item-title>
          {{ $t('common.duplicate') }}
        </v-list-item-title>
      </v-list-item>

      <v-list-item @click="copyWidgetId">
        <v-list-item-title>
          {{ $t('view.copyWidgetId') }}
        </v-list-item-title>
      </v-list-item>

      <v-list-item @click="showDeleteWidgetModal">
        <v-list-item-title class="text-error">
          {{ $t('common.delete') }}
        </v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { mapActions } from 'pinia';

import { MODALS, SIDE_BARS_BY_WIDGET_TYPES, WIDGET_TYPES } from '@/constants';

import { useViewStore } from '@/stores/modules/view';
import { useWidgetStore } from '@/stores/modules/widget';

import { writeTextToClipboard } from '@/helpers/clipboard';
import { calculateNewWidgetGridParametersY } from '@/helpers/entities/widget/grid';
import { setSeveralFields } from '@/helpers/immutable';

import { activeViewMixin } from '@/mixins/active-view';
import { viewRouterMixin } from '@/mixins/view/router';

const MIGRATED_SETTINGS_WIDGET_TYPES = [
  WIDGET_TYPES.alarmList,
  WIDGET_TYPES.context,
  WIDGET_TYPES.text,
  WIDGET_TYPES.counter,
  WIDGET_TYPES.map,
  WIDGET_TYPES.externalDataTable,
  WIDGET_TYPES.serviceWeather,
];

export default {
  mixins: [activeViewMixin, viewRouterMixin],
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
  computed: {
    hasSettingsSidebar() {
      return MIGRATED_SETTINGS_WIDGET_TYPES.includes(this.widget.type);
    },
  },
  methods: {
    ...mapActions(useWidgetStore, {
      removeWidget: 'remove',
    }),
    ...mapActions(useViewStore, {
      fetchViewTab: 'fetchViewTab',
    }),

    showSettings() {
      this.$sidebar.show({
        name: SIDE_BARS_BY_WIDGET_TYPES[this.widget.type],
        config: {
          widget: {
            ...this.widget,
            tab: this.tab._id,
          },
        },
      });
    },

    /**
     * Copy a widget parameters and open corresponding settings sidebar.
     *
     * @param {Object} payload
     * @param {string} payload.viewId
     * @param {string} payload.tabId
     * @returns {Promise<void>}
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
     * Show select view tab modal window.
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
     * Copy widget id to clipboard.
     */
    async copyWidgetId() {
      try {
        await writeTextToClipboard(this.widget._id);

        this.$popups.success({ text: this.$t('success.widgetIdCopied') });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
