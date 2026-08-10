<template>
  <c-layout column>
    <external-data-table-database-field
      :value="type"
      :disabled="!isNew"
      required
      @update:value="type = $event"
    />
    <c-name-field
      v-model="name"
      :regex="nameRegex"
      :tooltip="t('externalData.tableNameTooltip')"
      :disabled="fromConfig"
      required
    />
    <c-text-field v-model="description" :label="t('common.description')" />
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ExternalDataTableDatabaseField from './fields/exterrnal-data-table-database-field.vue';

export default {
  components: { ExternalDataTableDatabaseField },

  props: {
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
    const { t } = useI18n();

    const { type, name, description } = useFormModels(['type', 'name', 'description']);

    const nameRegex = /^[a-z_][\w_]*$/i;

    return {
      t,
      type,
      name,
      description,
      nameRegex,
    };
  },
};
</script>
