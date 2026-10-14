<template>
  <c-information-block :title="$t('storageSetting.metrics.title')">
    <c-enabled-duration-field
      v-field="form.delete_after"
      :label="$t('storageSetting.metrics.deleteAfter')"
      :suffix="$t('common.after')"
      name="metrics.delete_after"
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
