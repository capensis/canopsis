<template>
  <v-select
    v-field="value"
    :label="$t('availability.filterByValue')"
    :items="valueFilterMethods"
    hide-details
  />
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: String,
      required: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const valueFilterMethods = computed(
      () => Object.values(AVAILABILITY_VALUE_FILTER_METHODS).map(valueFilterMethod => ({
        value: valueFilterMethod,
        text: t(`availability.valueFilterMethods.${valueFilterMethod}`),
      })),
    );

    return {
      valueFilterMethods,
    };
  },
};
</script>
