<template>
  <widget-settings-group :title="$tc('common.chart', 2)">
    <v-layout
      class="pa-3"
      column
    >
      <field-draggable-list
        v-field="charts"
        @edit="showEditChartModal"
        @remove="showRemoveChartModal"
      >
        <template #title="{ item }">
          <v-layout align-center>
            <v-icon large>
              {{ $constants.WIDGET_ICONS[item.type] }}
            </v-icon>
            <span class="ml-3">{{ item.title }}</span>
          </v-layout>
        </template>
      </field-draggable-list>
      <v-menu bottom>
        <template #activator="{ on }">
          <v-flex>
            <v-btn
              class="ml-0 mt-3"
              color="primary"
              v-on="on"
            >
              {{ $t('common.add') }}
            </v-btn>
          </v-flex>
        </template>
        <v-list>
          <v-list-item
            v-for="{ type, text, icon } in chartTypes"
            :key="text"
            @click="showCreateChartModal(type)"
          >
            <v-icon>{{ icon }}</v-icon>
            <span class="ml-3">{{ text }}</span>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-layout>
  </widget-settings-group>
</template>

<script>
import { MODALS, WIDGET_ICONS, WIDGET_TYPES } from '@/constants';

import { addKeyInEntity } from '@/helpers/array';

import { formArrayMixin } from '@/mixins/form';

import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import FieldDraggableList from '@/components/sidebars/form/fields/draggable-list.vue';

export default {
  components: { WidgetSettingsGroup, FieldDraggableList },
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

    showRemoveChartModal(chart, index) {
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
