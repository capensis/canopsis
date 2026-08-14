<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createAckEvent.title') }}</span>
      </template>

      <template #text>
        <c-layout gap="2" column>
          <alarm-general-table
            v-if="config.items"
            :items="config.items"
          />
          <ack-event-form
            :is-note-required="isNoteRequired"
            :hide-ack-resources="!isAllComponentAlarms"
          />
          <c-radio-group
            v-model="actionType"
            :label="t('alarm.actionsRequired')"
            name="actionType"
            hide-details
          >
            <v-radio
              v-for="type in actionTypes"
              :key="type.value"
              :value="type.value"
              :label="type.label"
              color="primary"
            />
          </c-radio-group>
        </c-layout>
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitting"
          :loading="submitting"
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
import { computed, ref } from 'vue';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import AckEventForm from '@/components/widgets/alarm/forms/ack-event-form.vue';
import { ACK_MODAL_ACTIONS_TYPES, MODALS } from '@/constants';
import { isEntityComponentType } from '@/helpers/entities/entity/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

/**
 * Modal to create an ack event
 */
export default {
  name: MODALS.createAckEvent,

  components: { AlarmGeneralTable, AckEventForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);
    const actionType = ref(ACK_MODAL_ACTIONS_TYPES.ack);

    const form = useForm({
      initialValues: {
        comment: '',
        ack_resources: false,
      },
    });

    const isAlarmsHasDeclareTicketRules = computed(() =>
      (config.value.items ?? []).some(
        ({ assigned_declare_ticket_rules: assignedDeclareTicketRules }) =>
          assignedDeclareTicketRules?.length,
      ),
    );

    const isAllComponentAlarms = computed(() =>
      (config.value.items ?? []).every(({ entity }) => isEntityComponentType(entity.type)),
    );

    const actionTypes = computed(() => {
      const types = [
        {
          label: t('alarm.acknowledge'),
          value: ACK_MODAL_ACTIONS_TYPES.ack,
        },
        {
          label: t('alarm.acknowledgeAndAssociateTicket'),
          value: ACK_MODAL_ACTIONS_TYPES.ackAndAssociateTicket,
        },
      ];

      if (isAlarmsHasDeclareTicketRules.value) {
        types.push({
          label: t('alarm.acknowledgeAndDeclareTicket'),
          value: ACK_MODAL_ACTIONS_TYPES.ackAndDeclareTicket,
        });
      }

      return types;
    });

    const isNoteRequired = computed(() => config.value?.isNoteRequired);

    const { submitting, submit } = useSubmittableForm({
      method: form.handleSubmit(async () => {
        if (config.value.action) {
          const needDeclareTicket = actionType.value === ACK_MODAL_ACTIONS_TYPES.ackAndDeclareTicket;
          const needAssociateTicket =
            actionType.value === ACK_MODAL_ACTIONS_TYPES.ackAndAssociateTicket;

          await config.value.action(form.values.value, {
            needDeclareTicket,
            needAssociateTicket,
          });
        }

        close();
      }),
    });

    useFormConfirmableCloseModal({
      form: form.values,
      submit,
      close,
    });

    return {
      t,
      config,
      actionType,
      actionTypes,
      isAllComponentAlarms,
      isNoteRequired,
      submitting,
      submit,
      close,
    };
  },
};
</script>
