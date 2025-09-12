<template>
  <v-layout class="gap-3" column>
    <v-form data-vv-scope="test-name">
      <template-testing-test-field
        :value="''"
        name="test.name"
        required
      />
    </v-form>
    <v-form data-vv-scope="test-data">
      <template-testing-test-variables-test-data
        v-field="form"
        :type="type"
      />
    </v-form>
    <template-testing-test-variables-template-item
      v-for="field in variablesFields"
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
    variablesFields: {
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
