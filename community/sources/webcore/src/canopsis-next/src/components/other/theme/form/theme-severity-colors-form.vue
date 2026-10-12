<template>
  <v-row class="theme-form__colors" no-gutters>
    <v-col cols="12">
      <theme-color-picker-field v-model="okField" :label="t('theme.state.ok')" />
      <theme-color-picker-field v-model="minorField" :label="t('theme.state.minor')" />
      <theme-color-picker-field v-model="majorField" :label="t('theme.state.major')" />
      <theme-color-picker-field v-model="criticalField" :label="t('theme.state.critical')" />
    </v-col>
  </v-row>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ThemeColorPickerField from './fields/theme-color-picker-field.vue';

export default {
  components: { ThemeColorPickerField },

  setup() {
    const { t } = useI18n();
    const { okField, minorField, majorField, criticalField } = useFormModels({
      okField: 'colors.state.ok',
      minorField: 'colors.state.minor',
      majorField: 'colors.state.major',
      criticalField: 'colors.state.critical',
    });

    return { t, okField, minorField, majorField, criticalField };
  },
};
</script>
