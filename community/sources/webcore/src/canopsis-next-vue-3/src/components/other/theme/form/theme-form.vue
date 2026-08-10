<template>
  <c-layout class="theme-form" column>
    <c-name-field v-model="nameField" autofocus required />

    <c-information-block :title="t('theme.main.title')">
      <c-layout class="theme-form__colors" column>
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
        <theme-color-picker-field v-model="errorField" :label="t('theme.main.error')" />
        <theme-color-picker-field
          v-model="errorBackgroundField"
          :label="t('theme.main.errorBackground')"
        />
        <theme-color-picker-field v-model="warningField" :label="t('theme.main.warning')" />
        <theme-color-picker-field
          v-model="warningBackgroundField"
          :label="t('theme.main.warningBackground')"
        />
        <theme-color-picker-field v-model="successField" :label="t('theme.main.success')" />
        <theme-color-picker-field
          v-model="successBackgroundField"
          :label="t('theme.main.successBackground')"
        />
        <theme-color-picker-field v-model="infoField" :label="t('theme.main.info')" />
        <theme-color-picker-field
          v-model="infoBackgroundField"
          :label="t('theme.main.infoBackground')"
        />
        <theme-color-picker-field
          v-model="activeColorField"
          :label="t('theme.main.activeColor')"
          :help-text="t('theme.main.activeColorHelpText')"
        />
        <theme-color-picker-field v-model="backgroundField" :label="t('theme.main.background')" />
      </c-layout>
    </c-information-block>

    <c-information-block :title="t('theme.fontSize.title')">
      <c-layout class="theme-form__colors" column>
        <theme-color-font-size-field v-model="fontSizeField" />
      </c-layout>
    </c-information-block>

    <c-information-block :title="t('theme.table.title')">
      <c-layout class="theme-form__colors" column>
        <theme-color-picker-field
          v-model="tableBackgroundField"
          :label="t('theme.table.background')"
        />
        <theme-color-picker-field v-model="tableRowColorField" :label="t('theme.table.rowColor')" />
        <theme-enabled-color-picker-field
          v-model="tableShiftRowField"
          :enable-label="t('theme.table.shiftRowEnable')"
          :enable-help-text="t('theme.table.shiftRowEnableHelpText')"
          :label="t('theme.table.shiftRowColor')"
        />
        <theme-enabled-color-picker-field
          v-model="tableHoverRowField"
          :enable-label="t('theme.table.hoverRowEnable')"
          :label="t('theme.table.hoverRowColor')"
        />
        <theme-colors-preview-list :items="tableColors" />
      </c-layout>
    </c-information-block>

    <c-information-block :title="t('theme.state.title')">
      <c-layout class="theme-form__colors" column>
        <theme-color-picker-field v-model="stateOkField" :label="t('theme.state.ok')" />
        <theme-color-picker-field v-model="stateMinorField" :label="t('theme.state.minor')" />
        <theme-color-picker-field v-model="stateMajorField" :label="t('theme.state.major')" />
        <theme-color-picker-field v-model="stateCriticalField" :label="t('theme.state.critical')" />
        <theme-colors-preview-list :items="stateColors" />
      </c-layout>
    </c-information-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { THEME_FONT_PIXEL_SIZES } from '@/constants';
import { getMostReadableTextColor } from '@/helpers/color';
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ThemeColorsPreviewList from '../partials/theme-colors-preview-list.vue';

import ThemeColorFontSizeField from './fields/theme-color-font-size-field.vue';
import ThemeColorPickerField from './fields/theme-color-picker-field.vue';
import ThemeEnabledColorPickerField from './fields/theme-enabled-color-picker-field.vue';

