<template>
  <v-layout column>
    <external-data-table-database-field
      v-field="form.type"
      :disabled="!isNew"
      required
    />
    <c-name-field
      v-field="form.name"
      :regex="nameRegex"
      :tooltip="$t('externalData.tableNameTooltip')"
      :disabled="fromConfig"
      required
    />
    <v-text-field
      v-field="form.description"
      :label="$t('common.description')"
    />
  </v-layout>
</template>

<script>
import { onMounted, onBeforeUnmount } from 'vue';

import ExternalDataTableDatabaseField from './fields/exterrnal-data-table-database-field.vue';

export default {
  components: { ExternalDataTableDatabaseField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    fromConfig: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const nameRegex = /^[a-z_][\w_]*$/i;

    const attachRegexRule = () => {};
    const detachRegexRule = () => {};

    onMounted(attachRegexRule);
    onBeforeUnmount(detachRegexRule);

    return {
      nameRegex,
    };
  },
};
</script>
