<template>
  <v-row class="theme-form__columns" no-gutters>
    <v-col class="theme-form__column" cols="12" md="6">
      <theme-color-picker-field
        v-model="primaryField"
        :label="t('theme.main.primary')"
        :help-text="t('theme.main.primaryHelpText')"
      />
      <theme-color-picker-field
        v-model="secondaryField"
        :label="t('theme.main.secondary')"
        :help-text="t('theme.main.secondaryHelpText')"
      />
      <theme-color-picker-field v-model="accentField" :label="t('theme.main.accent')" />
      <theme-color-picker-field v-model="backgroundField" :label="t('theme.main.background')" />
      <theme-color-picker-field
        v-model="activeColorField"
        :label="t('theme.main.activeColor')"
        :help-text="t('theme.main.activeColorHelpText')"
      />
    </v-col>

    <v-col class="theme-form__column" cols="12" md="6">
      <theme-color-picker-field v-model="errorField" :label="t('theme.main.error')" />
      <theme-color-picker-field v-model="errorBackgroundField" :label="t('theme.main.errorBackground')" />
      <theme-color-picker-field v-model="warningField" :label="t('theme.main.warning')" />
      <theme-color-picker-field v-model="warningBackgroundField" :label="t('theme.main.warningBackground')" />
      <theme-color-picker-field v-model="successField" :label="t('theme.main.success')" />
      <theme-color-picker-field v-model="successBackgroundField" :label="t('theme.main.successBackground')" />
      <theme-color-picker-field v-model="infoField" :label="t('theme.main.info')" />
      <theme-color-picker-field v-model="infoBackgroundField" :label="t('theme.main.infoBackground')" />
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
    const {
      primaryField,
      secondaryField,
      accentField,
      backgroundField,
      activeColorField,
      errorField,
      errorBackgroundField,
      warningField,
      warningBackgroundField,
      successField,
      successBackgroundField,
      infoField,
      infoBackgroundField,
    } = useFormModels({
      primaryField: 'colors.main.primary',
      secondaryField: 'colors.main.secondary',
      accentField: 'colors.main.accent',
      backgroundField: 'colors.main.background',
      activeColorField: 'colors.main.active_color',
      errorField: 'colors.main.error',
      errorBackgroundField: 'colors.main.error_background',
      warningField: 'colors.main.warning',
      warningBackgroundField: 'colors.main.warning_background',
      successField: 'colors.main.success',
      successBackgroundField: 'colors.main.success_background',
      infoField: 'colors.main.info',
      infoBackgroundField: 'colors.main.info_background',
    });

    return {
      t,
      primaryField,
      secondaryField,
      accentField,
      backgroundField,
      activeColorField,
      errorField,
      errorBackgroundField,
      warningField,
      warningBackgroundField,
      successField,
      successBackgroundField,
      infoField,
      infoBackgroundField,
    };
  },
};
</script>
