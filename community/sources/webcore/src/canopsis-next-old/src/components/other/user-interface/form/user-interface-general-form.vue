<template>
  <c-form-block>
    <c-form-block-row
      :label="$t('userInterface.appTitle')"
      indented
    >
      <text-editor-field
        v-field="form.app_title"
        :disabled="disabled"
        :label="$t('userInterface.appTitle')"
        :variables="variables"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('userInterface.language')">
      <c-language-field
        v-field="form.language"
        :label="$t('userInterface.language')"
        :disabled="disabled"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('userInterface.defaultTheme')">
      <c-theme-field
        v-field="form.default_color_theme"
        :label="$t('userInterface.defaultTheme')"
        :disabled="disabled"
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('common.timezone')">
      <c-timezone-field
        v-field="form.timezone"
        disabled
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('userInterface.logo')" indented>
      <file-selector
        ref="fileSelectorElement"
        :max-file-size="maxFileSize"
        :disabled="disabled"
        accept="image/*"
        name="logo"
        with-files-list
        @change="changeLogoFile"
      />
    </c-form-block-row>

    <c-form-block-row
      :label="$t('userInterface.versionDescriptionTooltip')"
      indented
    >
      <text-editor-field
        v-field="form.version_description"
        :label="$t('userInterface.versionDescriptionTooltip')"
        :config="textEditorConfig"
        :variables="versionDescriptionVariables"
        public
      />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import { computed, ref } from 'vue';

import { MAX_ICON_SIZE_IN_KB } from '@/constants';

import { objectToVariables, variableTemplatePreparer } from '@/helpers/variables';

import { useModelField } from '@/hooks/form/model-field';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  inject: ['$validator'],
  components: {
    FileSelector,
    TextEditorField,
  },
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
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);
    const { templateVars } = useTemplateVars();

    const maxFileSize = MAX_ICON_SIZE_IN_KB;
    const fileSelectorElement = ref(null);

    const variables = computed(() => objectToVariables({ env: templateVars.value }));

    const versionDescriptionVariables = computed(() => ([
      ...variables.value,
      ...[
        'edition',
        'versionUpdated',
        'serialName',
      ].map(text => ({ text, value: variableTemplatePreparer(text) })),
    ]));

    const textEditorConfig = computed(() => ({ disabled: props.disabled }));

    const changeLogoFile = ([file] = []) => updateField('logo', file);

    return {
      maxFileSize,
      fileSelectorElement,
      variables,
      versionDescriptionVariables,
      textEditorConfig,

      changeLogoFile,
    };
  },
};
</script>
