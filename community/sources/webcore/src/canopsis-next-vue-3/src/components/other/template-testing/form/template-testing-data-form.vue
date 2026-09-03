<template>
  <c-layout gap="2" column>
    <c-name-field
      v-model="nameField"
      :label="t('common.name')"
      name="name"
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('common.description')">
        <c-description-field
          v-model="descriptionField"
          :label="t('common.description')"
          :max-length="500"
          name="description"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.type')">
        <template-testing-data-type-field
          :value="typeField"
          :label="t('common.type')"
          :disabled="!isNew"
          name="type"
          required
          @update:value="typeField = $event"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('templateTesting.testData')">
        <c-json-field
          ref="jsonFieldElement"
          :value="bodyField"
          :label="t('common.value')"
          name="value"
          rows="11"
          validate-on="button"
          @update:value="bodyField = $event"
        >
          <template v-if="isEventType" #append>
            <c-layout justify="end">
              <v-btn color="primary" variant="outlined" @click="showSetPreFilledTemplateModal">
                {{ t('templateTesting.usePreFilledTemplate') }}
              </v-btn>
            </c-layout>
          </template>
        </c-json-field>
      </c-form-block-row>

      <c-form-block-row v-if="!isEventType" :label="t('common.header', 2)">
        <request-headers-field
          :headers="headersField"
          name="headers"
          @update:headers="headersField = $event"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed, ref } from 'vue';

import {
  MODALS,
  TEMPLATE_TESTING_DATA_EVENT_PRE_FILLED_TEMPLATE,
  TEMPLATE_TESTING_DATA_TYPES,
} from '@/constants';

import { useFormModels, useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import RequestHeadersField from '@/components/forms/request/fields/request-headers-field.vue';

import TemplateTestingDataTypeField from './fields/template-testing-data-type-field.vue';

export default {
  components: {
    RequestHeadersField,
    TemplateTestingDataTypeField,
  },
  props: {
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const form = useFormContext();
    const jsonFieldElement = ref(null);

    const { nameField, descriptionField, typeField, bodyField, headersField } = useFormModels({
      nameField: 'name',
      descriptionField: 'description',
      typeField: 'type',
      bodyField: 'body',
      headersField: 'headers',
    });

    const isEventType = computed(() => typeField.value === TEMPLATE_TESTING_DATA_TYPES.event);

    const setPreFilledTemplate = () => {
      form.setValue('body', TEMPLATE_TESTING_DATA_EVENT_PRE_FILLED_TEMPLATE);
      jsonFieldElement.value?.reset?.();
    };

    const showSetPreFilledTemplateModal = () => {
      const bodyToCompare = bodyField.value?.trim?.();

      if (!bodyToCompare || bodyToCompare === '{}') {
        setPreFilledTemplate();

        return;
      }

      modals.show({
        name: MODALS.confirmation,
        config: {
          title: t('templateTesting.usePreFilledTemplate'),
          text: t('templateTesting.usePreFilledTemplateWarning'),
          action: setPreFilledTemplate,
        },
      });
    };

    return {
      t,
      jsonFieldElement,
      nameField,
      descriptionField,
      typeField,
      bodyField,
      headersField,
      isEventType,
      showSetPreFilledTemplateModal,
    };
  },
};
</script>
