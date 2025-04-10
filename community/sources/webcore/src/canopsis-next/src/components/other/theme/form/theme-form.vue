<template>
  <v-layout
    class="theme-form"
    column
  >
    <c-name-field v-field="form.name" autofocus required />
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
        />
        <theme-color-picker-field
          v-field="form.colors.main.error_icons"
          :label="$t('theme.main.errorIcons')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.error"
          :label="$t('theme.main.error')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.warning_icons"
          :label="$t('theme.main.warningIcons')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.warning"
          :label="$t('theme.main.warning')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.success_icons"
          :label="$t('theme.main.successIcons')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.success"
          :label="$t('theme.main.success')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.info_icons"
          :label="$t('theme.main.infoIcons')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.info"
          :label="$t('theme.main.info')"
        />
        <theme-color-picker-field
          v-field="form.colors.main.background"
          :label="$t('theme.main.background')"
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
        />
        <theme-color-picker-field
          v-field="form.colors.table.row_color"
          :label="$t('theme.table.rowColor')"
        />
        <theme-enabled-color-picker-field
          v-field="form.colors.table.shift_row_color"
          :enable-label="$t('theme.table.shiftRowEnable')"
          :enable-help-text="$t('theme.table.shiftRowEnableHelpText')"
          :label="$t('theme.table.shiftRowColor')"
        />
        <theme-enabled-color-picker-field
          v-field="form.colors.table.hover_row_color"
          :enable-label="$t('theme.table.hoverRowEnable')"
          :label="$t('theme.table.hoverRowColor')"
        />
        <theme-colors-preview-list :items="tableColors" />
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
        />
        <theme-color-picker-field
          v-field="form.colors.state.minor"
          :label="$t('theme.state.minor')"
        />
        <theme-color-picker-field
          v-field="form.colors.state.major"
          :label="$t('theme.state.major')"
        />
        <theme-color-picker-field
          v-field="form.colors.state.critical"
          :label="$t('theme.state.critical')"
        />
        <theme-colors-preview-list :items="stateColors" />
      </v-layout>
    </c-information-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { THEME_FONT_PIXEL_SIZES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import ThemeEnabledColorPickerField from '@/components/other/theme/form/fields/theme-enabled-color-picker-field.vue';
import ThemeColorFontSizeField from '@/components/other/theme/form/fields/theme-color-font-size-field.vue';

import ThemeColorsPreviewList from '../partials/theme-colors-preview-list.vue';

import ThemeColorPickerField from './fields/theme-color-picker-field.vue';

export default {
  inject: ['$validator'],
  components: {
    ThemeColorFontSizeField,
    ThemeColorsPreviewList,
    ThemeEnabledColorPickerField,
    ThemeColorPickerField,
  },
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
  setup(props) {
    const { t } = useI18n();

    const fontSize = computed(() => THEME_FONT_PIXEL_SIZES[props.form.font_size]);

    const attachFontSizeAndActiveColor = items => (
      items.map(item => ({ color: '#000', fontSize: fontSize.value, ...item }))
    );

    const tableColors = computed(() => {
      const result = [
        {
          backgroundColor: props.form.colors.table.background,
          text: t('theme.table.background'),
        },
        {
          backgroundColor: props.form.colors.table.row_color,
          text: t('theme.table.rowColor'),
        },
      ];

      if (props.form.colors.table.shift_row_color?.enabled) {
        result.push({
          backgroundColor: props.form.colors.table.shift_row_color.color,
          text: t('theme.table.shiftRowColor'),
        });
      }

      if (props.form.colors.table.hover_row_color?.enabled) {
        result.push({
          backgroundColor: props.form.colors.table.hover_row_color.color,
          text: t('theme.table.hoverRowColor'),
        });
      }

      return attachFontSizeAndActiveColor(result);
    });

    const stateColors = computed(() => attachFontSizeAndActiveColor([
      {
        backgroundColor: props.form.colors.state.ok,
        text: t('theme.state.ok'),
        color: '#fff',
      },
      {
        backgroundColor: props.form.colors.state.minor,
        text: t('theme.state.minor'),
        color: '#000',
      },
      {
        backgroundColor: props.form.colors.state.major,
        text: t('theme.state.major'),
        color: '#000',
      },
      {
        backgroundColor: props.form.colors.state.critical,
        text: t('theme.state.critical'),
        color: '#fff',
      },
    ]));

    return {
      tableColors,
      stateColors,
    };
  },
};
</script>

<style lang="scss">
.theme-form {
  --item-divider-border: 1px solid var(--v-divider-border-color);

  &__colors {
    padding: 16px 0 24px 0;

    & > * {
      padding: 10px 0 10px 16px;

      &:not(:last-of-type) {
        border-bottom: var(--item-divider-border);
      }
    }
  }
}
</style>
