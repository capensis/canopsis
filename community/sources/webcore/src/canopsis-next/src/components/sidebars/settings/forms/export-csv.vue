<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile {{ $t('settings.exportCsv.title') }}
    v-container
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
      h4.subheading.my-4 {{ $t('settings.exportColumnNames') }}
      c-columns-field.subheading(v-field="form.widgetExportColumns")
</template>

<script>
import {
  EXPORT_CSV_SEPARATORS,
  EXPORT_CSV_DATETIME_FORMATS,
} from '@/constants';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
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
};
</script>
