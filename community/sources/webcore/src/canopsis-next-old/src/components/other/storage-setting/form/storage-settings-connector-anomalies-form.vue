<template>
  <c-information-block :title="$t('storageSetting.connectorAnomalies.title')">
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.connectorAnomalies.deleteHourlyData')"
      :suffix="$t('storageSetting.olderThan')"
      name="connector_anomalies.delete_after"
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
