<template>
  <c-form-block>
    <c-form-block-row
      :label="$t('userInterface.footer')"
      indented
    >
      <text-editor-field
        v-field="form.footer"
        :label="$t('userInterface.footer')"
        :config="textEditorConfig"
        :variables="variables"
        public
      />
    </c-form-block-row>

    <c-form-block-row
      :label="$t('userInterface.description')"
      indented
    >
      <text-editor-field
        v-field="form.login_page_description"
        :label="$t('userInterface.description')"
        :config="textEditorConfig"
        :variables="variables"
        public
      />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import { computed } from 'vue';

import { objectToVariables } from '@/helpers/variables';

import { useTemplateVars } from '@/hooks/store/modules/template-vars';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  inject: ['$validator'],
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { templateVars } = useTemplateVars();

    const variables = computed(() => objectToVariables({ env: templateVars.value }));
    const textEditorConfig = computed(() => ({ disabled: props.disabled }));

    return {
      variables,
      textEditorConfig,
    };
  },
};
</script>
