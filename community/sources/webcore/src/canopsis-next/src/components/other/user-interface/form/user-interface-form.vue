<template lang="pug">
  v-layout(column)
    v-layout(row)
      v-text-field(
        v-field="form.app_title",
        :disabled="disabled",
        :label="$t('userInterface.appTitle')"
      )
    c-duration-field(
      v-field="form.popup_timeout.info",
      :label="$t('userInterface.infoPopupTimeout')",
      name="popup_timeout.info"
    )
    c-duration-field(
      v-field="form.popup_timeout.error",
      :label="$t('userInterface.errorPopupTimeout')",
      name="popup_timeout.error"
    )
    v-layout(row)
      c-language-field(
        v-field="form.language",
        :label="$t('userInterface.language')"
      )
    v-layout(row)
      c-number-field(
        v-field="form.max_matched_items",
        :label="$t('userInterface.maxMatchedItems')",
        :min="1",
        name="max_matched_items"
      )
        template(#append="")
          c-help-icon(
            :text="$t('userInterface.tooltips.maxMatchedItems')",
            color="grey darken-1",
            icon="help",
            left
          )
    v-layout(row)
      c-number-field(
        v-field="form.check_count_request_timeout",
        :label="$t('userInterface.checkCountRequestTimeout')",
        :min="1",
        name="check_count_request_timeout"
      )
        template(#append="")
          c-help-icon(
            :text="$t('userInterface.tooltips.checkCountRequestTimeout')",
            color="grey darken-1",
            icon="help",
            left
          )
    v-layout(row)
      c-timezone-field(v-field="form.timezone", disabled)
    v-layout(row)
      v-switch(
        v-field="form.allow_change_severity_to_info",
        :label="$t('userInterface.allowChangeSeverityToInfo')",
        color="primary"
      )
    v-layout(row)
      v-flex
        text-editor-field(
          v-field="form.footer",
          :label="$t('userInterface.footer')",
          :config="textEditorConfig",
          public
        )
    v-layout.mt-3(row)
      v-flex
        text-editor-field(
          v-field="form.login_page_description",
          :label="$t('userInterface.description')",
          :config="textEditorConfig",
          public
        )
    v-layout.mt-3(row)
      v-flex
        span.v-label.file-selector__label {{ $t('userInterface.logo') }}
        v-layout(row)
          file-selector.mt-1(
            ref="fileSelector",
            v-validate="`image|size:${$config.MAX_LOGO_SIZE_IN_KB}`",
            :error-messages="errors.collect('logo')",
            :disabled="disabled",
            accept="image/*",
            name="logo",
            with-files-list,
            @change="changeLogoFile"
          )
</template>

<script>
import { formMixin } from '@/mixins/form';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  inject: ['$validator'],
  components: {
    FileSelector,
    TextEditorField,
  },
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    textEditorConfig() {
      return { disabled: this.disabled };
    },
  },
  methods: {
    async changeLogoFile([file] = []) {
      this.updateField('logo', file);
    },

    reset() {
      this.$refs.fileSelector.clear();
    },
  },
};
</script>

<style lang="scss" scoped>
.file-selector {
  &__label {
    font-size: .85em;
    display: block;
  }
}
</style>
