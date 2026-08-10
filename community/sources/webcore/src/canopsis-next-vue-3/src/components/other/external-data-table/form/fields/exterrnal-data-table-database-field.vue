<template>
  <c-select
    :model-value="value"
    :items="databases"
    :label="$t('common.database')"
    :disabled="disabled"
    :required="required"
    item-title="text"
    item-value="value"
    @update:model-value="updateModel"
  />
</template>
<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_TABLES_TYPES } from '@/constants';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  emits: ['update:value'],
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const databases = computed(() =>
      Object.values(EXTERNAL_DATA_TABLES_TYPES).map(value => ({
        value,
        text: t(`externalData.tableTypes.${value}`),
      })),
    );

    return {
      databases,
      updateModel,
    };
  },
};
</script>
