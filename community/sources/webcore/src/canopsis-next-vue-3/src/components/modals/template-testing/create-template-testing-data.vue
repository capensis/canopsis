<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="pa-0" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <div class="pa-3">
          <template-testing-data-form :is-new="isNew" />
        </div>
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
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import {
  formToTemplateTestingData,
  templateTestingDataToForm,
} from '@/helpers/entities/template-testing-data/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import TemplateTestingDataForm from '@/components/other/template-testing/form/template-testing-data-form.vue';

export default {
  components: {
    TemplateTestingDataForm,
  },

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

    const isNew = computed(() => !config.value.templateTestingData?._id);
    const title = computed(
      () => config.value.title || t('modals.createTemplateTestingData.create.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: templateTestingDataToForm(config.value.templateTestingData),
    });
    const { values, handleSubmit, isSubmitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToTemplateTestingData(values.value));

      closeModal();
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
      isNew,
      isSubmitting,
      closeModal,
      submit,
    };
  },
};
</script>
