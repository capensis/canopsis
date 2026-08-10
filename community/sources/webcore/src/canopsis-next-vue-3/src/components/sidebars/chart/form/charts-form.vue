<template>
  <widget-settings-item :title="$t('common.chart', 2)">
    <c-layout gap="2" column>
      <c-draggable-list-field
        :model-value="modelValue"
        handle=".chart-drag-handler"
        @update:model-value="$emit('update:modelValue', $event)"
      >
        <template #item="{ element, index }">
          <c-layout align="center" gap="2">
            <v-icon class="chart-drag-handler" icon="drag_indicator" />
            <v-icon :icon="widgetIcons[element.type]" />
            <span class="flex-grow-1">{{ element.title }}</span>
            <v-btn
              icon="edit"
              variant="text"
              size="small"
              @click="editChart(element, index)"
            />
            <v-btn
              icon="delete"
              variant="text"
              size="small"
              @click="removeChart(index)"
            />
          </c-layout>
        </template>
      </c-draggable-list-field>
      <div class="align-self-start">
        <v-menu>
          <template #activator="{ props }">
            <v-btn v-bind="props" color="primary" size="small">
              {{ $t('common.add') }}
            </v-btn>
          </template>
          <v-list>
            <v-list-item v-for="type in chartTypes" :key="type" @click="createChart(type)">
              <template #prepend>
                <v-icon :icon="widgetIcons[type]" />
              </template>
              <v-list-item-title>{{
                $t(`modals.createWidget.types.${type}.title`)
              }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </div>
    </c-layout>
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import { MODALS, WIDGET_ICONS, WIDGET_TYPES } from '@/constants';
import { addKeyInEntity } from '@/helpers/array';

export default {
  components: { WidgetSettingsItem },
  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
  },
  emits: ['update:modelValue'],
  data() {
    return {
      widgetIcons: WIDGET_ICONS,
      chartTypes: [WIDGET_TYPES.barChart, WIDGET_TYPES.lineChart, WIDGET_TYPES.numbers],
    };
  },
  methods: {
    createChart(type) {
      this.$modals.show({
        name: MODALS.createAlarmChart,
        config: {
          chart: { type },
          title: this.$t(`modals.createAlarmChart.${type}.create.title`),
          onlyExternal: true,
          action: chart =>
            this.$emit('update:modelValue', [...this.modelValue, addKeyInEntity(chart)]),
        },
      });
    },
    editChart(chart, index) {
      this.$modals.show({
        name: MODALS.createAlarmChart,
        config: {
          chart,
          title: this.$t(`modals.createAlarmChart.${chart.type}.edit.title`),
          onlyExternal: true,
          action: (value) => {
            const charts = [...this.modelValue];
            charts[index] = { ...value, key: chart.key };
            this.$emit('update:modelValue', charts);
          },
        },
      });
    },
    removeChart(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () =>
            this.$emit(
              'update:modelValue',
              this.modelValue.filter((_, i) => i !== index),
            ),
        },
      });
    },
  },
};
</script>
