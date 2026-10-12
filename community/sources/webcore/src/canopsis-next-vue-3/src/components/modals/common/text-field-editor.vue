<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ config.title }}</span>
      </template>

      <template #text>
        <c-text-field
          v-model="text"
          v-bind="fieldProps"
          :required="field.validationRules?.required"
          autofocus
        />

        <c-alert v-if="config.alert" type="warning">
          {{ config.alert.text }}
        </c-alert>
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
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
import { omit } from 'lodash';
import { computed } from 'vue';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
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
    const field = computed(() => config.value.field ?? { name: 'text', label: 'Text' });
    const fieldName = computed(() => field.value.name ?? 'text');
    const fieldProps = computed(() => omit(field.value, ['validationRules', 'value']));

    const form = useForm({
      initialValues: {
        [fieldName.value]: field.value.value ?? field.value.text ?? '',
      },
    });
    const { model, handleSubmit, isSubmitting } = form;
    const text = model(fieldName);

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value?.action?.(text.value);

      closeModal();
    });

    useFormConfirmableCloseModal({
      form: text,
      submit,
      close: closeModal,
    });

    return {
      t,
      text,
      field,
      fieldProps,
      submit,
      isSubmitting,
      config,
      closeModal,
    };
  },
};
</script>
