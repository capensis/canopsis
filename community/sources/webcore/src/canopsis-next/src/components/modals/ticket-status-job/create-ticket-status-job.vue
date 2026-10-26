<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <ticket-status-job-form :is-new="isNew" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import {
  formToTicketStatusJob,
  ticketStatusJobToForm,
} from '@/helpers/entities/ticket-status-job/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import TicketStatusJobForm from '@/components/other/ticket-status-job/ticket-status-job-form.vue';

export default {
  components: {
    TicketStatusJobForm,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const ticketStatusJob = computed(() => config.value.ticketStatusJob);
    const isNew = computed(() => !ticketStatusJob.value?._id);
    const title = computed(
      () => config.value.title || t('modals.createTicketStatusJob.edit.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.save'));

    const form = useForm({
      initialValues: ticketStatusJobToForm(ticketStatusJob.value),
    });
    const { values, handleSubmit, isSubmitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });
    const submit = handleSubmit(async () => {
      await config.value.action?.(formToTicketStatusJob(values.value));

      closeModal();
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      title,
      values,
      isNew,
      isSubmitting,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
