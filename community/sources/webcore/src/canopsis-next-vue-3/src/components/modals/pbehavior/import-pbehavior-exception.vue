<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ t('modals.importPbehaviorException.title') }}</span>
      </template>

      <template #text>
        <pbehavior-exception-import-form />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
:disabled="submitting"
:loading="submitting"
color="primary"
type="submit">
          {{ t('common.import') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { pbehaviorExceptionImportToForm } from '@/helpers/entities/pbehavior/exception/form';

import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import PbehaviorExceptionImportForm from '@/components/other/pbehavior/exceptions/form/pbehavior-exception-import-form.vue';

export default {
  components: { PbehaviorExceptionImportForm },

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

    const form = useForm({
      initialValues: pbehaviorExceptionImportToForm(config.value.exceptionImport),
    });
    const { values, handleSubmit, submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(values.value);

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
      submitting,
      closeModal,
      submit,
    };
  },
};
</script>
