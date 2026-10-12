<template>
  <widget-settings-item :title="$t('settings.exportCsv.title')">
    <c-csv-separator-field v-field="form.exportCsvSeparator" />
    <v-select
      v-if="datetimeFormat"
      v-field="form.exportCsvDatetimeFormat"
      :items="formats"
      :label="$t('settings.exportCsv.fields.datetimeFormat')"
    />
    <v-layout column>
      <h4 class="text-subtitle-1 my-4">
        {{ $t('settings.exportColumnNames') }}
      </h4>
      <c-columns-field
        v-if="withoutTemplate"
        v-field="form.widgetExportColumns"
        :items="items"
        :label="$t('settings.exportColumnNames')"
        :type="type"
        :with-instructions="withInstructions"
        :with-template="withTemplate"
        :with-simple-template="withSimpleTemplate"
        :variables="variables"
        :optional-infos-attributes="optionalInfosAttributes"
        :without-infos-attributes="withoutInfosAttributes"
        :excluded-columns="excludedColumns"
        :without-custom-label="withoutCustomLabel"
      />
      <c-columns-with-template-field
        v-else
        v-field="form.widgetExportColumns"
        :items="items"
        :template="form.widgetExportColumnsTemplate"
        :templates="templates"
        :templates-pending="templatesPending"
        :label="$t('settings.exportColumnNames')"
        :type="type"
        :with-instructions="withInstructions"
        :variables="variables"
        :optional-infos-attributes="optionalInfosAttributes"
        :with-template="withTemplate"
        :with-simple-template="withSimpleTemplate"
        :without-infos-attributes="withoutInfosAttributes"
        :excluded-columns="excludedColumns"
        @update:template="updateTemplate"
      />
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { EXPORT_CSV_DATETIME_FORMATS } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

import WidgetSettingsItem from '../partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  mixins: [formBaseMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
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
      required: false,
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
      required: false,
    },
  },
  computed: {
    formats() {
      return Object.values(EXPORT_CSV_DATETIME_FORMATS);
    },
  },
  methods: {
    updateTemplate(template, columns) {
      this.updateModel({
        ...this.form,

        widgetExportColumnsTemplate: template,
        widgetExportColumns: columns,
      });
    },
  },
};
</script>
