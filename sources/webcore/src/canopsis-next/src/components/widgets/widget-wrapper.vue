<template lang="pug">
  v-card.full-height
    v-card-title.lighten-1.pa-0
      v-layout(justify-space-between, align-center)
        v-flex
          h4.ml-2.font-weight-regular {{ widget.title }}
        v-spacer
        v-layout(justify-end)
          v-menu(offset-y)
            v-btn.ma-0(icon, small, slot="activator")
              v-icon(small) more_horiz
            v-list(dense)
              v-list-tile(
                @click="showSettings({ tabId: tab._id, rowId: row._id, widget })",
                :data-test="`editWidgetButton-${widget._id}`"
              )
                div {{ $t('common.edit') }}
              v-list-tile(@click="showSelectViewTabModal", :data-test="`copyWidgetButton-${widget._id}`")
                div {{ $t('common.duplicate') }}
              v-list-tile(@click="showDeleteWidgetModal", :data-test="`deleteWidgetButton-${widget._id}`")
                v-list-tile-title.error--text {{ $t('common.delete') }}
    v-container.pa-0(fill-height, fluid)
      v-card-text.pa-0
        component(
          :is="widgetsComponentsMap[widget.type]",
          :widget="widget",
          :tabId="tab._id",
          :isEditingMode="isEditingMode"
        )
</template>

<script>
import { cloneDeep } from 'lodash';
import { GridItem } from 'vue-grid-layout';

import { WIDGET_TYPES, MODALS, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import modalMixin from '@/mixins/modal';
import sideBarMixin from '@/mixins/side-bar/side-bar';

import { generateWidgetByType } from '@/helpers/entities';

import AlarmsListWidget from '@/components/other/alarm/alarms-list.vue';
import EntitiesListWidget from '@/components/other/context/entities-list.vue';
import WeatherWidget from '@/components/other/service-weather/weather.vue';
import StatsHistogramWidget from '@/components/other/stats/histogram/stats-histogram.vue';
import StatsCurvesWidget from '@/components/other/stats/curves/stats-curves.vue';
import StatsTableWidget from '@/components/other/stats/stats-table.vue';
import StatsCalendarWidget from '@/components/other/stats/calendar/stats-calendar.vue';
import StatsNumberWidget from '@/components/other/stats/stats-number.vue';
import StatsParetoWidget from '@/components/other/stats/pareto/stats-pareto.vue';
import TextWidget from '@/components/other/text/text.vue';

export default {
  components: {
    AlarmsListWidget,
    EntitiesListWidget,
    WeatherWidget,
    StatsHistogramWidget,
    StatsCurvesWidget,
    StatsTableWidget,
    StatsCalendarWidget,
    StatsNumberWidget,
    StatsParetoWidget,
    TextWidget,
    GridItem,
  },
  mixins: [modalMixin, sideBarMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    updateTabMethod: {
      type: Function,
      required: true,
    },
  },
  data() {
    return {
      widgetsComponentsMap: {
        [WIDGET_TYPES.alarmList]: 'alarms-list-widget',
        [WIDGET_TYPES.context]: 'entities-list-widget',
        [WIDGET_TYPES.weather]: 'weather-widget',
        [WIDGET_TYPES.statsHistogram]: 'stats-histogram-widget',
        [WIDGET_TYPES.statsCurves]: 'stats-curves-widget',
        [WIDGET_TYPES.statsTable]: 'stats-table-widget',
        [WIDGET_TYPES.statsCalendar]: 'stats-calendar-widget',
        [WIDGET_TYPES.statsNumber]: 'stats-number-widget',
        [WIDGET_TYPES.statsPareto]: 'stats-pareto-widget',
        [WIDGET_TYPES.text]: 'text-widget',
      },
    };
  },
  methods: {
    deleteWidgetFromTab(widgetId) {
      const newTab = cloneDeep(this.tab);

      const widgetIndex = newTab.widgets.findIndex(widget => widget._id === widgetId);

      newTab.widgets.splice(widgetIndex, 1);

      newTab.layout = this.deleteWidgetFromLayout(widgetId);

      return newTab;
    },

    deleteWidgetFromLayout(widgetId) {
      const newLayout = cloneDeep(this.tab.layout);

      const widgetIndex = newLayout.findIndex(widget => widget._id === widgetId);

      newLayout.splice(widgetIndex, 1);

      return newLayout;
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
      this.showModal({
        name: MODALS.selectViewTab,
        config: {
          action: ({ tabId, viewId }) => this.cloneWidget({ tabId, viewId }),
        },
      });
    },

    showDeleteWidgetModal() {
      this.showModal({
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

<style lang="scss" scoped>
  .full-height {
    height: 100%;
    position: relative;
  }
</style>
