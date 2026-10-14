<template>
  <c-layout gap="2">
    <div class="w-50">
      <external-data-table-database-field :value="form.type" disabled />
    </div>
    <div class="w-50">
      <c-text-field :model-value="form.name" :label="t('common.name')" disabled />
    </div>
  </c-layout>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

import ExternalDataTableDatabaseField from '@/components/other/external-data-table/form/fields/exterrnal-data-table-database-field.vue';

export default {
  components: { ExternalDataTableDatabaseField },

  props: {
    form: {
      type: Object,
      required: true,
    },
  },

  setup() {
    const { t } = useI18n();

    return {
      t,
    };
  },
};
</script>
