<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <text-editor-with-template-field
          v-model="form"
          :templates="config.templates"
          :variables="config.variables"
          :required="required"
          :label="config.label"
          :sanitize-options="config.sanitizeOptions"
          autofocus
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="close">
          {{ $t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { normalizeHandlebarsNbsp } from '@/helpers/handlebars/normalize';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import TextEditorWithTemplateField from '@/components/common/text-editor/text-editor-with-template.vue';

export default {
  components: { TextEditorWithTemplateField },

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
    const title = computed(() => config.value.title ?? t('modals.textEditor.title'));
    const required = computed(() => (
      config.value.rules?.required ||
      config.value.rules?.includes?.('required') ||
      config.value.required
    ));

    const formContext = useForm({
      initialValues: {
        text: config.value.text ?? '',
        template: config.value.template ?? CUSTOM_WIDGET_TEMPLATE,
      },
    });
    const { values: form, isSubmitting } = formContext;
    const close = () => modals.hide({ id: props.modal.id });
    const submit = formContext.handleSubmit(async () => {
      await config.value.action?.({
        ...form.value,

        text: normalizeHandlebarsNbsp(form.value.text),
      });

      close();
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      config,
      title,
      required,
      form,
      isSubmitting,
      close,
      submit,
    };
  },
};
</script>
