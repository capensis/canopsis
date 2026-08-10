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
          :disabled="isDisabled"
          :loading="submitting"
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
import TextEditorWithTemplateField from '@/components/common/text-editor/text-editor-with-template.vue';
import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';
import { normalizeHandlebarsNbsp } from '@/helpers/handlebars/normalize';

export default {
  components: { TextEditorWithTemplateField },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    const { text = '', template = CUSTOM_WIDGET_TEMPLATE } = this.modal.config ?? {};

    return {
      submitting: false,
      form: {
        text,
        template,
      },
    };
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },

    title() {
      return this.config.title ?? this.$t('modals.textEditor.title');
    },

    required() {
      return (
        this.config.rules?.required ||
        this.config.rules?.includes?.('required') ||
        this.config.required
      );
    },

    isDisabled() {
      return this.required && !this.form.text;
    },
  },

  methods: {
    close() {
      this.$modals.hide({ id: this.modal.id });
    },

    async submit() {
      if (this.isDisabled) {
        return;
      }

      this.submitting = true;

      try {
        if (this.config.action) {
          await this.config.action({
            ...this.form,

            text: normalizeHandlebarsNbsp(this.form.text),
          });
        }

        this.close();
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>
