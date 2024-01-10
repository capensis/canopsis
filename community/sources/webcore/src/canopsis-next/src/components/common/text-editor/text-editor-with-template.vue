<template>
  <v-layout column>
    <v-select
      :value="form.template"
      :items="templatesWithCustom"
      :label="$t('common.template')"
      return-object
      @input="updateTemplate"
    />
    <text-editor-field
      v-validate="rules"
      :value="form.text"
      :label="label"
      :error-messages="errors.collect('text')"
      :variables="variables"
      :dark="$system.dark"
      name="text"
      @input="updateText"
    />
  </v-layout>
</template>

<script>
import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { formMixin } from '@/mixins/form';

import TextEditorField from './text-editor.vue';

export default {
  inject: ['$validator', '$system'],
  components: { TextEditorField },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    templates: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      required: false,
    },
    rules: {
      type: Object,
      required: false,
    },
  },
  computed: {
    templatesWithCustom() {
      return [
        { value: CUSTOM_WIDGET_TEMPLATE, text: this.$t('common.custom'), content: '' },

        ...this.templates.map(template => ({
          ...template,

          value: template._id,
          text: template.title,
        })),
      ];
    },
  },
  methods: {
    updateText(text) {
      if (this.form.template !== CUSTOM_WIDGET_TEMPLATE && text !== this.form.text) {
        this.updateModel({
          text,
          template: CUSTOM_WIDGET_TEMPLATE,
        });

        return;
      }

      this.updateField('text', text);
    },

    updateTemplate({ value, content }) {
      if (value === this.form.template) {
        return;
      }

      this.updateModel({
        template: value,
        text: content,
      });
    },
  },
};
</script>
