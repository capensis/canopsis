<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    v-layout(v-for="row in rows", :key="row._id" wrap)
      v-flex(xs12)
        v-layout.hide-on-full-screen(justify-end)
          v-btn.ma-2(
          v-if="isEditingMode && hasUpdateAccess",
          @click.stop="showDeleteRowModal(row)",
          small,
          color="error",
          ) {{ $t('view.deleteRow') }} - {{ row.title }}
      v-flex(
      v-for="widget in row.widgets",
      :key="widget._id",
      :class="getWidgetFlexClass(widget)"
      )
        v-layout.hide-on-full-screen(align-center, justify-space-between)
          h3.my-1.mx-2(v-show="widget.title") {{ widget.title }}
          v-layout(justify-end)
            template(v-if="isEditingMode && hasUpdateAccess")
              v-btn.ma-1(
              @click="showDeleteWidgetModal(row._id, widget)",
              small,
              color="error",
              ) {{ $t('view.deleteWidget') }}
              v-btn.ma-1(
              @click="showSettings(tab._id, row._id, widget)",
              icon
              )
                v-icon settings
        component(
        :is="widgetsComponentsMap[widget.type]",
        :widget="widget",
        :tabId="tab._id",
        :isEditingMode="isEditingMode",
        )
</template>

<script>
import { MODALS, WIDGET_TYPES, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import AlarmsListWidget from '@/components/other/alarm/alarms-list.vue';
import EntitiesListWidget from '@/components/other/context/entities-list.vue';
import WeatherWidget from '@/components/other/service-weather/weather.vue';
import StatsHistogramWidget from '@/components/other/stats/histogram/stats-histogram.vue';
import StatsCurvesWidget from '@/components/other/stats/curves/stats-curves.vue';
import StatsTableWidget from '@/components/other/stats/stats-table.vue';
import StatsCalendarWidget from '@/components/other/stats/calendar/stats-calendar.vue';
import StatsNumberWidget from '@/components/other/stats/stats-number.vue';
import TextWidget from '@/components/other/text/text.vue';

import popupMixin from '@/mixins/popup';
import modalMixin from '@/mixins/modal';
import sideBarMixin from '@/mixins/side-bar/side-bar';

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
    TextWidget,
  },
  mixins: [
    popupMixin,
    modalMixin,
    sideBarMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
    hasUpdateAccess: {
      type: Boolean,
      default: false,
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
        [WIDGET_TYPES.text]: 'text-widget',
      },
    };
  },
  computed: {
    rows() {
      return this.tab.rows || [];
    },

    getWidgetFlexClass() {
      return widget => [
        `xs${widget.size.sm}`,
        `md${widget.size.md}`,
        `lg${widget.size.lg}`,
      ];
    },
  },
  methods: {
    showSettings(tabId, rowId, widget) {
      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[widget.type],
        config: {
          tabId,
          rowId,
          widget,
        },
      });
    },

    showDeleteRowModal(row = {}) {
      const widgets = row.widgets || [];

      if (widgets.length > 0) {
        this.addErrorPopup({ text: this.$t('errors.lineNotEmpty') });
      } else {
        this.showModal({
          name: MODALS.confirmation,
          config: {
            action: () => {
              const newTab = {
                ...this.tab,

                rows: this.rows.filter(tabRow => tabRow._id !== row._id),
              };

              return this.updateTabMethod(newTab);
            },
          },
        });
      }
    },

    showDeleteWidgetModal(rowId, widget = {}) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => {
            const newTab = {
              ...this.tab,

              rows: this.rows.map((tabRow) => {
                if (tabRow._id === rowId) {
                  return {
                    ...tabRow,
                    widgets: tabRow.widgets.filter(rowWidget => rowWidget._id !== widget._id),
                  };
                }

                return tabRow;
              }),
            };

            return this.updateTabMethod(newTab);
          },
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .full-screen {
    .hide-on-full-screen {
      display: none;
    }
  }
</style>
