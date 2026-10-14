<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createDeclareTicketEvent.title') }}</span>
      </template>

      <template #text>
        <c-enabled-field
          v-if="config.items.length > 1"
          v-model="singleMode"
          :label="t('declareTicket.oneTicketForAlarms')"
        />
        <declare-ticket-events-form
          :form="values"
          :alarms="config.items"
          :tickets-by-alarms="config.ticketsByAlarms"
          :alarms-by-tickets="config.alarmsByTickets"
          :hide-ticket-resource="!isAllComponentAlarms"
          :single-mode="singleMode"
          @update:form="updateForm"
        />
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
          {{ t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { keyBy } from 'lodash';
import { computed, ref } from 'vue';

import { MODALS } from '@/constants';

import { isSuccessStepTicketType } from '@/helpers/entities/alarm/step/entity';
import {
  alarmsToDeclareTicketEventForm,
  formToDeclareTicketEvents,
} from '@/helpers/entities/declare-ticket/event/form';
import { isEntityComponentType } from '@/helpers/entities/entity/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import DeclareTicketEventsForm from '@/components/other/declare-ticket/form/declare-ticket-events-form.vue';

/**
 * Modal to declare a ticket.
 */
export default {
  name: MODALS.createDeclareTicketEvent,

  components: { DeclareTicketEventsForm },

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
    const singleMode = ref(false);

    const form = useForm({
      initialValues: alarmsToDeclareTicketEventForm(config.value.alarmsByTickets),
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const close = () => modals.hide({ id: props.modal.id });
    const updateForm = (nextForm) => {
      values.value = nextForm;
    };

    const isAllComponentAlarms = computed(
      () => (config.value.items ?? []).every(({ entity }) => isEntityComponentType(entity.type)),
    );

    const isSomeOneSuccessTicketExist = () => {
      const alarmsById = keyBy(config.value.items, '_id');

      return Object.entries(values.value.alarms_by_tickets).some(
        ([ticketRuleId, alarmIds]) => alarmIds.some((alarmId) => {
          const alarm = alarmsById[alarmId];

          return alarm.v.tickets?.some(
            ticket => ticketRuleId === ticket.ticket_rule_id && isSuccessStepTicketType(ticket._t),
          );
        }),
      );
    };

    const callActionAndHideModal = async () => {
      const preparedForm = formToDeclareTicketEvents(values.value, singleMode.value);

      await config.value.action?.(preparedForm, singleMode.value);

      close();
    };

    const submit = handleSubmit(async () => {
      if (isSomeOneSuccessTicketExist()) {
        const isSingleAlarm = config.value.items.length === 1;

        modals.show({
          name: MODALS.confirmation,
          config: {
            title: isSingleAlarm
              ? t('modals.confirmationCreateNewTicketForAlarm.title')
              : t('modals.confirmationCreateNewTicketForAlarms.title'),
            text: isSingleAlarm
              ? t('modals.confirmationCreateNewTicketForAlarm.text')
              : t('modals.confirmationCreateNewTicketForAlarms.text'),
            action: callActionAndHideModal,
          },
        });
      } else {
        await callActionAndHideModal();
      }
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
      singleMode,
      isAllComponentAlarms,
      submitting,
      submit,
      close,
      updateForm,
    };
  },
};
</script>
