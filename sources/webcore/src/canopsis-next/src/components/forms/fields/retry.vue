<template lang="pug">
  v-layout
    v-flex
      v-text-field(
        v-field="value.count",
        v-validate="'required|numeric|min_value:0'",
        :label="$t('modals.createWebhook.fields.retryCount')",
        :min="0",
        data-vv-name="retryCount",
        type="number"
      )
    v-flex
      v-text-field(
        v-field="value.delay",
        v-validate="'required|numeric|min_value:0'",
        :label="$t('modals.createWebhook.fields.retryDelay')",
        :min="0",
        data-vv-name="retryDelay",
        type="number"
      )
</template>

<script>
import { RETRY_DEFAULT_DELAY } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      default: () => ({
        count: 0,
        delay: RETRY_DEFAULT_DELAY,
      }),
    },
  },
};
</script>
