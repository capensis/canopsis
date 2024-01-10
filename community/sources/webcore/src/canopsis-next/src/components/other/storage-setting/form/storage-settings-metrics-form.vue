<template>
  <c-information-block
    :title="$t('storageSetting.metrics.title')"
    :help-text="$t('storageSetting.metrics.titleHelp')"
    help-icon-color="info"
  >
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.metrics.deleteAfter')"
      :help-text="$t('storageSetting.metrics.deleteAfterHelpText')"
      :name="metricsDeleteAfterFieldName"
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
    metricsDeleteAfterFieldName() {
      return 'metrics.delete_after';
    },
  },
};
</script>
