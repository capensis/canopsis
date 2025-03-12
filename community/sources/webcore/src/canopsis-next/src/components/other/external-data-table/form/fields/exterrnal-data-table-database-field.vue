<template>
  <c-select-field
    v-field="value"
    :items="databases"
    :label="$t('common.database')"
    :disabled="disabled"
    :required="required"
  />
</template>
<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_TABLES_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: Number,
      default: EXTERNAL_DATA_TABLES_TYPES.mongo,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const databases = computed(() => (
      Object.values(EXTERNAL_DATA_TABLES_TYPES).map(value => ({ value, text: t(`externalData.tableTypes.${value}`) }))
    ));

    return {
      databases,
    };
  },
};
</script>
