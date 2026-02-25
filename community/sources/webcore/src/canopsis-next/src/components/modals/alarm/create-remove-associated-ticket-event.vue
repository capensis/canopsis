<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.removeAssociatedTicketEvent.title') }}</span>
      </template>
      <template #text="">
        <v-layout class="gap-4" column>
          <alarm-general-table v-if="items.length" :items="items" item-class="itemClass" />
          <c-collapse-panel class="c-alternative-bg-panel" expanded disabled>
            <template #header>
              <span class="font-weight-medium text-uppercase">
                {{ $t('modals.removeAssociatedTicketEvent.title') }}
              </span>
            </template>
            <remove-associated-ticket-event-form
              v-model="form"
              :items="items"
            />
          </c-collapse-panel>
        </v-layout>
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="isDisabled"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS } from '@/constants';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import RemoveAssociatedTicketEventForm from '@/components/widgets/alarm/forms/remove-associated-ticket-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to remove associated ticket from alarm
 */
export default {
  name: MODALS.removeAssociatedTicketEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, RemoveAssociatedTicketEventForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const form = ref({
      ticket: null,
      comment: '',
    });

    const { config, close } = useInnerModal(props);

    /**
     * Checks if an alarm item should be visually disabled (greyed out).
     * Item is disabled when a ticket is selected but the alarm does not have that ticket associated.
     *
     * @param {Object} item - Alarm item with v.tickets array
     * @param {string|null} ticketValue - Selected ticket value to match
     * @return {boolean} True if item should be disabled
     */
    const isDisabledItem = (item, ticketValue) => (
      ticketValue && !item.v.tickets.find(ticket => ticket.ticket === ticketValue)
    );

    const items = computed(() => (config.value.items ?? []).map(item => ({
      ...item,

      itemClass: isDisabledItem(item, form.value.ticket) ? 'grey--text' : '',
    })));

    const { submitting, isDisabled, submit } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action(form.value);
          close();
        }
      },
    });

    useFormConfirmableCloseModal({
      form,
      submit,
      close,
    });

    return {
      form,
      items,
      submitting,
      isDisabled,
      submit,
    };
  },
};
</script>
