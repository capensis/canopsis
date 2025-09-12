<template>
  <v-card class="test-data-card">
    <v-card-text>
      <v-layout column>
        <span class="text-subtitle-2">{{ $t('templateTesting.testData') }}</span>
        <component
          v-for="field in fields"
          :is="field.is"
          :key="field.key"
          v-bind="field.bind"
          v-on="field.on"
        />
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { isUndefined } from 'lodash';
import { computed } from 'vue';

import {
  TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES,
  TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES_TO_DATA_TYPE,
  USER_PERMISSIONS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useCRUDPermissions } from '@/hooks/auth';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import TemplateTestingTestDataField from './form/fields/template-testing-data-field.vue';

const TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES_TO_COPONENTS = {
  [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event]: 'template-testing-test-data-field',
  [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response]: 'template-testing-test-data-field',
  [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.alarm]: 'c-alarm-field',
  [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.user]: 'c-user-field',
  [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.entity]: 'c-entity-field',
};

export default {
  components: { TemplateTestingTestDataField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Array,
      default: () => [],
    },
    type: {
      type: Number,
      required: false,
    },
    name: {
      type: String,
      default: 'testdata',
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const { hasReadAccess: hasReadAccessForAnyUser } = useCRUDPermissions(USER_PERMISSIONS.user);
    const { updateFieldInArrayItem } = useArrayModelField(props, emit);

    /**
     * Gets the appropriate label for a test data item based on its type and index
     * @param {Object} item - The test data item
     * @param {string} item.type - The type of the test data item
     * @param {number} [item.index] - Optional index for response-type items (0-based)
     * @param {string} item.key - The unique key identifier for the item
     * @returns {string} The translated label for the item
     */
    const getLabel = (item) => {
      if (isUndefined(item.index)) {
        return t(`templateTesting.testDataLabels.${item.type}`);
      }

      return t(`templateTesting.testDataResponseLabels.${props.type}`, { index: item.index + 1 });
    };

    const hasFilledSomeRequiredField = computed(() => props.form.some(item => item.someRequired && item.value));

    const fields = computed(() => props.form.map((item, index) => ({
      is: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES_TO_COPONENTS[item.type],
      key: item.key,
      bind: {
        value: item.value,
        label: getLabel(item),
        params: item.params,
        name: `${props.name}.${item.type === TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response ? `responses.${index}` : item.type}`,
        required: item.required || (item.someRequired && !hasFilledSomeRequiredField.value),
        disabled: item.type === TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.user && !hasReadAccessForAnyUser.value,
        type: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES_TO_DATA_TYPE[item.type],
      },
      on: {
        input: value => updateFieldInArrayItem(index, 'value', value),
      },
    })));

    return {
      fields,
    };
  },
};
</script>

<style lang="scss" scoped>
.theme--light .test-data-card {
  background-color: var(--v-application-background-darken1);
}

.theme--dark .test-data-card {
  background-color: var(--v-application-background-lighten2);
}
</style>
