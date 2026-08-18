<template>
  <widget-settings-item :title="$t('settings.colorsSelector.title')">
    <c-layout column>
      <c-color-picker-field
        v-for="level in levels"
        :key="level"
        :model-value="modelValue[level]"
        :label="getButtonText(level)"
        :type="colorType"
        class="mb-3"
        splitted
        @update:model-value="updateLevel(level, $event)"
      />
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { ALARM_LEVELS_COLORS, STATS_CRITICITY } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
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

  emits: ['update:modelValue'],

  computed: {
    levels() {
      return Object.values(STATS_CRITICITY);
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

    updateLevel(level, value) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        [level]: value,
      });
    },
  },
};
</script>
