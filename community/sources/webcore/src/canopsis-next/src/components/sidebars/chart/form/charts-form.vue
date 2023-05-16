<template lang="pug">
  widget-settings-group(:title="$tc('common.chart', 2)")
    v-layout.pa-3(column)
      c-draggable-list-field(v-field="charts", :handle="`.${chartHandleClass}`")
        v-layout(v-for="(chart, index) in charts", :key="chart.key", row, align-center)
          v-flex(xs1)
            v-icon.draggable(:class="chartHandleClass") drag_indicator
          v-flex(xs8)
            v-layout(row, align-center)
              v-icon(large) {{ $constants.WIDGET_ICONS[chart.type] }}
              span.ml-3 {{ chart.title }}
          v-flex(xs3)
            c-action-btn(type="edit", @click="showEditChartModal(chart, index)")
            c-action-btn(type="delete", @click="showRemoveChartModal(index)")
      v-menu(bottom)
        template(#activator="{ on }")
          v-flex
            v-btn.ml-0.mt-3(v-on="on", color="primary") {{ $t('common.add') }}
        v-list
          v-list-tile(
            v-for="{ type, text, icon } in chartTypes",
            :key="text",
            @click="showCreateChartModal(type)"
          )
            v-icon {{ icon }}
            span.ml-3 {{ text }}
</template>

<script>
import { MODALS, WIDGET_ICONS, WIDGET_TYPES } from '@/constants';

import { addKeyInEntity } from '@/helpers/array';

import { formArrayMixin } from '@/mixins/form';

import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';

export default {
  components: { WidgetSettingsGroup },
  mixins: [formArrayMixin],
  model: {
    prop: 'charts',
    event: 'input',
  },
  props: {
    charts: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    chartHandleClass() {
      return 'chart-drag-handler';
    },

    chartTypes() {
      return [
        WIDGET_TYPES.barChart,
        WIDGET_TYPES.lineChart,
        WIDGET_TYPES.numbers,
      ].map(type => ({
        text: this.$t(`modals.createWidget.types.${type}.title`),
        icon: WIDGET_ICONS[type],
        type,
      }));
    },
  },
  methods: {
    showCreateChartModal(type) {
      this.$modals.show({
        name: MODALS.createAlarmChart,
        config: {
          chart: { type },
          title: this.$t(`modals.createAlarmChart.${type}.create.title`),
          onlyExternal: true,
          action: newChart => this.addItemIntoArray(addKeyInEntity(newChart)),
        },
      });
    },

    showEditChartModal(chart, index) {
      this.$modals.show({
        name: MODALS.createAlarmChart,
        config: {
          chart,
          title: this.$t(`modals.createAlarmChart.${chart.type}.edit.title`),
          onlyExternal: true,
          action: newChart => this.updateItemInArray(index, { ...newChart, key: chart.key }),
        },
      });
    },

    showRemoveChartModal(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeItemFromArray(index),
        },
      });
    },
  },
};
</script>
