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
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';
import { normalizeHandlebarsNbsp } from '@/helpers/handlebars/normalize';

export default {
  components: { TextEditorField },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    const text = this.modal.config?.text ?? '';

    return {
      submitting: false,
      form: {
        text,
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

    variables() {
      return this.config.variables;
    },

    textarea() {
      return this.config.textarea;
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
          await this.config.action(normalizeHandlebarsNbsp(this.form.text));
        }

        this.close();
      } finally {
        this.submitting = false;
      }
    },
  },
};
</script>
