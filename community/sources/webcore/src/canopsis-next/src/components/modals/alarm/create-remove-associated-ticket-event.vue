<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.removeAssociatedTicketEvent.title') }}</span>
      </template>
      <template #text="">
        <v-layout class="gap-4" column>
          <template v-if="items.length">
            <alarm-general-table :items="items" />
          </template>
          <remove-associated-ticket-event-form
            v-model="form"
            :items="items"
          />
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
      reason: '',
    });

    const { config, close } = useInnerModal(props);

    const items = computed(() => config.value.items ?? []);

    const { submitting, isDisabled, submit } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action({
            ticket: form.value.ticket,
            reason: form.value.reason,
          });
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
