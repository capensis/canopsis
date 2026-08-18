<template>
  <widget-settings-item :title="$t('settings.exportCsv.title')">
    <c-csv-separator-field
      :value="form.exportCsvSeparator"
      @update:value="updateField('exportCsvSeparator', $event)"
    />
    <c-select
      v-if="datetimeFormat"
      :model-value="form.exportCsvDatetimeFormat"
      :items="formats"
      :label="$t('settings.exportCsv.fields.datetimeFormat')"
      item-title="text"
      item-value="value"
      @update:model-value="updateField('exportCsvDatetimeFormat', $event)"
    />
    <h4 class="text-subtitle-1 my-4">
      {{ $t('settings.exportColumnNames') }}
    </h4>
    <c-columns-field
      v-if="withoutTemplate"
      :model-value="form.widgetExportColumns"
      :items="items"
      :label="$t('settings.exportColumnNames')"
      :type="type"
      :name="name"
      :with-instructions="withInstructions"
      :with-template="withTemplate"
      :with-simple-template="withSimpleTemplate"
      :variables="variables"
      :optional-infos-attributes="optionalInfosAttributes"
      :without-infos-attributes="withoutInfosAttributes"
      :excluded-columns="excludedColumns"
      :without-custom-label="withoutCustomLabel"
      @update:model-value="updateField('widgetExportColumns', $event)"
    />
    <c-columns-with-template-field
      v-else
      :model-value="form.widgetExportColumns"
      :template="form.widgetExportColumnsTemplate"
      :templates="templates"
      :templates-pending="templatesPending"
      :type="type"
      :name="name"
      :with-instructions="withInstructions"
      :with-template="withTemplate"
      :with-simple-template="withSimpleTemplate"
      :variables="variables"
      :optional-infos-attributes="optionalInfosAttributes"
      :without-infos-attributes="withoutInfosAttributes"
      :excluded-columns="excludedColumns"
      @update:model-value="updateField('widgetExportColumns', $event)"
      @update:template="updateTemplate"
    />
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import { EXPORT_CSV_DATETIME_FORMATS } from '@/constants';

import WidgetSettingsItem from '../partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'widgetExportColumns',
    },
    type: {
      type: String,
      required: false,
    },
    templates: {
      type: Array,
      default: () => [],
    },
    templatesPending: {
      type: Boolean,
      default: false,
    },
    datetimeFormat: {
      type: Boolean,
      default: false,
    },
    withInstructions: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => [],
    },
    optionalInfosAttributes: {
      type: Boolean,
      default: false,
    },
    withTemplate: {
      type: Boolean,
      default: false,
    },
    withSimpleTemplate: {
      type: Boolean,
      default: false,
    },
    withoutInfosAttributes: {
      type: Boolean,
      default: false,
    },
    withoutTemplate: {
      type: Boolean,
      default: false,
    },
    withoutCustomLabel: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      required: false,
    },
    excludedColumns: {
      type: Array,
      default: undefined,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const form = computed(() => props.modelValue);
    const formats = computed(() => Object.values(EXPORT_CSV_DATETIME_FORMATS));

    const updateField = (field, value) => {
      emit('update:modelValue', { ...props.modelValue, [field]: value });
    };

    const updateTemplate = (template, columns) => {
      emit('update:modelValue', {
        ...props.modelValue,
        widgetExportColumnsTemplate: template,
        widgetExportColumns: columns,
      });
    };

    return {
      form,
      formats,
      updateField,
      updateTemplate,
    };
  },
};
</script>
