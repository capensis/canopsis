<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createTicketStatusJob.edit.title') }}</span>
      </template>
      <template #text="">
        <ticket-status-job-form v-model="form" />
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
          @click="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { ticketStatusJobToForm, formToTicketStatusJob } from '@/helpers/entities/ticket-status-job/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import TicketStatusJobForm from '@/components/other/ticket-status-job/ticket-status-job-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createTicketStatusJob,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    TicketStatusJobForm,
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

    const form = ref(ticketStatusJobToForm(config.value.ticketStatusJob ?? {}));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToTicketStatusJob(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      isDisabled,
      submitting,
      submit,
      close,
    };
  },
};
</script>
