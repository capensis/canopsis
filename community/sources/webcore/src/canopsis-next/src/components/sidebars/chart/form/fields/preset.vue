<template>
  <widget-settings-item :title="$t('settings.chart.preset')">
    <c-select
      v-model="preset"
      :items="presetsWithCustom"
      :label="$t('settings.chart.preset')"
      item-title="text"
      item-value="value"
    >
      <template #item="{ props, item }">
        <v-list-item v-bind="getItemProps(props)">
          <v-list-item-title>{{ getItem(item).text }}</v-list-item-title>

          <template v-if="getItem(item).helpText" #append>
            <c-help-icon
              :text="getItem(item).helpText"
              icon="help"
              size="20"
              left
            />
          </template>
        </v-list-item>
      </template>
    </c-select>
  </widget-settings-item>
</template>

<script>
import { isEqual, isUndefined } from 'lodash';

import { CHART_PRESET_CUSTOM_ITEM_VALUE, KPI_PIE_CHART_SHOW_MODS } from '@/constants';

import { formToMetricPresets } from '@/helpers/entities/metric/form';
import {
  getWidgetChartPresetParameters,
  getWidgetChartPresetTypesByWidgetType,
} from '@/helpers/entities/metric/widget';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    type: {
      type: String,
      required: true,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    preset: {
      get() {
        return this.presets.find(({ parameters }) => {
          const { metrics: oldMetrics, ...oldParameters } = this.modelValue;
          const { metrics: newMetrics, ...newParameters } = { ...this.modelValue, ...parameters };

          return isEqual(oldParameters, newParameters)
            && isEqual(formToMetricPresets(oldMetrics), formToMetricPresets(newMetrics));
        })?.value ?? CHART_PRESET_CUSTOM_ITEM_VALUE;
      },

      set(preset) {
        if (preset !== CHART_PRESET_CUSTOM_ITEM_VALUE) {
          this.$emit('update:modelValue', this.getParametersByPreset(preset));
        }
      },
    },

    availablePresetTypes() {
      return getWidgetChartPresetTypesByWidgetType(this.type);
    },

    presets() {
      return this.availablePresetTypes.map((value) => {
        const parameters = this.getParametersByPreset(value);
        const helpText = this.getPresetHelpTextByParameters(parameters);

        return {
          value,
          text: this.$t(`settings.chart.presets.${value}`),
          parameters,
          helpText,
        };
      });
    },

    presetsWithCustom() {
      return [
        {
          value: CHART_PRESET_CUSTOM_ITEM_VALUE,
          text: this.$t('common.custom'),
        },
        ...this.presets,
      ];
    },
  },

  methods: {
    getItem(item) {
      return item?.raw ?? item ?? {};
    },

    getItemProps(props) {
      const { title, ...itemProps } = props;

      return itemProps;
    },

    getParametersByPreset(preset) {
      return {
        ...this.modelValue,
        ...getWidgetChartPresetParameters(this.type, preset),
        chart_title: this.$t(`settings.chart.presetChartHeaders.${preset}`),
      };
    },

    getPresetHelpTextByParameters(parameters) {
      const result = [{
        label: this.$t('common.header'),
        value: parameters.chart_title,
      }];

      if (!isUndefined(parameters.stacked)) {
        result.push({
          label: this.$t('settings.chart.graphType'),
          value: parameters.stacked
            ? this.$t('settings.chart.stackedBars')
            : this.$t('settings.chart.separateBars'),
        });
      }

      if (parameters.aggregate_func) {
        result.push({
          label: this.$t('kpi.calculationMethod'),
          value: parameters.aggregate_func,
        });
      }

      if (!isUndefined(parameters.comparison)) {
        result.push({
          label: this.$t('settings.chart.showComparison'),
          value: parameters.comparison ? this.$t('common.enabled') : this.$t('common.disabled'),
        });
      }

      if (parameters.metrics) {
        result.push({
          label: this.$t('settings.chart.selectMetrics'),
          value: parameters.metrics.map(({ metric }) => `\n  - ${this.$t(`alarm.metrics.${metric}`)}`).join(''),
        });
      }

      if (parameters.default_time_range) {
        result.push({
          label: this.$t('settings.defaultTimeRange'),
          value: this.$t(`quickRanges.types.${parameters.default_time_range}`),
        });
      }

      if (parameters.sampling) {
        result.push({
          label: this.$t('settings.defaultSampling'),
          value: parameters.sampling,
        });
      }

      if (parameters.show_mode) {
        result.push({
          label: this.$t('settings.chart.sharesType'),
          value: this.$t(
            `common.${parameters.show_mode === KPI_PIE_CHART_SHOW_MODS.numbers ? 'number' : 'percent'}`,
            2,
          ),
        });
      }

      return result.map(({ label, value }) => `${label}: ${value}`).join('\n');
    },
  },
};
</script>
