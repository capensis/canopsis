<template>
  <c-form-block>
    <c-form-block-row :label="t('userInterface.appTitle')" indented>
      <text-editor-field
        v-model="appTitleField"
        :disabled="disabled"
        :label="t('userInterface.appTitle')"
        :variables="variables"
        name="app_title"
      />
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.language')">
      <c-language-field
        v-model="languageField"
        :label="t('userInterface.language')"
        :disabled="disabled"
      />
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.defaultTheme')">
      <c-theme-field
        v-model="defaultColorThemeField"
        :label="t('userInterface.defaultTheme')"
        :disabled="disabled"
      />
    </c-form-block-row>

    <c-form-block-row :label="t('common.timezone')">
      <c-timezone-field v-model="timezoneField" disabled />
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.logo')" indented align-center>
      <file-selector
        v-model="logoField"
        ref="fileSelectorElement"
        :max-file-size="maxFileSize"
        :disabled="disabled"
        accept="image/*"
        name="logo"
        with-files-list
      >
        <template #activator="{ disabled: activatorDisabled, loading, on }">
          <v-btn :disabled="activatorDisabled" :loading="loading" color="primary" v-on="on">
            <v-icon icon="cloud_upload" />
          </v-btn>
        </template>
      </file-selector>
    </c-form-block-row>

    <c-form-block-row :label="t('userInterface.versionDescriptionTooltip')" indented>
      <text-editor-field
        v-model="versionDescriptionField"
        :label="t('userInterface.versionDescriptionTooltip')"
        :config="textEditorConfig"
        :variables="versionDescriptionVariables"
        name="version_description"
        public
      />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import { computed, ref } from 'vue';

import { MAX_ICON_SIZE_IN_KB } from '@/constants';

import { objectToVariables, variableTemplatePreparer } from '@/helpers/variables';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  components: {
    FileSelector,
    TextEditorField,
  },
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { templateVars } = useTemplateVars();
    const {
      appTitleField,
      languageField,
      defaultColorThemeField,
      timezoneField,
      logoField,
      versionDescriptionField,
    } = useFormModels({
      appTitleField: 'app_title',
      languageField: 'language',
      defaultColorThemeField: 'default_color_theme',
      timezoneField: 'timezone',
      logoField: 'logo',
      versionDescriptionField: 'version_description',
    });

    const maxFileSize = MAX_ICON_SIZE_IN_KB;
    const fileSelectorElement = ref(null);
    const variables = computed(() => objectToVariables({ env: templateVars.value }));
    const versionDescriptionVariables = computed(() => [
      ...variables.value,
      ...['edition', 'versionUpdated', 'serialName'].map((text) => ({
        text,
        value: variableTemplatePreparer(text),
      })),
    ]);
    const textEditorConfig = computed(() => ({ disabled: props.disabled }));

    return {
      t,
      maxFileSize,
      fileSelectorElement,
      variables,
      versionDescriptionVariables,
      textEditorConfig,
      appTitleField,
      languageField,
      defaultColorThemeField,
      timezoneField,
      logoField,
      versionDescriptionField,
    };
  },
};
</script>
