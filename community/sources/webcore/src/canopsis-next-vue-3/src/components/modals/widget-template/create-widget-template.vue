<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <widget-template-form />
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="closeModal">
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

import { MODALS } from '@/constants';

import { formToWidgetTemplate, widgetTemplateToForm } from '@/helpers/entities/widget/template/form';
import { getWidgetTemplateModalTitle } from '@/helpers/entities/widget/template/modal-title';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import WidgetTemplateForm from '@/components/other/widget-template/form/widget-template-form.vue';

export default {
  name: MODALS.createWidgetTemplate,

  components: { WidgetTemplateForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t, te } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const form = useForm({
      initialValues: widgetTemplateToForm(config.value.widgetTemplate),
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const title = computed(() => {
      if (config.value.title) {
        return config.value.title;
      }

      return getWidgetTemplateModalTitle({
        t,
        te,
        type: values.value.type ?? config.value.widgetTemplate?.type,
        isEdit: !!config.value.widgetTemplate?._id,
      });
    });
    const submitLabel = computed(() =>
      t(config.value.widgetTemplate?._id ? 'common.save' : 'common.create'),
    );
    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToWidgetTemplate(values.value));

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
      submitting,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
