<template>
  <v-select
    v-field="value"
    :items="typeItems"
    :label="$t('common.actionType')"
    name="type"
    required
  />
</template>

<script>
import { computed } from 'vue';

import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: DYNAMIC_INFO_INFORMATION_TYPES.setConstantValue,
    },
  },
  setup() {
    const { t } = useI18n();

    const typeItems = computed(() => Object.values(DYNAMIC_INFO_INFORMATION_TYPES).map(type => ({
      text: t(`dynamicInfo.informationTypes.${type}`),
      value: type,
    })));

    return {
      typeItems,
    };
  },
};
</script>
