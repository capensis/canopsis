<template>
  <widget-settings-item :title="$t('settings.chart.preset')">
    <v-select
      v-model="preset"
      :items="presetsWithCustom"
      :label="$t('settings.chart.preset')"
    >
      <template #item="{ item, attrs, on }">
        <v-list-item
          v-bind="attrs"
          v-on="on"
        >
          <v-list-item-content>{{ item.text }}</v-list-item-content>
          <v-list-item-action v-if="item.helpText">
            <c-help-icon
              :text="item.helpText"
              icon="help"
              size="20"
              left
            />
          </v-list-item-action>
        </v-list-item>
      </template>
    </v-select>
  </widget-settings-item>
</template>

<script>
import { isEqual, isUndefined } from 'lodash';

import { KPI_PIE_CHART_SHOW_MODS, CHART_PRESET_CUSTOM_ITEM_VALUE } from '@/constants';

import {
  getWidgetChartPresetParameters,
  getWidgetChartPresetTypesByWidgetType,
} from '@/helpers/entities/metric/widget';
import { formToMetricPresets } from '@/helpers/entities/metric/form';

import { formBaseMixin } from '@/mixins/form';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  mixins: [formBaseMixin],
  model: {
    prop: 'parameters',
    event: 'input',
  },
  props: {
    parameters: {
      type: Object,
      required: true,
    },
    type: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'preset',
    },
  },
  computed: {
    preset: {
      get() {
        return this.presets.find(({ parameters }) => {
          const { metrics: oldMetrics, ...oldParameters } = this.parameters;
          const { metrics: newMetrics, ...newParameters } = { ...this.parameters, ...parameters };

          return isEqual(oldParameters, newParameters)
            && isEqual(formToMetricPresets(oldMetrics), formToMetricPresets(newMetrics));
        }) ?? CHART_PRESET_CUSTOM_ITEM_VALUE;
      },

      set(preset) {
        if (preset !== CHART_PRESET_CUSTOM_ITEM_VALUE) {
          this.updateModel(this.getParametersByPreset(preset));
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
    getParametersByPreset(preset) {
      return {
        ...this.parameters,
        ...getWidgetChartPresetParameters(this.type, preset),
        chart_title: this.$t(`settings.chart.presetChartHeaders.${preset}`),
      };
    },

    getPresetHelpTextByParameters(parameters) {
      const result = [{
        label: this.$tc('common.header'),
        value: parameters.chart_title,
      }];

      if (!isUndefined(parameters.stacked)) {
        result.push({
          label: this.$t('settings.chart.graphType'),
          value: parameters.stacked ? this.$t('settings.chart.stackedBars') : this.$t('settings.chart.separateBars'),
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
          value: this.$tc(
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
