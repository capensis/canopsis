<template>
  <v-layout class="gap-3" column>
    <template-testing-test-field
      name="test"
      required
    />
    <template-testing-test-variables-test-data
      v-field="form"
      :type="type"
    />
    <template-testing-test-variables-template-item
      v-for="field in fields"
      :key="field.key"
      :template="field.value"
      :title="$t(field.textKey)"
      :textarea="field.textarea"
      @input="field.updateField"
    />
  </v-layout>
</template>

<script>
import TemplateTestingTestField from './form/fields/template-testing-test-field.vue';
import TemplateTestingTestVariablesTestData from './template-testing-test-variables-test-data.vue';
import TemplateTestingTestVariablesTemplateItem from './template-testing-test-variables-template-item.vue';

export default {
  components: {
    TemplateTestingTestField,
    TemplateTestingTestVariablesTestData,
    TemplateTestingTestVariablesTemplateItem,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Array,
      default: () => [],
    },
    fields: {
      type: Array,
      default: () => [],
    },
    type: {
      type: Number,
      required: false,
    },
  },
};
</script>
