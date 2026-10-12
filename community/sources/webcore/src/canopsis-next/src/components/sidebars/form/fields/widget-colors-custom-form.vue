<template>
  <div class="widget-colors-custom-form">
    <template v-if="withTable">
      <span class="widget-colors-custom-form__section-title text-subtitle-2">
        {{ t('theme.table.title') }}
      </span>

      <widget-settings-flat-item :title="t('theme.table.background')">
        <template #actions>
          <c-color-picker-field
            :model-value="table.background"
            :label="t('common.text')"
            class="widget-colors-custom-form__color-picker"
            size="small"
            justify-end
            hide-details
            @update:model-value="updateField('table.background', $event)"
          />
        </template>
      </widget-settings-flat-item>
      <widget-settings-flat-item :title="t('theme.table.rowColor')">
        <template #actions>
          <c-color-picker-field
            :model-value="table.row_color"
            :label="t('common.text')"
            class="widget-colors-custom-form__color-picker"
            size="small"
            justify-end
            hide-details
            @update:model-value="updateField('table.row_color', $event)"
          />
        </template>
      </widget-settings-flat-item>
      <widget-settings-flat-item>
        <template #title>
          <v-row density="comfortable" align="center">
            {{ t('theme.table.shiftRowColor') }}
            <c-help-icon :text="t('theme.table.shiftRowEnableHelpText')" icon="help" top />
            <c-enabled-field
              :model-value="table.shift_row?.enabled"
              :aria-label="t('theme.table.shiftRowColor')"
              hide-label
              hide-details
              @update:model-value="updateField('table.shift_row.enabled', $event)"
            />
          </v-row>
        </template>
        <template #actions>
          <c-color-picker-field
            v-if="table.shift_row?.enabled"
            :model-value="table.shift_row?.color"
            :label="t('common.text')"
            class="widget-colors-custom-form__color-picker"
            size="small"
            justify-end
            hide-details
            @update:model-value="updateField('table.shift_row.color', $event)"
          />
        </template>
      </widget-settings-flat-item>
      <widget-settings-flat-item :title="t('theme.table.hoverRowEnable')">
        <template #title>
          <v-row density="comfortable" align="center">
            {{ t('theme.table.hoverRowEnable') }}
            <c-enabled-field
              :model-value="table.hover_row?.enabled"
              :aria-label="t('theme.table.hoverRowEnable')"
              hide-label
              hide-details
              @update:model-value="updateField('table.hover_row.enabled', $event)"
            />
          </v-row>
        </template>
        <template #actions>
          <v-row density="comfortable" align="center">
            <c-color-picker-field
              v-if="table.hover_row?.enabled"
              :model-value="table.hover_row?.color"
              :label="t('common.text')"
              class="widget-colors-custom-form__color-picker"
              size="small"
              justify-end
              hide-details
              @update:model-value="updateField('table.hover_row.color', $event)"
            />
          </v-row>
        </template>
      </widget-settings-flat-item>
    </template>

    <span class="widget-colors-custom-form__section-title text-subtitle-2">
      {{ t('theme.widget.title') }}
    </span>

    <widget-settings-flat-item :title="t('theme.widget.titleBackground')">
      <template #actions>
        <c-color-picker-field
          :model-value="widget.title_background"
          :label="t('common.text')"
          class="widget-colors-custom-form__color-picker"
          size="small"
          justify-end
          hide-details
          @update:model-value="updateField('widget.title_background', $event)"
        />
      </template>
    </widget-settings-flat-item>
    <widget-settings-flat-item :title="t('theme.widget.contentBackground')">
      <template #actions>
        <c-color-picker-field
          :model-value="widget.content_background"
          :label="t('common.text')"
          class="widget-colors-custom-form__color-picker"
          size="small"
          justify-end
          hide-details
          @update:model-value="updateField('widget.content_background', $event)"
        />
      </template>
    </widget-settings-flat-item>
    <widget-settings-flat-item>
      <template #title>
        <v-row density="comfortable" align="center">
          {{ t('theme.widget.border') }}
          <c-enabled-field
            :model-value="widget.border?.enabled"
            :aria-label="t('theme.widget.border')"
            hide-label
            inline
            @update:model-value="updateWidgetBorder('enabled', $event)"
          />
        </v-row>
      </template>
      <template #actions>
        <v-row density="comfortable" align="center">
          <c-number-field
            v-if="widget.border?.enabled"
            :value="widget.border?.radius"
            :label="t('theme.widget.borderRadius')"
            :min="0"
            class="widget-colors-custom-form__border-radius"
            type="number"
            hide-details
            @update:value="updateField('widget.border.radius', $event)"
          />
          <c-color-picker-field
            :model-value="widget.border?.color"
            class="widget-colors-custom-form__color-picker"
            size="small"
            hide-label
            justify-end
            hide-details
            @update:model-value="updateField('widget.border.color', $event)"
          />
        </v-row>
      </template>
    </widget-settings-flat-item>

    <widget-table-colors-preview
      :table="tablePreview"
      :widget="widgetPreview"
      :with-table="withTable"
      class="widget-colors-custom-form__preview"
    />
  </div>
</template>

<script>
import { computed } from 'vue';

import { DEFAULT_WIDGET_BORDER_RADIUS } from '@/constants';

import { useModelField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import WidgetSettingsFlatItem from '@/components/sidebars/partials/widget-settings-flat-item.vue';
import WidgetTableColorsPreview from '@/components/widgets/partials/widget-table-colors-preview.vue';

export default {
  components: {
    WidgetTableColorsPreview,
    WidgetSettingsFlatItem,
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
    const { updateField } = useModelField(props, emit);

    const table = computed(() => props.modelValue.table ?? {});
    const widget = computed(() => props.modelValue.widget ?? {});
    const tablePreview = computed(() => ({
      background: table.value.background,
      row_color: table.value.row_color,
      shift_row: table.value.shift_row,
      hover_row: table.value.hover_row,
    }));
    const widgetPreview = computed(() => ({
      title_background: widget.value.title_background,
      content_background: widget.value.content_background,
      border: widget.value.border,
    }));
    const updateWidgetBorder = (field, value) => {
      const nextValue = field === 'enabled' && value && !widget.value.border?.radius
        ? DEFAULT_WIDGET_BORDER_RADIUS
        : value;

      updateField(`widget.border.${field}`, nextValue);
    };

    return {
      t,
      table,
      widget,
      updateField,
      tablePreview,
      widgetPreview,
      updateWidgetBorder,
    };
  },
};
</script>

<style lang="scss">
.widget-colors-custom-form {
  &__section-title {
    display: block;
    margin: 16px 0 4px;
  }

  &__border-radius {
    max-width: 72px;
    max-height: 50px;
    margin: 0 !important;
  }

  &__color-picker .v-btn {
    margin: 0 !important;
    width: 54px;
  }

  &__preview {
    margin-top: 16px;
  }
}
</style>
