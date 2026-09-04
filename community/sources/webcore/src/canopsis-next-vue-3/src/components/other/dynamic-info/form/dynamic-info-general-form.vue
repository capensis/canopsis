<template>
  <c-layout gap="2" column>
    <c-name-field
      v-model="nameField"
      :autofocus="isDisabledIdField"
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('common.description')">
        <c-description-field v-model="descriptionField" required />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.disableDuringPeriods')">
        <c-disable-during-periods-field v-model="disableDuringPeriodsField" />
      </c-form-block-row>

      <c-form-block-row :label="$t('modals.createDynamicInfo.infosSection.title')">
        <dynamic-info-infos-form
          :variables="templateVars.value"
          :copy-variables="copyVars.value"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useValidationHeader } from '@/hooks/validator/validation-header';

import DynamicInfoInfosForm from './dynamic-info-infos-form.vue';

export default {
  components: {
    DynamicInfoInfosForm,
  },

  props: {
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    copyVars: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { hasAnyError } = useValidationHeader();
    const { nameField, descriptionField, disableDuringPeriodsField } = useFormModels({
      nameField: 'name',
      descriptionField: 'description',
      disableDuringPeriodsField: 'disable_during_periods',
    });

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs.
       */
      hasAnyError,
      nameField,
      descriptionField,
      disableDuringPeriodsField,
    };
  },
};
</script>
