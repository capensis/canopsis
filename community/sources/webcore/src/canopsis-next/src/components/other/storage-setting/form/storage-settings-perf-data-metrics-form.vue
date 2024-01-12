<template>
  <c-information-block
    :title="$t('storageSetting.perfDataMetrics.title')"
    :help-text="$t('storageSetting.perfDataMetrics.titleHelp')"
    help-icon-color="info"
  >
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.perfDataMetrics.deleteAfter')"
      :help-text="$t('storageSetting.perfDataMetrics.deleteAfterHelpText')"
      :name="perfDataMetricsDeleteAfterFieldName"
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
  computed: {
    perfDataMetricsDeleteAfterFieldName() {
      return 'perf_data_metrics.delete_after';
    },
  },
};
</script>
