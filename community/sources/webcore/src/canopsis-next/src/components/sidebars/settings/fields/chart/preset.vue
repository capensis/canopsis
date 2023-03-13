<template lang="pug">
  widget-settings-item(:title="$t('settings.chart.preset')")
    v-select(
      v-model="preset",
      :items="presets",
      :label="$t('settings.chart.preset')"
    )
</template>

<script>
import { isEqual } from 'lodash';

import { getWidgetChartPresetParameters, getWidgetChartPresetTypesByWidgetType } from '@/helpers/entities/widget';

import { formBaseMixin } from '@/mixins/form';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

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
        return this.availablePresetTypes.find(preset => isEqual(this.parameters, this.getParametersByPreset(preset)));
      },

      set(preset) {
        this.updateModel(this.getParametersByPreset(preset));
      },
    },

    availablePresetTypes() {
      return getWidgetChartPresetTypesByWidgetType(this.type);
    },

    presets() {
      return this.availablePresetTypes.map(value => ({
        value,
        text: this.$t(`settings.chart.presets.${value}`),
      }));
    },
  },
  methods: {
    getParametersByPreset(preset) {
      return {
        ...this.parameters,
        ...getWidgetChartPresetParameters(this.type, preset),
        chart_header: this.$t(`settings.chart.presetChartHeaders.${preset}`),
      };
    },
  },
};
</script>
