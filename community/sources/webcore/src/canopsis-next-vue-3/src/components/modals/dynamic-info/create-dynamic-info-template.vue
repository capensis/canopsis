<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close fill-height>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <dynamic-info-template-form />
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
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import DynamicInfoTemplateForm from '@/components/other/dynamic-info/form/dynamic-info-template-form.vue';
import { formToTemplate, templateToForm } from '@/helpers/entities/dynamic-info/template/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  components: { DynamicInfoTemplateForm },

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
    const isNew = computed(() => !config.value.template);
    const title = computed(
      () => config.value.title || t('modals.createDynamicInfoTemplate.create.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));
    const form = useForm({
      initialValues: templateToForm(config.value.template),
    });
    const namesField = form.array('names');

    const closeModal = () => modals.hide({ id: props.modal.id });
    const submit = form.handleSubmit(async () => {
      if (!namesField.value.length) {
        return;
      }

      await config.value.action?.(formToTemplate(form.values.value));

      closeModal();
    });

    useFormConfirmableCloseModal({
      submit,

      form: form.values,
      close: closeModal,
    });

    return {
      t,
      config,
      submitting: form.submitting,
      title,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
