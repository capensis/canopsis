<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createAssociateTicketEvent.title') }}</span>
      </template>

      <template #text>
        <c-layout gap="3" column>
          <alarm-general-table v-if="config.items" :items="config.items" />
          <v-divider />
          <c-checkbox
            v-if="isAllComponentAlarms"
            v-model="ticketResources"
            :label="t('alarm.associateTicketResources')"
            color="primary"
          />
          <associate-ticket-event-form :form="values" class="mt-3" @update:form="updateForm" />
          <c-description-field
            v-model="comment"
            :label="t('common.comment')"
            :max-length="255"
            name="comment"
          />
          <c-alert v-if="itemsWithoutAck.length" type="info">
            {{ alertMessage }}
          </c-alert>
        </c-layout>
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="submitting"
          color="primary"
          type="submit"
        >
          {{ t('common.saveChanges') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import {
  eventToAssociateTicketForm,
  formToAssociateTicketEvent,
} from '@/helpers/entities/associate-ticket/event/form';
import { isEntityComponentType } from '@/helpers/entities/entity/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import AssociateTicketEventForm from '@/components/other/declare-ticket/form/associate-ticket-event-form.vue';
import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

/**
 * Modal to associate a ticket to an alarm.
 */
export default {
  name: MODALS.createAssociateTicketEvent,

  components: { AssociateTicketEventForm, AlarmGeneralTable },

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

    const form = useForm({
      initialValues: eventToAssociateTicketForm(config.value.associateTicket),
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const close = () => modals.hide({ id: props.modal.id });
    const updateForm = (nextForm) => {
      values.value = nextForm;
    };

    const ticketResources = form.model('ticket_resources');
    const comment = form.model('comment');

    const itemsWithoutAck = computed(
      () => (config.value.ignoreAck ? [] : (config.value.items ?? []).filter(item => !item.v.ack)),
    );
    const alertMessage = computed(() => {
      const { length: count } = itemsWithoutAck.value;

      return t('declareTicket.noAckItems', count, { count });
    });
    const isAllComponentAlarms = computed(
      () => (config.value.items ?? []).every(({ entity }) => isEntityComponentType(entity.type)),
    );

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToAssociateTicketEvent(values.value));

      close();
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close,
    });

    return {
      t,
      config,
      values,
      ticketResources,
      comment,
      itemsWithoutAck,
      alertMessage,
      isAllComponentAlarms,
      submitting,
      submit,
      close,
      updateForm,
    };
  },
};
</script>
