<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createRrule.title') }}</span>
      </template>
      <template #text="">
        <v-layout
          class="gap-2"
          column
        >
          <recurrence-rule-form
            v-model="form"
            :with-exdate-type="config.withExdateType"
          />
        </v-layout>
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.saveChanges') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import {
  recurrenceRuleModalConfigToForm,
  formToReccurenceRuleModalConfig,
} from '@/helpers/entities/shared/recurrence-rule/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import RecurrenceRuleForm from '@/components/forms/recurrence-rule/recurrence-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRecurrenceRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    RecurrenceRuleForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const form = ref(recurrenceRuleModalConfigToForm(config.value));

    const { submit, submitting, isDisabled } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToReccurenceRuleModalConfig(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      close,
      submit,
      submitting,
      isDisabled,
    };
  },
};
</script>
