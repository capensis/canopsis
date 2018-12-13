<template lang="pug">
  div
    div#view
      v-layout(v-for="(row, rowKey) in rows", :key="row._id", row, wrap)
        v-flex(xs12)
          v-layout.notDisplayedFullScreen(align-center)
            h2.ml-1 {{ row.title }}
            v-tooltip.ml-2(v-if="isEditModeEnable", left)
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
            v-flex(xs1, v-if="isEditModeEnable")
              v-btn.ma-0(v-if="hasUpdateAccess", icon, @click="showSettings(row._id, widget)")
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
          :rowId="row._id",
          :hasUpdateAccess="hasUpdateAccess"
          )
    .fab
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, color="secondary", @click.stop="refreshView")
          v-icon refresh
        span {{ $t('common.refresh') }}
      v-speed-dial(
      v-model="fab",
      direction="left",
      transition="slide-y-reverse-transition"
      )
        v-btn(slot="activator", color="primary", dark, fab, v-model="fab")
          v-icon menu
          v-icon close
        v-tooltip(top)
          v-btn(
          slot="activator",
          v-model="isFullScreenModeEnable"
          fab,
          dark,
          small,
          @click="fullScreenModeToggle",
          )
            v-icon fullscreen
            v-icon fullscreen_exit
          span alt + enter / command + enter
        v-tooltip(top)
          v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateWidgetModal")
            v-icon add
          span {{ $t('common.addWidget') }}
        v-tooltip(v-if="hasUpdateAccess", top)
          v-btn(slot="activator", fab, dark, small, @click.stop="viewEditModeToggle", v-model="isEditModeEnable")
            v-icon edit
            v-icon done
          span {{ $t('common.toggleEditView') }} (ctrl + e / command + e)
</template>

<script>
import get from 'lodash/get';
import pullAt from 'lodash/pullAt';

import { WIDGET_TYPES, MODALS, USERS_RIGHTS_MASKS, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';
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
import modalMixin from '@/mixins/modal/modal';
import entitiesViewMixin from '@/mixins/entities/view';
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
    entitiesViewMixin,
    sideBarMixin,
  ],
  props: {
    id: {
      type: [String, Number],
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
      widgetKeyPrefix: uid(),
      isEditModeEnable: false,
      isFullScreenModeEnable: false,
      fab: false,
    };
  },
  computed: {
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
    rows() {
      return get(this.view, 'rows', []);
    },
  },
  created() {
    document.addEventListener('keydown', this.keyDownListener);
    this.fetchView({ id: this.id });
  },
  beforeDestroy() {
    this.$fullscreen.exit();
    document.removeEventListener('keydown', this.keyDownListener);
  },
  methods: {
    keyDownListener(event) {
      if (event.keyCode === 13 && event.altKey) {
        this.fullScreenModeToggle();
        event.preventDefault();
      } else if (event.keyCode === 69 && event.ctrlKey) {
        this.viewEditModeToggle();
        event.preventDefault();
      }
    },

    fullScreenModeToggle() {
      const element = document.getElementById('view');

      if (element) {
        this.$fullscreen.toggle(element, {
          fullscreenClass: 'fullscreen',
          background: 'white',
          callback: value => this.isFullScreenModeEnable = value,
        });
      }
    },

    showSettings(rowId, widget) {
      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[widget.type],
        config: {
          widget,
          rowId,
        },
      });
    },

    async refreshView() {
      await this.fetchView({ id: this.id });

      this.widgetKeyPrefix = uid();
    },

    showCreateWidgetModal() {
      this.showModal({
        name: MODALS.createWidget,
      });
    },

    viewEditModeToggle() {
      this.isEditModeEnable = !this.isEditModeEnable;
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

<style lang="scss" scoped>
.fullscreen {
  .notDisplayedFullScreen {
    display: none;
  }
}
</style>
