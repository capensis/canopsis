<template>
  <c-information-block :title="$t('storageSetting.perfDataMetrics.title')">
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.perfDataMetrics.deleteAfter')"
      :suffix="$t('common.after')"
      name="perf_data_metrics.delete_after"
      switcher
      hide-value-on-false
    />
  </c-information-block>
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
