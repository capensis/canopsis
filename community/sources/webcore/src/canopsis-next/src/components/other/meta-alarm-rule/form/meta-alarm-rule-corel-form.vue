<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field.number="config.corel_id",
      v-validate="'required'",
      :label="$t('metaAlarmRule.corelId')",
      :error-messages="errors.collect('corelId')",
      name="corelId"
    )
      v-tooltip(slot="append", left)
        v-icon(slot="activator") help
        span(v-html="$t('metaAlarmRule.corelIdHelp')")
    v-text-field(
      v-field.number="config.corel_status",
      v-validate="'required'",
      :label="$t('metaAlarmRule.corelStatus')",
      :error-messages="errors.collect('corelStatus')",
      name="corelStatus"
    )
      v-tooltip(slot="append", left)
        v-icon(slot="activator") help
        span(v-html="$t('metaAlarmRule.corelStatusHelp')")
    v-text-field(
      v-field.number="config.corel_parent",
      v-validate="'required'",
      :label="$t('metaAlarmRule.corelParent')",
      :error-messages="errors.collect('corelParent')",
      name="corelParent"
    )
    v-text-field(
      v-field.number="config.corel_child",
      v-validate="'required'",
      :label="$t('metaAlarmRule.corelChild')",
      :error-messages="errors.collect('corelChild')",
      name="corelChild"
    )
    c-number-field(
      v-field="config.threshold_count",
      :label="$t('metaAlarmRule.thresholdCount')",
      :min="0",
      name="thresholdCount",
      required
    )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'config',
    event: 'input',
  },
  props: {
    config: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
