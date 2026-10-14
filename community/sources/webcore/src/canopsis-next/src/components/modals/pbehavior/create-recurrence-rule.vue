<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createRrule.title') }}</span>
      </template>

      <template #text>
        <c-layout gap="2" column>
          <recurrence-rule-form
            :form="values"
            :with-exdate-type="config.withExdateType"
            @update:form="updateForm"
          />
        </c-layout>
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ t('common.saveChanges') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script setup>
import { computed } from 'vue';

import {
  formToReccurenceRuleModalConfig,
  recurrenceRuleModalConfigToForm,
} from '@/helpers/entities/shared/recurrence-rule/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import RecurrenceRuleForm from '@/components/forms/recurrence-rule/recurrence-rule-form.vue';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const modals = useModals();
const config = computed(() => props.modal.config ?? {});

const form = useForm({
  initialValues: recurrenceRuleModalConfigToForm(config.value),
});
const { values, handleSubmit, submitting } = form;

const closeModal = () => modals.hide({ id: props.modal.id });
const updateForm = (value) => {
  values.value = value;
};

const submit = handleSubmit(async () => {
  await config.value.action?.(formToReccurenceRuleModalConfig(values.value));

  closeModal();
});

useFormConfirmableCloseModal({
  form: values,
  submit,
  close: closeModal,
});
</script>
