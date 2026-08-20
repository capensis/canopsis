<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.removeAssociatedTicketEvent.title') }}</span>
      </template>

      <template #text>
        <c-layout gap="4" column>
          <alarm-general-table v-if="items.length" :items="items" item-class="itemClass" />
          <c-collapse-panel class="c-alternative-bg-panel" expanded disabled>
            <template #header>
              <span class="font-weight-medium text-uppercase">
                {{ t('modals.removeAssociatedTicketEvent.title') }}
              </span>
            </template>
            <remove-associated-ticket-event-form
              :form="values"
              :items="config.items"
              @update:form="updateForm"
            />
          </c-collapse-panel>
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
          {{ t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import RemoveAssociatedTicketEventForm from '@/components/widgets/alarm/forms/remove-associated-ticket-event-form.vue';


/**
 * Modal to remove associated ticket from alarm
 */
export default {
  name: MODALS.removeAssociatedTicketEvent,

  components: { AlarmGeneralTable, RemoveAssociatedTicketEventForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = useForm({
      initialValues: {
        ticket: null,
        comment: '',
      },
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const submittingItem = (item, ticketValue) => (
      ticketValue && !item.v.tickets?.find?.(ticket => ticket.ticket === ticketValue)
    );

    const items = computed(() => (config.value.items ?? []).map(item => ({
      ...item,

      itemClass: submittingItem(item, values.value.ticket) ? 'text-medium-emphasis' : '',
    })));

    const updateForm = nextForm => values.value = nextForm;

    const submit = handleSubmit(async () => {
      if (config.value.action) {
        await config.value.action(values.value);
        close();
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
      items,
      submitting,
      submit,
      close,
      updateForm,
    };
  },
};
</script>
