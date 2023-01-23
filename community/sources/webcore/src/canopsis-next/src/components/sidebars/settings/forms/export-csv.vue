<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.exportCsv.title') }}
    v-container
      v-layout(row)
        v-select(
          v-field="form.exportCsvSeparator",
          :items="separators",
          :label="$t('settings.exportCsv.fields.separator')"
        )
      v-layout(v-if="datetimeFormat", row)
        v-select(
          v-field="form.exportCsvDatetimeFormat",
          :items="formats",
          :label="$t('settings.exportCsv.fields.datetimeFormat')"
        )
      v-layout(column)
        h4.subheading {{ $t('settings.exportColumnNames') }}
        c-columns-with-template-field(
          v-field="form.widgetExportColumns",
          :template="form.widgetExportColumnsTemplate",
          :templates="templates",
          :templates-pending="templatesPending",
          :label="$t('settings.exportColumnNames')",
          :type="type",
          :alarm-infos="alarmInfos",
          :entity-infos="entityInfos",
          :infos-pending="infosPending",
          @update:template="updateTemplate"
        )
</template>

<script>
import {
  EXPORT_CSV_SEPARATORS,
  EXPORT_CSV_DATETIME_FORMATS,
} from '@/constants';

import { formBaseMixin } from '@/mixins/form';

import FieldColumns from '../fields/common/columns.vue';

export default {
  components: { FieldColumns },
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
    alarmInfos: {
      type: Array,
      default: () => [],
    },
    entityInfos: {
      type: Array,
      default: () => [],
    },
    infosPending: {
      type: Boolean,
      default: false,
    },
    datetimeFormat: {
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