export default {
  components: {
    ThemeColorFontSizeField,
    ThemeColorPickerField,
    ThemeColorsPreviewList,
    ThemeEnabledColorPickerField,
  },

  setup() {
    const { t } = useI18n();
    const {
      nameField,
      fontSizeField,
      primaryField,
      secondaryField,
      accentField,
      errorField,
      errorBackgroundField,
      warningField,
      warningBackgroundField,
      successField,
      successBackgroundField,
      infoField,
      infoBackgroundField,
      activeColorField,
      backgroundField,
      tableBackgroundField,
      tableRowColorField,
      tableShiftRowField,
      tableHoverRowField,
      stateOkField,
      stateMinorField,
      stateMajorField,
      stateCriticalField,
    } = useFormModels({
      nameField: 'name',
      fontSizeField: 'font_size',
      primaryField: 'colors.main.primary',
      secondaryField: 'colors.main.secondary',
      accentField: 'colors.main.accent',
      errorField: 'colors.main.error',
      errorBackgroundField: 'colors.main.error_background',
      warningField: 'colors.main.warning',
      warningBackgroundField: 'colors.main.warning_background',
      successField: 'colors.main.success',
      successBackgroundField: 'colors.main.success_background',
      infoField: 'colors.main.info',
      infoBackgroundField: 'colors.main.info_background',
      activeColorField: 'colors.main.active_color',
      backgroundField: 'colors.main.background',
      tableBackgroundField: 'colors.table.background',
      tableRowColorField: 'colors.table.row_color',
      tableShiftRowField: 'colors.table.shift_row',
      tableHoverRowField: 'colors.table.hover_row',
      stateOkField: 'colors.state.ok',
      stateMinorField: 'colors.state.minor',
      stateMajorField: 'colors.state.major',
      stateCriticalField: 'colors.state.critical',
    });

    const fontSize = computed(() => THEME_FONT_PIXEL_SIZES[fontSizeField.value]);
    const attachFontSizeAndActiveColor = items =>
      items.map(item => ({
        color: activeColorField.value,
        fontSize: fontSize.value,
        ...item,
      }));

    const tableColors = computed(() => {
      const result = [
        {
          backgroundColor: tableBackgroundField.value,
          text: t('theme.table.background'),
        },
        {
          backgroundColor: tableRowColorField.value,
          text: t('theme.table.rowColor'),
        },
      ];

      if (tableShiftRowField.value?.enabled) {
        result.push({
          backgroundColor: tableShiftRowField.value.color,
          text: t('theme.table.shiftRowColor'),
        });
      }

      if (tableHoverRowField.value?.enabled) {
        result.push({
          backgroundColor: tableHoverRowField.value.color,
          text: t('theme.table.hoverRowColor'),
        });
      }

      return attachFontSizeAndActiveColor(result);
    });

    const stateColors = computed(() =>
      attachFontSizeAndActiveColor([
        {
          backgroundColor: stateOkField.value,
          text: t('theme.state.ok'),
          color: getMostReadableTextColor(stateOkField.value),
        },
        {
          backgroundColor: stateMinorField.value,
          text: t('theme.state.minor'),
          color: getMostReadableTextColor(stateMinorField.value),
        },
        {
          backgroundColor: stateMajorField.value,
          text: t('theme.state.major'),
          color: getMostReadableTextColor(stateMajorField.value),
        },
        {
          backgroundColor: stateCriticalField.value,
          text: t('theme.state.critical'),
          color: getMostReadableTextColor(stateCriticalField.value),
        },
      ]),
    );

    return {
      t,
      nameField,
      fontSizeField,
      primaryField,
      secondaryField,
      accentField,
      errorField,
      errorBackgroundField,
      warningField,
      warningBackgroundField,
      successField,
      successBackgroundField,
      infoField,
      infoBackgroundField,
      activeColorField,
      backgroundField,
      tableBackgroundField,
      tableRowColorField,
      tableShiftRowField,
      tableHoverRowField,
      stateOkField,
      stateMinorField,
      stateMajorField,
      stateCriticalField,
      tableColors,
      stateColors,
    };
  },
};
</script>

<style lang="scss">
.theme-form {
  --item-divider-border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));

  &__colors {
    padding: 16px 0 24px;

    & > * {
      padding: 10px 0 10px 16px;

      &:not(:last-of-type) {
        border-bottom: var(--item-divider-border);
      }
    }
  }
}
</style>
