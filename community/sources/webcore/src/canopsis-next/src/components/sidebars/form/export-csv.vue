<template lang="pug">
  widget-settings-item(:title="$t('settings.exportCsv.title')")
    v-select(
      v-field="form.exportCsvSeparator",
      :items="separators",
      :label="$t('settings.exportCsv.fields.separator')"
    )
    v-select(
      v-if="datetimeFormat",
      v-field="form.exportCsvDatetimeFormat",
      :items="formats",
      :label="$t('settings.exportCsv.fields.datetimeFormat')"
    )
    v-layout(column)
      h4.subheading.my-4 {{ $t('settings.exportColumnNames') }}
      c-columns-with-template-field(
        v-field="form.widgetExportColumns",
        :template="form.widgetExportColumnsTemplate",
        :templates="templates",
        :templates-pending="templatesPending",
        :label="$t('settings.exportColumnNames')",
        :type="type",
        :with-instructions="withInstructions",
        :variables="variables",
        :optional-infos-attributes="optionalInfosAttributes",
        :with-simple-template="withSimpleTemplate",
        :without-infos-attributes="withoutInfosAttributes",
        @update:template="updateTemplate"
      )
</template>

<script>
import { EXPORT_CSV_SEPARATORS, EXPORT_CSV_DATETIME_FORMATS } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

import FieldColumns from '@/components/sidebars/form/fields/columns.vue';

import WidgetSettingsItem from '../partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem, FieldColumns },
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
      required: true,
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
    withSimpleTemplate: {
      type: Boolean,
      default: false,
    },
    withoutInfosAttributes: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    separators() {
      return Object.values(EXPORT_CSV_SEPARATORS);
    },

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
