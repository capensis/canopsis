<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      :label="$t('common.testName')"
      name="name"
      required
    />

    <template-testing-test-type-field
      :value="form.type"
      :label="$t('common.ruleType')"
      :disabled="!isNew"
      name="type"
      required
    />

    <c-name-field
      :value="form.rule?.name"
      :label="$t('common.ruleName')"
      :disabled="!isNew"
      name="rule"
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
import TemplateTestingTestTypeField from './fields/template-testing-test-type-field.vue';

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
