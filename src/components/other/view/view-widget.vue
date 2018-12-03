<template lang="pug">
  div#view
    v-layout(v-for="(row, rowKey) in rows", :key="row._id", row, wrap)
      v-flex(xs12)
        v-layout.notDisplayedFullScreen(align-center)
          h2.ml-1 {{ row.title }}
          v-tooltip.ml-2(left, v-if="isEditingMode")
            v-btn.ma-0(slot="activator", icon, @click.stop="deleteRow(rowKey)")
              v-icon.error--text delete
            span {{ $t('common.delete') }}
      v-flex(
      v-for="(widget, widgetKey) in row.widgets",
      :key="`${widgetKeyPrefix}_${widget._id}`",
      :class="getWidgetFlexClass(widget)"
      )
        v-layout.notDisplayedFullScreen(justify-space-between, align-center)
          v-flex
            h3.my-1.ml-2(v-show="widget.title") {{ widget.title }}
          v-flex(xs1, v-if="isEditingMode")
            v-btn.ma-0(v-if="hasUpdateAccess", icon, @click="showSettings(row._id, tab._id, widget)")
              v-icon settings
            v-tooltip(left)
              v-btn.ma-0(
              slot="activator",
              icon,
              @click="deleteWidget(widgetKey, rowKey)"
              )
                v-icon.error--text delete
              span {{ $t('common.delete') }}
        component(
        :is="widgetsComponentsMap[widget.type]",
        :widget="widget",
        )
</template>

<script>
import get from 'lodash/get';
import pullAt from 'lodash/pullAt';

import { MODALS, WIDGET_TYPES, USERS_RIGHTS_MASKS, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';
import uid from '@/helpers/uid';

import AlarmsList from '@/components/other/alarm/alarms-list.vue';
import EntitiesList from '@/components/other/context/entities-list.vue';
import Weather from '@/components/other/service-weather/weather.vue';
import StatsHistogram from '@/components/other/stats/histogram/stats-histogram-wrapper.vue';
import StatsCurves from '@/components/other/stats/curves/stats-curves-wrapper.vue';
import StatsTable from '@/components/other/stats/stats-table.vue';
import StatsCalendar from '@/components/other/stats/stats-calendar.vue';
import StatsNumber from '@/components/other/stats/stats-number.vue';

import authMixin from '@/mixins/auth';
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
    authMixin,
    popupMixin,
    modalMixin,
    sideBarMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      widgetKeyPrefix: uid(),
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
      return get(this.tab, 'rows', []);
    },
    hasUpdateAccess() {
      return this.checkUpdateAccess(this.id, USERS_RIGHTS_MASKS.update);
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
    showSettings(rowId, tabId, widget) {
      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[widget.type],
        config: {
          widget,
          rowId,
          tabId,
        },
      });
    },

    deleteRow(rowKey) {
      if (this.view.rows[rowKey].widgets.length > 0) {
        this.addErrorPopup({ text: this.$t('errors.lineNotEmpty') });
      } else {
        this.showModal({
          name: MODALS.confirmation,
          config: {
            action: () => {
              const view = { ...this.view };
              pullAt(view.rows, rowKey);
              this.updateView({ id: this.id, data: view });
            },
          },
        });
      }
    },

    deleteWidget(widgetKey, rowKey) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => {
            const view = { ...this.view };
            pullAt(view.rows[rowKey].widgets, widgetKey);
            this.updateView({ id: this.id, data: view });
          },
        },
      });
    },
  },
};
</script>
