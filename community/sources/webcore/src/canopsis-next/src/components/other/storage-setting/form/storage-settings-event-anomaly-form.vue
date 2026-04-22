<template>
  <c-information-block :title="$t('storageSetting.eventAnomaly.title')">
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.eventAnomaly.deleteHourlyData')"
      :suffix="$t('storageSetting.olderThan')"
      name="event_anomaly.delete_after"
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
