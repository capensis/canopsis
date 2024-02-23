<template>
  <widget-settings-item :title="$t('settings.colorsSelector.title')">
    <v-layout column>
      <c-color-picker-field
        v-field="levelsColors[level]"
        v-for="level in $constants.STATS_CRITICITY"
        :key="level"
        :label="getButtonText(level)"
        :type="colorType"
        class="mb-3"
        splitted
      />
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { ALARM_LEVELS_COLORS, STATS_CRITICITY } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  model: {
    prop: 'levelsColors',
    event: 'input',
  },
  props: {
    levelsColors: {
      type: Object,
      default: () => ({ ...ALARM_LEVELS_COLORS }),
    },
    hideSuffix: {
      type: Boolean,
      default: false,
    },
    colorType: {
      type: String,
      default: 'rgba',
    },
  },
  methods: {
    getButtonText(key) {
      let suffix = '';

      if (!this.hideSuffix) {
        if (key === STATS_CRITICITY.ok) {
          suffix = ` / ${this.$t('common.yes')}`;
        } else if (key === STATS_CRITICITY.critical) {
          suffix = ` / ${this.$t('common.no')}`;
        }
      }

      return this.$t(`settings.colorsSelector.statsCriticity.${key}`) + suffix;
    },
  },
};
</script>
