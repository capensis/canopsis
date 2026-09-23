<template>
  <widget-settings-item :title="t('settings.colors.title')">
    <c-select
      v-model="mode"
      :items="colorModes"
      :label="t('settings.colors.title')"
      item-title="text"
      item-value="value"
      hide-details
      single-line
    />

    <v-expand-transition>
      <widget-colors-custom-form
        v-if="mode === WIDGET_COLORS_MODES.custom"
        v-model="customColors"
        :with-table="withTable"
      />
    </v-expand-transition>
  </widget-settings-item>
</template>

<script>
import { computed, inject } from 'vue';

import { WIDGET_COLORS_MODES } from '@/constants';

import { themeTableColorsToForm, themeWidgetColorsToForm } from '@/helpers/entities/theme/form';

import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

import WidgetColorsCustomForm from './widget-colors-custom-form.vue';

export default {
  components: {
    WidgetColorsCustomForm,
    WidgetSettingsItem,
  },

  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    withTable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const system = inject('$system', {});

    const updateMode = (value) => {
      if (value !== WIDGET_COLORS_MODES.custom) {
        emit('update:modelValue', { ...props.modelValue, mode: value });

        return;
      }

      const themeColors = system.theme?.colors ?? {};

      emit('update:modelValue', {
        ...props.modelValue,
        mode: value,
        table: themeTableColorsToForm(themeColors.table),
        widget: themeWidgetColorsToForm(themeColors.widget),
      });
    };

    const mode = computed({
      get: () => props.modelValue.mode,
      set: updateMode,
    });
    const customColors = computed({
      get: () => props.modelValue,
      set: value => emit('update:modelValue', value),
    });
    const colorModes = computed(() => [
      { value: WIDGET_COLORS_MODES.theme, text: t('settings.colors.theme') },
      { value: WIDGET_COLORS_MODES.custom, text: t('settings.colors.custom') },
    ]);

    return {
      t,
      mode,
      customColors,
      colorModes,
      WIDGET_COLORS_MODES,
    };
  },
};
</script>
