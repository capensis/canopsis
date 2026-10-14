<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createSnoozeEvent.title') }}</span>
      </template>

      <template #text>
        <snooze-event-form :is-note-required="config.isNoteRequired" />
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
import { MODALS } from '@/constants';

import { snoozeToForm } from '@/helpers/entities/alarm/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import SnoozeEventForm from '@/components/widgets/alarm/forms/snooze-event-form.vue';

/**
 * Modal to put a snooze on an alarm
 */
export default {
  name: MODALS.createSnoozeEvent,

  components: { SnoozeEventForm },

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
      initialValues: snoozeToForm(),
    });
    const { handleSubmit, isSubmitting: submitting } = form;

    const submit = handleSubmit(async () => {
      await config.value?.action?.(form.values.value);

      close();
    });

    useFormConfirmableCloseModal({
      form: form.values,
      submit,
      close,
    });

    return {
      t,
      config,
      submitting,
      submit,
      close,
    };
  },
};
</script>
