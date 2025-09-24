<template>
  <v-layout class="gap-3" column>
    <v-form data-vv-scope="test-name">
      <template-testing-test-field
        :value="selectedTest"
        :rule-id="preparedRuleId"
        name="test.name"
        return-object
        required
        @input="updateSelectedTest"
      />
      <c-alert v-if="!ruleId" type="warning">
        {{ $t('templateTesting.isNewForm') }}
      </c-alert>
    </v-form>
    <v-form data-vv-scope="test-data">
      <template-testing-test-variables-test-data
        v-field="form"
        :type="type"
      />
    </v-form>
    <v-form data-vv-scope="test-result">
      <v-layout class="gap-3" column>
        <template-testing-test-variables-template-item
          v-for="field in variablesFields"
          :key="field.key"
          :name="field.key"
          :template="field.value"
          :title="$t(field.textKey, field.textArgs)"
          :textarea="field.textarea"
          :json="field.json"
          :result="testResult[field.key]"
          :variables="templateVars[field.templateVarsKey]"
          :last-run-value="lastRunVariables[field.formKey ?? field.key]"
          @input="field.updateField"
        />
      </v-layout>
    </v-form>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import TemplateTestingTestField from '../../form/fields/template-testing-test-field.vue';
import TemplateTestingTestVariablesTestData from '../partials/template-testing-test-variables-test-data.vue';
import TemplateTestingTestVariablesTemplateItem from '../partials/template-testing-test-variables-template-item.vue';

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
    selectedTest: {
      type: Object,
      default: () => ({}),
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    variablesFields: {
      type: Array,
      default: () => [],
    },
    type: {
      type: Number,
      required: false,
    },
    ruleId: {
      type: String,
      required: false,
    },
    testResult: {
      type: Object,
      default: () => ({}),
    },
    lastRunVariables: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const preparedRuleId = computed(() => props.ruleId ?? 'fake-id');

    /**
     * Updates the selected test by emitting the new selected test value
     * @param {Object} newSelectedTest - The new selected test object
     */
    const updateSelectedTest = newSelectedTest => emit('update:selected-test', newSelectedTest);

    return {
      preparedRuleId,

      updateSelectedTest,
    };
  },
};
</script>
