<template>
  <c-layout gap="3" column>
    <c-name-field
      v-model="nameField"
      :max-length="FORM_MAX_STRING_LENGTH"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('linkRule.type')">
        <c-radio-group :model-value="typeField" mandatory @update:model-value="updateType">
          <v-radio
            v-for="type in types"
            :key="type.value"
            :value="type.value"
            :label="type.label"
            color="primary"
          />
        </c-radio-group>
      </c-form-block-row>

      <c-form-block-row :label="$t('externalData.title')">
        <external-data-field
          v-model:external-data="externalDataField"
          :types="externalDataTypes"
          :variables="templateVars.external_data"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import {
  EXTERNAL_DATA_TYPES,
  FORM_MAX_STRING_LENGTH,
  LINK_RULE_TYPES,
  LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES,
} from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useValidationHeader } from '@/hooks/validator/validation-header';

import ExternalDataField from '@/components/forms/external-data/external-data-field.vue';

export default {
  components: { ExternalDataField },
  props: {
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const { hasAnyError } = useValidationHeader();

    const { nameField, typeField, sourceCodeField, externalDataField } = useFormModels({
      nameField: 'name',
      typeField: 'type',
      sourceCodeField: 'source_code',
      externalDataField: 'external_data',
    });

    const isAlarmType = computed(() => typeField.value === LINK_RULE_TYPES.alarm);

    const types = computed(() => Object.values(LINK_RULE_TYPES).map(type => ({
        value: type,
        label: t(`linkRule.types.${type}`),
      })),
    );

    const externalDataTypes = computed(() => [
      {
        text: t(`externalData.types.${EXTERNAL_DATA_TYPES.table}`),
        value: EXTERNAL_DATA_TYPES.table,
      },
    ]);

    const updateType = (type) => {
      typeField.value = type;
      sourceCodeField.value = LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES[type] ?? '';
    };

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      FORM_MAX_STRING_LENGTH,
      nameField,
      typeField,
      externalDataField,
      isAlarmType,
      types,
      externalDataTypes,
      updateType,
    };
  },
};
</script>
