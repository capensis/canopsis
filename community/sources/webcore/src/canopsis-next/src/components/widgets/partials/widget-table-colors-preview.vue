<template>
  <div class="widget-table-colors-preview">
    <span class="text-subtitle-2 d-block mb-2">{{ t('theme.preview.title') }}</span>
    <v-card :style="widgetBorderStyle" class="widget-table-colors-preview__table">
      <div
        :style="widgetTitleRowStyle"
        class="widget-table-colors-preview__row widget-table-colors-preview__row--widget-title"
      >
        {{ t('theme.preview.widgetTitle') }}
      </div>
      <div v-if="withTable" :style="headerRowStyle" class="widget-table-colors-preview__row">
        {{ t('theme.preview.header') }}
      </div>
      <div v-if="withTable" :style="tableRowStyle" class="widget-table-colors-preview__row">
        {{ t('theme.preview.tableRow') }}
      </div>
      <div
        v-if="withTable && tableModel.shift_row?.enabled"
        :style="shiftRowStyle"
        class="widget-table-colors-preview__row"
      >
        {{ t('theme.preview.secondTableRow') }}
      </div>
      <div
        v-if="withTable && tableModel.hover_row?.enabled"
        :style="hoverRowStyle"
        class="widget-table-colors-preview__row"
      >
        {{ t('theme.preview.hoveredRow') }}
      </div>
    </v-card>
  </div>
</template>

<script>
import { computed, inject } from 'vue';

import { THEME_FONT_PIXEL_SIZES, THEME_FONT_SIZES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    table: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    fontStyle: {
      type: Object,
      default: undefined,
    },
    withTable: {
      type: Boolean,
      default: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const system = inject('$system', {});
    const tableModel = computed(() => props.table);
    const currentFontStyle = computed(() => props.fontStyle ?? {
      color: system.theme?.colors?.main?.active_color,
      fontSize: `${THEME_FONT_PIXEL_SIZES[system.theme?.font_size]
        ?? THEME_FONT_PIXEL_SIZES[THEME_FONT_SIZES.medium]}px`,
    });
    const getRowStyle = backgroundColor => ({
      backgroundColor,
      ...currentFontStyle.value,
    });
    const widgetTitleRowStyle = computed(() => getRowStyle(props.widget.title_background));
    const headerRowStyle = computed(() => getRowStyle(props.table.background));
    const tableRowStyle = computed(() => getRowStyle(props.table.row_color));
    const shiftRowStyle = computed(() => getRowStyle(props.table.shift_row?.color));
    const hoverRowStyle = computed(() => getRowStyle(props.table.hover_row?.color));
    const widgetBorderStyle = computed(() => {
      if (!props.widget.border?.enabled) {
        return {};
      }

      return {
        border: `1px solid ${props.widget.border.color}`,
        borderRadius: `${props.widget.border.radius}px`,
      };
    });

    return {
      t,
      tableModel,
      widgetTitleRowStyle,
      headerRowStyle,
      tableRowStyle,
      shiftRowStyle,
      hoverRowStyle,
      widgetBorderStyle,
    };
  },
};
</script>

<style lang="scss">
.widget-table-colors-preview {
  &__table {
    overflow: hidden;
  }

  &__row {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 44px;
    padding: 0 16px;
    line-height: 1.25;
    text-transform: uppercase;
    border-bottom: 1px solid var(--v-border-color-base, rgba(0, 0, 0, 0.12));

    &--widget-title {
      justify-content: flex-start;
      text-transform: none;
      font-size: 1.125rem;
      font-weight: 600;
    }

    &:last-child {
      border-bottom: none;
    }
  }
}
</style>
