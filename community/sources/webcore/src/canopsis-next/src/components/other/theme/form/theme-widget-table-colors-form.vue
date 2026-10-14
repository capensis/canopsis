<template>
  <v-row class="theme-form__colors" no-gutters>
    <v-col cols="12">
      <span class="text-subtitle-2">{{ t('theme.table.title') }}</span>

    <theme-color-picker-field
      v-model="tableBackgroundField"
      :label="t('theme.table.background')"
      :color-label="t('common.text')"
    />
    <theme-color-picker-field
      v-model="tableRowColorField"
      :label="t('theme.table.rowColor')"
      :color-label="t('common.text')"
    />
    <theme-enabled-color-picker-field
      v-model="tableShiftRowField"
      :enable-label="t('theme.table.shiftRowEnable')"
      :enable-help-text="t('theme.table.shiftRowEnableHelpText')"
      :color-label="t('common.text')"
    >
    </theme-enabled-color-picker-field>
    <theme-enabled-color-picker-field
      v-model="tableHoverRowField"
      :enable-label="t('theme.table.hoverRowEnable')"
      :color-label="t('common.text')"
    >
    </theme-enabled-color-picker-field>

    <span class="text-subtitle-2 mt-3">{{ t('theme.widget.title') }}</span>

    <theme-color-picker-field
      v-model="widgetTitleBackgroundField"
      :label="t('theme.widget.titleBackground')"
      :color-label="t('common.text')"
    />
    <theme-color-picker-field
      v-model="widgetContentBackgroundField"
      :label="t('theme.widget.contentBackground')"
      :color-label="t('common.text')"
    />
    <theme-enabled-color-picker-field
      v-model="widgetBorderField"
      :enable-label="t('theme.widget.border')"
    >
      <template #append>
        <c-number-field
          v-if="widgetBorderField.enabled"
          v-model:value="widgetBorderField.radius"
          :label="t('theme.widget.borderRadius')"
          :min="0"
          name="colors.widget.border.radius"
          class="theme-widget-table-colors-form__border-radius"
          type="number"
          hide-details
        />
      </template>
    </theme-enabled-color-picker-field>

      <widget-table-colors-preview
        :table="tablePreview"
        :widget="widgetPreview"
        :font-style="fontStyle"
        class="mt-3"
      />
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { DEFAULT_WIDGET_BORDER_RADIUS, THEME_FONT_PIXEL_SIZES } from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import WidgetTableColorsPreview from '@/components/widgets/partials/widget-table-colors-preview.vue';

import ThemeColorPickerField from './fields/theme-color-picker-field.vue';
import ThemeEnabledColorPickerField from './fields/theme-enabled-color-picker-field.vue';

export default {
  components: {
    ThemeColorPickerField,
    ThemeEnabledColorPickerField,
    WidgetTableColorsPreview,
  },

  setup() {
    const { t } = useI18n();
    const {
      fontSizeField,
      activeColorField,
      tableBackgroundField,
      tableRowColorField,
      tableShiftRowField,
      tableHoverRowField,
      widgetTitleBackgroundField,
      widgetContentBackgroundField,
      widgetBorderField: widgetBorderModel,
    } = useFormModels({
      fontSizeField: 'font_size',
      activeColorField: 'colors.main.active_color',
      tableBackgroundField: 'colors.table.background',
      tableRowColorField: 'colors.table.row_color',
      tableShiftRowField: 'colors.table.shift_row',
      tableHoverRowField: 'colors.table.hover_row',
      widgetTitleBackgroundField: 'colors.widget.title_background',
      widgetContentBackgroundField: 'colors.widget.content_background',
      widgetBorderField: 'colors.widget.border',
    });
    const widgetBorderField = computed({
      get: () => widgetBorderModel.value,
      set: (value) => {
        widgetBorderModel.value = value.enabled && !value.radius
          ? { ...value, radius: DEFAULT_WIDGET_BORDER_RADIUS }
          : value;
      },
    });

    const fontStyle = computed(() => ({
      color: activeColorField.value,
      fontSize: `${THEME_FONT_PIXEL_SIZES[fontSizeField.value]}px`,
    }));
    const tablePreview = computed(() => ({
      background: tableBackgroundField.value,
      row_color: tableRowColorField.value,
      shift_row: tableShiftRowField.value,
      hover_row: tableHoverRowField.value,
    }));
    const widgetPreview = computed(() => ({
      title_background: widgetTitleBackgroundField.value,
      content_background: widgetContentBackgroundField.value,
      border: widgetBorderField.value,
    }));

    return {
      t,
      fontStyle,
      tablePreview,
      widgetPreview,
      tableBackgroundField,
      tableRowColorField,
      tableShiftRowField,
      tableHoverRowField,
      widgetTitleBackgroundField,
      widgetContentBackgroundField,
      widgetBorderField,
    };
  },
};
</script>

<style lang="scss">
.theme-widget-table-colors-form {
  &__border-radius {
    flex: 0 0 100px;
    margin-right: 8px;
  }
}
</style>
