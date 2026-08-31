<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <c-textarea
          v-if="textarea"
          v-model="form.text"
          :label="config.label"
          :required="required"
          name="text"
          autofocus
        />

        <text-editor-field
          v-else
          v-model="form.text"
          :label="config.label"
          :sanitize-options="config.sanitizeOptions"
          :variables="variables"
          :required="required"
          name="text"
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

import { normalizeHandlebarsNbsp } from '@/helpers/handlebars/normalize';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  components: { TextEditorField },

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
    const variables = computed(() => config.value.variables);
    const textarea = computed(() => config.value.textarea);
    const required = computed(() => (
      config.value.rules?.required ||
      config.value.rules?.includes?.('required') ||
      config.value.required
    ));

    const form = useForm({
      initialValues: {
        text: config.value.text ?? '',
      },
    });
    const { values, isSubmitting } = form;
    const text = form.model('text');
    const close = () => modals.hide({ id: props.modal.id });
    const submit = form.handleSubmit(async () => {
      await config.value.action?.(normalizeHandlebarsNbsp(values.value.text));

      close();
    });

    useFormConfirmableCloseModal({ form: values, submit, close });

    return {
      config,
      title,
      variables,
      textarea,
      required,
      form: values,
      text,
      isSubmitting,
      close,
      submit,
    };
  },
};
</script>
