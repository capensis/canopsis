<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <external-auth-token-form />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn :disabled="submitting" :loading="submitting" color="primary" type="submit">
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import {
  externalAuthTokenToForm,
  formToExternalAuthToken,
} from '@/helpers/entities/external-auth-token/form';

import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import ExternalAuthTokenForm from '@/components/other/external-auth-token/form/external-auth-token-form.vue';

export default {
  components: { ExternalAuthTokenForm },

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
    const externalAuthToken = computed(() => config.value.externalAuthToken);
    const isNew = computed(() => !externalAuthToken.value?._id);
    const title = computed(
      () =>
        config.value.title ||
        t(`modals.createExternalAuthToken.${isNew.value ? 'create' : 'edit'}.title`),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: externalAuthTokenToForm(externalAuthToken.value),
    });
    const { values, handleSubmit, submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToExternalAuthToken(values.value));

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      title,
      submitLabel,
      submitting,
      closeModal,
      submit,
    };
  },
};
</script>
