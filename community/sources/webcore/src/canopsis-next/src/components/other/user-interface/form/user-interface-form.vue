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
      :label="$t('userInterface.infoPopupTimeout')"
    )
    c-duration-field(
      v-field="form.popup_timeout.error",
      :label="$t('userInterface.errorPopupTimeout')"
    )
    v-layout(row)
      c-language-field(
        v-field="form.language",
        :label="$t('userInterface.language')"
      )
    v-layout(row)
      v-text-field(
        v-field.number="form.max_matched_items",
        v-validate="'numeric|min_value:1'",
        :label="$t('userInterface.maxMatchedItems')",
        :error-messages="errors.collect('max_matched_items')",
        type="number",
        name="max_matched_items"
      )
        v-tooltip(slot="append", left)
          v-icon(slot="activator") help
          div(v-html="$t('userInterface.tooltips.maxMatchedItems')")
    v-layout(row)
      v-text-field(
        v-field.number="form.check_count_request_timeout",
        v-validate="'numeric|min_value:1'",
        :label="$t('userInterface.checkCountRequestTimeout')",
        :error-messages="errors.collect('check_count_request_timeout')",
        type="number",
        name="check_count_request_timeout"
      )
        v-tooltip(slot="append", left)
          v-icon(slot="activator") help
          div(v-html="$t('userInterface.tooltips.checkCountRequestTimeout')")
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
          :config="textEditorConfig"
        )
    v-layout.mt-3(row)
      v-flex
        text-editor-field(
          v-field="form.login_page_description",
          :label="$t('userInterface.description')",
          :config="textEditorConfig"
        )
    v-layout.mt-3(row)
      v-flex
        span.theme--light.v-label.file-selector__label {{ $t('userInterface.logo') }}
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
