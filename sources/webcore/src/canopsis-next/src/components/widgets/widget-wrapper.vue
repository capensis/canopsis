<template lang="pug">
  v-card
    template(v-if="widget.title || isEditingMode")
      v-card-title.lighten-1.pa-1
        v-layout(justify-space-between, align-center)
          v-flex
            h4.ml-2.font-weight-regular {{ widget.title }}
          v-spacer
          v-layout(justify-end, v-if="isEditingMode")
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
      v-divider
    v-card-text.pa-0
      component(
        v-bind="widgetsComponentsMap(widget.type).bind",
        :widget="widget",
        :tabId="tab._id",
        :isEditingMode="isEditingMode"
      )
</template>

<script>
import { cloneDeep } from 'lodash';

import { WIDGET_TYPES, WIDGET_TYPES_RULES, MODALS, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

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
import AlertOverlay from '@/components/layout/alert/alert-overlay.vue';

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
    AlertOverlay,
  },
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
    row: {
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
  computed: {
    widgetsComponentsMap() {
      return (widgetType) => {
        const baseMap = {
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
        };

        let widgetSpecificsProp = {};

        Object.entries(WIDGET_TYPES_RULES).forEach(([key, rule]) => {
          if (rule.edition !== this.edition) {
            baseMap[key] = 'alert-overlay';
            widgetSpecificsProp = {
              message: this.$t('errors.statsWrongEditionError'),
              value: true,
            };
          }
        });

        return {
          bind: {
            ...widgetSpecificsProp,
            is: baseMap[widgetType],
          },
        };
      };
    },
  },
  methods: {
    deleteWidgetFromTabRow(widgetId) {
      const newTab = cloneDeep(this.tab);

      const rowIndex = this.tab.rows.findIndex(row => row.widgets.find(widget => widget._id === widgetId));

      const widgetIndex = this.tab.rows[rowIndex].widgets.findIndex(widget => widget._id === widgetId);

      newTab.rows[rowIndex].widgets.splice(widgetIndex, 1);

      return newTab;
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
      rowId,
    }) {
      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[widget.type],
        config: {
          viewId,
          tabId,
          rowId,
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
            const updatedTab = this.deleteWidgetFromTabRow(this.widget._id);

            return this.updateTabMethod(updatedTab);
          },
        },
      });
    },
  },
};
</script>
