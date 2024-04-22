<template>
  <v-layout column>
    <c-id-field
      v-field="form._id"
      :disabled="disabledIdField"
      :help-text="$t('metaAlarmRule.idHelp')"
    />
    <c-name-field
      v-field="form.name"
      required
    />
    <c-description-field
      v-field="form.output_template"
      :label="$t('metaAlarmRule.outputTemplate')"
      :help-text="$t('metaAlarmRule.outputTemplateHelp')"
    />
    <c-enabled-field
      v-field="form.auto_resolve"
      :label="$t('metaAlarmRule.autoResolve')"
    />
  </v-layout>
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
      default: () => ({}),
    },
    disabledIdField: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
