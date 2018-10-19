<template lang="pug">
  v-container
    div
      v-layout(v-for="(row, rowKey) in rows", :key="row._id", row, wrap)
        v-flex(xs12)
          v-layout(align-center)
            h2 {{ row.title }}
            v-tooltip.ml-2(left, v-if="isEditModeEnable")
              v-btn.ma-0(slot="activator", fab, small, dark, color="red darken-3", @click.stop="deleteRow(rowKey)")
                v-icon delete
              span {{ $t('common.delete') }}
        v-flex(
        v-for="(widget, widgetKey) in row.widgets",
        :key="`${widgetKeyPrefix}_${widget._id}`",
        :class="getWidgetFlexClass(widget)"
        )
          v-layout(justify-space-between)
            h3 {{ widget.title }}
            v-tooltip(left, v-if="isEditModeEnable")
              v-btn.ma-0(
              slot="activator",
              fab,
              small,
              dark,
              color="red darken-3",
              @click="deleteWidget(widgetKey, rowKey)"
              )
                v-icon delete
              span {{ $t('common.delete') }}
          component(
          :is="widgetsComponentsMap[widget.type]",
          :widget="widget",
          :rowId="row._id"
          )
    .fab
      v-speed-dial(
      v-model="fab",
      direction="top",
      transition="slide-y-reverse-transition"
      )
        v-btn(slot="activator", color="green darken-3", dark, fab, v-model="fab")
          v-icon menu
          v-icon close
        v-tooltip(bottom)
          v-btn(
          slot="activator",
          v-model="isFullScreenModeEnable"
          fab,
          dark,
          small,
          @click="fullScreenToggle",
          )
            v-icon fullscreen
            v-icon fullscreen_exit
          span alt + enter / command + enter
        v-tooltip(left)
          v-btn(slot="activator", fab, dark, small, color="info", @click.stop="refreshView")
            v-icon refresh
          span {{ $t('common.refresh') }}
        v-tooltip(left)
          v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateWidgetModal")
            v-icon add
          span {{ $t('common.addWidget') }}
        v-tooltip(left)
          v-btn(slot="activator", fab, dark, small, @click.stop="toggleViewEditMode", v-model="isEditModeEnable")
            v-icon edit
            v-icon done
          span {{ $t('common.toggleEditView') }}
</template>

<script>
import get from 'lodash/get';
import pullAt from 'lodash/pullAt';

import { WIDGET_TYPES, MODALS } from '@/constants';
import uid from '@/helpers/uid';

import AlarmsList from '@/components/other/alarm/alarms-list.vue';
import EntitiesList from '@/components/other/context/entities-list.vue';
import Weather from '@/components/other/service-weather/weather.vue';
import StatsHistogram from '@/components/other/stats/histogram/stats-histogram-wrapper.vue';
import StatsCurves from '@/components/other/stats/curves/stats-curves-wrapper.vue';
import StatsTable from '@/components/other/stats/stats-table.vue';
import StatsNumber from '@/components/other/stats/stats-number.vue';

import popupMixin from '@/mixins/popup';
import modalMixin from '@/mixins/modal/modal';
import entitiesViewMixin from '@/mixins/entities/view';

export default {
  components: {
    AlarmsList,
    EntitiesList,
    Weather,
    StatsHistogram,
    StatsCurves,
    StatsTable,
    StatsNumber,
  },
  mixins: [
    popupMixin,
    modalMixin,
    entitiesViewMixin,
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
        [WIDGET_TYPES.statsNumber]: 'stats-number',
      },
      widgetKeyPrefix: uid(),
      isEditModeEnable: false,
      isFullScreenModeEnable: false,
      fab: false,
    };
  },
  computed: {
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
    keyDownListener({ keyCode, altKey }) {
      if (keyCode === 13 && altKey) {
        this.fullScreenToggle();
      }
    },

    fullScreenToggle() {
      const element = document.getElementById('app');

      if (element) {
        this.$fullscreen.toggle(element, {
          wrap: false,
          fullscreenClass: '-fullscreen',
          callback: value => this.isFullScreenModeEnable = value,
        });
      }
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
    toggleViewEditMode() {
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

<style scoped>
  .fab {
    position: fixed;
    bottom: 0;
    right: 0;
  }
</style>
