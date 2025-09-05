<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      :label="$t('common.testName')"
      name="name"
      required
    />

    <template-testing-test-type-field
      v-field="form.rule_type"
      :label="$t('common.ruleType')"
      :disabled="!isNew"
      name="rule_type"
      required
    />

    <c-name-field
      v-field="form.rule_name"
      :label="$t('common.ruleName')"
      :disabled="!isNew"
      name="rule_name"
      required
    />

    <v-textarea
      v-field="form.description"
      :label="$t('common.description')"
      name="description"
    />
  </v-layout>
</template>

<script>
import TemplateTestingTestTypeField from '../partials/template-testing-test-type-field.vue';

export default {
  components: {
    TemplateTestingTestTypeField,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
