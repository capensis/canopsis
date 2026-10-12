<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createEntityResumeEvent.title') }}</span>
      </template>

      <template #text>
        <entity-resume-form />
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

import { entityResumeFormToRequest } from '@/helpers/entities/entity/pause';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import EntityResumeForm from '@/components/forms/entity-resume-form.vue';

export default {
  name: MODALS.createEntityResumeEvent,

  components: { EntityResumeForm },

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
        comment: '',
        ...config.value.initialValues,
      },
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const submit = handleSubmit(async () => {
      await config.value.action?.(entityResumeFormToRequest(values.value));
      close();
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close,
    });

    return {
      t,
      values,
      submitting,
      submit,
      close,
    };
  },
};
</script>
