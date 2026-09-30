<template>
  <widget-settings-item :title="$t('settings.colorsSelector.title')">
    <c-color-picker-field
      v-for="level in levels"
      :key="level"
      :model-value="modelValue[level]"
      :label="getButtonText(level)"
      :type="colorType"
      class="mb-3 w-100"
      fill-width
      @update:model-value="updateLevel(level, $event)"
    >
      <template #activator="{ props, color, label, style, disabled }">
        <v-row class="w-100" align="center" justify="space-between">
          <v-btn v-bind="props" :disabled="disabled">
            {{ label }}
          </v-btn>

          <div :style="style" class="pa-1 text-center">
            {{ color }}
          </div>
        </v-row>
      </template>
    </c-color-picker-field>
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import { ALARM_LEVELS_COLORS, STATS_CRITICITY } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

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

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);
    const levels = computed(() => Object.values(STATS_CRITICITY));

    const getButtonText = (key) => {
      let suffix = '';

      if (!props.hideSuffix) {
        if (key === STATS_CRITICITY.ok) {
          suffix = ` / ${t('common.yes')}`;
        } else if (key === STATS_CRITICITY.critical) {
          suffix = ` / ${t('common.no')}`;
        }
      }

      return t(`settings.colorsSelector.statsCriticity.${key}`) + suffix;
    };

    const updateLevel = (level, value) => {
      updateModel({
        ...props.modelValue,
        [level]: value,
      });
    };

    return {
      levels,
      getButtonText,
      updateLevel,
    };
  },
};
</script>
