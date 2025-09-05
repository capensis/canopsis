<template>
  <v-layout class="gap-3" column>
    <template-testing-test-field
      v-field="form.test"
      name="test"
      required
    />
    <template-testing-test-variables-test-data :items="items" />
    <template-testing-test-variables-template-item
      v-for="item in items"
      :key="item.key"
      :template="item.value"
      :title="$t(item.textKey)"
      :textarea="item.textarea"
      @input="item.updateField"
    />
  </v-layout>
</template>

<script>
import { isEqual, get } from 'lodash';
import { ref, watch } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { setField } from '@/helpers/immutable';
import { convertRuleToTemplateTestingTestValidateForm } from '@/helpers/entities/template-testing-validate/form';

import TemplateTestingTestField from './form/fields/template-testing-test-field.vue';
import TemplateTestingTestVariablesTestData from './template-testing-test-variables-test-data.vue';
import TemplateTestingTestVariablesTemplateItem from './template-testing-test-variables-template-item.vue';

export default {
  components: {
    TemplateTestingTestField,
    TemplateTestingTestVariablesTestData,
    TemplateTestingTestVariablesTemplateItem,
  },
  props: {
    mainForm: {
      type: Object,
      required: true,
    },
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const items = ref([]);

    watch(() => props.mainForm, (mainForm) => {
      const newItems = convertRuleToTemplateTestingTestValidateForm(mainForm, TEMPLATE_TESTING_TEST_TYPES.eventFilter);

      if (isEqual(items.value, newItems)) {
        return;
      }

      items.value = newItems.map(newItem => ({
        ...newItem,
        value: get(mainForm, newItem.key),
        updateField: value => emit('update:mainForm', setField(mainForm, newItem.key, value)),
      }));
    }, { immediate: true });

    return {
      items,
    };
  },
};
</script>
