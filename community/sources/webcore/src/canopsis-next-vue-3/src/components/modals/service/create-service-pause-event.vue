<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createPause.title') }}</span>
      </template>
      <template #text>
        <service-pause-event-form v-model="form" />
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
          {{ t('common.saveChanges') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS } from '@/constants';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import ServicePauseEventForm from '@/components/widgets/service-weather/service-pause-event-form.vue';

export default {
  name: MODALS.createServicePauseEvent,

  components: { ServicePauseEventForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const { values: form, handleSubmit, isSubmitting: submitting } = useForm({
      initialValues: {
        comment: '',
        reason: '',
      },
    });

    const submit = handleSubmit(async () => {
      await config.value.action?.(form.value);

      close();
    });

    useFormConfirmableCloseModal({
      form,
      submit,
      close,
    });

    return {
      t,
      form,
      submitting,
      close,
      submit,
    };
  },
};
</script>
