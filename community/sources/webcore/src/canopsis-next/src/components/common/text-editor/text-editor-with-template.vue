<template>
  <v-layout column>
    <c-widget-template-field
      :value="form.template"
      :templates="templates"
      @input="updateTemplate"
    />
    <text-editor-field
      v-validate="rules"
      :value="form.text"
      :label="label"
      :error-messages="errors.collect('text')"
      :variables="variables"
      :sanitize-options="sanitizeOptions"
      :dark="$system.dark"
      name="text"
      @input="updateText"
    />
  </v-layout>
</template>

<script>
import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

import TextEditorField from './text-editor.vue';

export default {
  inject: ['$validator', '$system'],

  components: { TextEditorField },

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
    sanitizeOptions: {
      type: Object,
      required: false,
    },
  },

  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);

    const updateTemplate = ({ value, content }) => {
      if (value === props.form.template) {
        return;
      }

      updateModel({ emit }, {
        template: value,
        text: content,
      });
    };

    const updateText = (text) => {
      if (props.form.template !== CUSTOM_WIDGET_TEMPLATE && text !== props.form.text) {
        updateModel({ emit }, {
          text,
          template: CUSTOM_WIDGET_TEMPLATE,
        });

        return;
      }

      updateField({ emit }, 'text', text);
    };

    return {
      updateTemplate,
      updateText,
    };
  },
};
</script>
