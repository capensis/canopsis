<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createChangeStateEvent.title') }}</span>
      </template>
      <template #text="">
        <alarm-general-table
          v-if="config.items"
          :items="config.items"
          class="mb-4"
        />
        <c-change-state-field
          v-model="form"
          :label="$t('modals.createChangeStateEvent.fields.output')"
        />
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
          :loading="submitting"
          :disabled="isDisabled"
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

import { MODALS, ALARM_STATES } from '@/constants';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create a 'change-state' event
 */
export default {
  name: MODALS.createChangeStateEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const { items } = props.modal.config;
    const [firstItem] = items ?? [];

    const form = ref({
      comment: '',
      state: firstItem && items.length === 1 ? firstItem.v.state.val : ALARM_STATES.major,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(form.value);
        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      config,
      form,
      isDisabled,
      submitting,
      submit,
      close,
    };
  },
};
</script>
