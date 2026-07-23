<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <pbehavior-exception-form />
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
  formToPbehaviorException,
  pbehaviorExceptionToForm,
} from '@/helpers/entities/pbehavior/exception/form';

import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import PbehaviorExceptionForm from '@/components/other/pbehavior/exceptions/form/pbehavior-exception-form.vue';

export default {
  components: { PbehaviorExceptionForm },

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
    const exception = computed(() => config.value.pbehaviorException ?? config.value.exception);
    const isNew = computed(() => !exception.value?._id);
    const title = computed(
      () =>
        config.value.title ||
        t(`modals.createPbehaviorException.${isNew.value ? 'create' : 'edit'}.title`),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: pbehaviorExceptionToForm(exception.value),
    });
    const { values, handleSubmit, submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToPbehaviorException(values.value));

      await config.value.afterSubmit?.(result);

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
