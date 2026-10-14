<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('apiToken.create') }}</span>
      </template>

      <template #text>
        <api-token-form />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ t('common.create') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { apiTokenToForm, formToApiToken } from '@/helpers/entities/api-token/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import ApiTokenForm from '@/components/other/api-token/form/api-token-form.vue';

export default {
  components: { ApiTokenForm },

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
    const form = useForm({ initialValues: apiTokenToForm() });
    const { values, handleSubmit, submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });
    const submit = handleSubmit(async () => {
      await config.value.action?.(formToApiToken(values.value));

      closeModal();
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      submitting,
      submit,
      closeModal,
    };
  },
};
</script>
