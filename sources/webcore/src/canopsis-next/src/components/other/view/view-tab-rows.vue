<template lang="pug">
  div.view(:id="`view-tab-${tab._id}`")
    v-layout(v-for="row in rows", :key="row._id", row, wrap)
      v-flex(xs12)
        v-layout.hide-on-full-screen(align-center)
          h2.ml-1 {{ row.title }}
          v-tooltip.ml-2(v-if="isEditingMode", left)
            v-btn.ma-0(slot="activator", icon, @click.stop="showDeleteRowModal(row)")
              v-icon.error--text delete
            span {{ $t('common.delete') }}
      v-flex(
      v-for="widget in row.widgets",
      :key="widget._id",
      :class="getWidgetFlexClass(widget)"
      )
        v-layout.hide-on-full-screen(justify-space-between, align-center)
          v-flex
            h3.my-1.ml-2(v-show="widget.title") {{ widget.title }}
          v-flex(xs1, v-if="isEditingMode")
            v-btn.ma-0(v-if="hasUpdateAccess", icon, @click="showSettings(tab._id, row._id, widget)")
              v-icon settings
            v-tooltip(left)
              v-btn.ma-0(
              slot="activator",
              icon,
              @click="showDeleteWidgetModal(row._id, widget)"
              )
                v-icon.error--text delete
              span {{ $t('common.delete') }}
        component(
        :is="widgetsComponentsMap[widget.type]",
        :widget="widget",
        :tabId="tab._id",
        :isEditingMode="isEditingMode",
        )
</template>

<script>
import { MODALS, WIDGET_TYPES, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import AlarmsList from '@/components/other/alarm/alarms-list.vue';
import EntitiesList from '@/components/other/context/entities-list.vue';
import Weather from '@/components/other/service-weather/weather.vue';
import StatsHistogram from '@/components/other/stats/histogram/stats-histogram-wrapper.vue';
import StatsCurves from '@/components/other/stats/curves/stats-curves-wrapper.vue';
import StatsTable from '@/components/other/stats/stats-table.vue';
import StatsCalendar from '@/components/other/stats/stats-calendar.vue';
import StatsNumber from '@/components/other/stats/stats-number.vue';

import popupMixin from '@/mixins/popup';
import modalMixin from '@/mixins/modal';
import sideBarMixin from '@/mixins/side-bar/side-bar';

export default {
  components: {
    AlarmsList,
    EntitiesList,
    Weather,
    StatsHistogram,
    StatsCurves,
    StatsTable,
    StatsCalendar,
    StatsNumber,
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
        [WIDGET_TYPES.alarmList]: 'alarms-list',
        [WIDGET_TYPES.context]: 'entities-list',
        [WIDGET_TYPES.weather]: 'weather',
        [WIDGET_TYPES.statsHistogram]: 'stats-histogram',
        [WIDGET_TYPES.statsCurves]: 'stats-curves',
        [WIDGET_TYPES.statsTable]: 'stats-table',
        [WIDGET_TYPES.statsCalendar]: 'stats-calendar',
        [WIDGET_TYPES.statsNumber]: 'stats-number',
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
