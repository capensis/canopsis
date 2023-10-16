<template>
  <v-layout
    class="theme-form"
    column
  >
    <c-name-field v-field="form.name" />
    <c-information-block :title="$t('theme.main.title')">
      <v-layout
        class="theme-form__colors"
        column
      >
        <theme-color-picker-field
          v-field="form.colors.main.primary"
          :label="$t('theme.main.primary')"
          :help-text="$t('theme.main.primaryHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.secondary"
          :label="$t('theme.main.secondary')"
          :help-text="$t('theme.main.secondaryHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.accent"
          :label="$t('theme.main.accent')"
          :help-text="$t('theme.main.secondaryHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.error"
          :label="$t('theme.main.error')"
          :help-text="$t('theme.main.secondaryHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.info"
          :label="$t('theme.main.info')"
          :help-text="$t('theme.main.secondaryHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.success"
          :label="$t('theme.main.success')"
          :help-text="$t('theme.main.secondaryHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.warning"
          :label="$t('theme.main.warning')"
          :help-text="$t('theme.main.secondaryHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.background"
          :label="$t('theme.main.background')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.active_color"
          :label="$t('theme.main.activeColor')"
          :help-text="$t('theme.main.activeColorHelpText')"
        />
        <theme-colors-preview
          :background-color="form.colors.main.background"
          :color="form.colors.main.active_color"
          :font-size="fontSize"
        />
      </v-layout>
    </c-information-block>
    <c-information-block :title="$t('theme.fontSize.title')">
      <v-layout
        class="theme-form__colors"
        column
      >
        <theme-color-font-size-field v-field="form.font_size" />
      </v-layout>
    </c-information-block>
    <c-information-block :title="$t('theme.table.title')">
      <v-layout
        class="theme-form__colors"
        column
      >
        <theme-color-picker-field
          v-field="form.colors.table.background"
          :label="$t('theme.table.background')"
          :help-text="$t('theme.table.backgroundHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.table.row_color"
          :label="$t('theme.table.rowColor')"
          :help-text="$t('theme.table.rowColorHelpText')"
        />
        <theme-enabled-color-picker-field
          v-field="form.colors.table.shift_row_color"
          :enable-label="$t('theme.table.shiftRowEnable')"
          :enable-help-text="$t('theme.table.shiftRowEnableHelpText')"
          :label="$t('theme.table.shiftRowColor')"
          :help-text="$t('theme.table.shiftRowColorHelpText')"
        />
        <theme-enabled-color-picker-field
          v-field="form.colors.table.hover_row_color"
          :enable-label="$t('theme.table.hoverRowEnable')"
          :enable-help-text="$t('theme.table.hoverRowEnableHelpText')"
          :label="$t('theme.table.hoverRowColor')"
        />
        <theme-colors-preview
          :background-color="form.colors.table.background"
          :color="form.colors.main.active_color"
          :font-size="fontSize"
        />
      </v-layout>
    </c-information-block>
    <c-information-block :title="$t('theme.state.title')">
      <v-layout
        class="theme-form__colors"
        column
      >
        <theme-color-picker-field
          v-field="form.colors.state.ok"
          :label="$t('theme.state.ok')"
          :help-text="$t('theme.state.okHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.state.minor"
          :label="$t('theme.state.minor')"
          :help-text="$t('theme.state.minorHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.state.major"
          :label="$t('theme.state.major')"
          :help-text="$t('theme.state.majorHelpText')"
        />
        <theme-color-picker-field
          v-field="form.colors.state.critical"
          :label="$t('theme.state.critical')"
          :help-text="$t('theme.state.criticalHelpText')"
        />
      </v-layout>
    </c-information-block>
  </v-layout>
</template>

<script>
import ThemeEnabledColorPickerField from '@/components/other/theme/form/fields/theme-enabled-color-picker-field.vue';
import ThemeColorFontSizeField from '@/components/other/theme/form/fields/theme-color-font-size-field.vue';

import { THEME_FONT_PIXEL_SIZES } from '@/constants/theme';

import ThemeColorsPreview from '../partials/theme-colors-preview.vue';

import ThemeColorPickerField from './fields/theme-color-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { ThemeColorFontSizeField, ThemeColorsPreview, ThemeEnabledColorPickerField, ThemeColorPickerField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    fontSize() {
      return THEME_FONT_PIXEL_SIZES[this.form.font_size];
    },
  },
};
</script>

<style lang="scss">
.theme-form {
  &__colors {
    padding: 16px 0 24px 16px;
    gap: 15px;
  }
}
</style>
