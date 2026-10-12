<template>
  <c-layout gap="2" column>
    <c-name-field
      v-model="nameField"
      :label="t('common.testName')"
      name="name"
      required
    />

    <template-testing-test-type-field
      :value="typeField"
      :label="t('common.ruleType')"
      :disabled="!isNew"
      name="type"
      required
      @update:value="typeField = $event"
    />

    <c-name-field
      :model-value="ruleField?.name"
      :label="t('common.ruleName')"
      :disabled="!isNew"
      name="rule"
      required
    />

    <c-description-field
      v-model="descriptionField"
      :label="t('common.description')"
      :max-length="500"
      name="description"
    />

    <c-form-block v-if="testDataFields.length">
      <c-form-block-row :label="t('templateTesting.testData')">
        <c-layout gap="3" column>
          <template-testing-data-field
            v-for="field in testDataFields"
            :key="field.key"
            :value="field.value"
            :label="field.label"
            :type="field.type"
            :name="field.name"
            clearable
            @update:value="updateDataField(field.path, $event)"
          />
        </c-layout>
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import {
  TEMPLATE_TESTING_DATA_TYPES,
  TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES,
} from '@/constants';

import { setField } from '@/helpers/immutable';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import TemplateTestingDataField from './fields/template-testing-data-field.vue';
import TemplateTestingTestTypeField from './fields/template-testing-test-type-field.vue';

export default {
  components: {
    TemplateTestingDataField,
    TemplateTestingTestTypeField,
  },
  props: {
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();
    const { nameField, typeField, ruleField, descriptionField, dataField } = useFormModels({
      nameField: 'name',
      typeField: 'type',
      ruleField: 'rule',
      descriptionField: 'description',
      dataField: 'data',
    });

    const getTestDataLabel = (type, index) => {
      if (type === TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event) {
        return t(`templateTesting.testDataLabels.${type}`);
      }

      if (type === TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.ticketStatusResponse) {
        return t(`templateTesting.testDataCheckTicketStatusResponseLabels.${typeField.value}`, {
          index: index + 1,
        });
      }

      return t(`templateTesting.testDataResponseLabels.${typeField.value}`, { index: index + 1 });
    };

    const getResponseFields = (field, type) => Object.entries(dataField.value?.[field] ?? {}).map(
        ([index, value]) => ({
          key: `${field}.${index}`,
          label: getTestDataLabel(type, Number(index)),
          name: `data.${field}.${index}`,
          path: [field, index],
          type: TEMPLATE_TESTING_DATA_TYPES.response,
          value,
        }),
      );

    const testDataFields = computed(() => [
      ...(dataField.value?.event
        ? [
            {
              key: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event,
              label: getTestDataLabel(TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event),
              name: 'data.event',
              path: [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.event],
              type: TEMPLATE_TESTING_DATA_TYPES.event,
              value: dataField.value.event,
            },
          ]
        : []),
      ...(dataField.value?.response
        ? [
            {
              key: TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response,
              label: getTestDataLabel(TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response, 0),
              name: 'data.response',
              path: [TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response],
              type: TEMPLATE_TESTING_DATA_TYPES.response,
              value: dataField.value.response,
            },
          ]
        : []),
      ...getResponseFields('responses', TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.response),
      ...getResponseFields(
        'ticket_status_responses',
        TEMPLATE_TESTING_TEST_VARIABLES_ENTITY_TYPES.ticketStatusResponse,
      ),
    ]);

    const updateDataField = (path, value) => {
      dataField.value = setField(dataField.value ?? {}, path, value);
    };

    return {
      t,
      nameField,
      typeField,
      ruleField,
      descriptionField,
      testDataFields,
      updateDataField,
    };
  },
};
</script>
