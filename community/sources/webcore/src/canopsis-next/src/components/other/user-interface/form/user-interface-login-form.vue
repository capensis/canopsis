<template>
  <c-form-block>
    <c-form-block-row :label="t('userInterface.footer')">
      <text-editor-field
        v-model="footerField"
        :label="t('userInterface.footer')"
        :config="textEditorConfig"
        :variables="variables"
        name="footer"
        public
      />
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.description')">
      <text-editor-field
        v-model="loginPageDescriptionField"
        :label="t('userInterface.description')"
        :config="textEditorConfig"
        :variables="variables"
        name="login_page_description"
        public
      />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import { computed } from 'vue';

import { objectToVariables } from '@/helpers/variables';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  components: { TextEditorField },
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { templateVars } = useTemplateVars();
    const { footerField, loginPageDescriptionField } = useFormModels({
      footerField: 'footer',
      loginPageDescriptionField: 'login_page_description',
    });

    const variables = computed(() => objectToVariables({ env: templateVars.value }));
    const textEditorConfig = computed(() => ({ disabled: props.disabled }));

    return {
      t,
      variables,
      textEditorConfig,
      footerField,
      loginPageDescriptionField,
    };
  },
};
</script>
