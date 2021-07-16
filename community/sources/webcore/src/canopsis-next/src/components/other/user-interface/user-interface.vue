<template lang="pug">
  v-form(@submit.prevent="submit")
    v-layout(row)
      v-flex.text-xs-center
        div.title {{ $t('parameters.userInterfaceForm.title') }}
    v-layout(row)
      v-text-field(
        v-model="form.app_title",
        :disabled="disabled",
        :label="$t('parameters.userInterfaceForm.fields.appTitle')"
      )
    popup-timeout-field(
      v-model="form.popup_timeout.info",
      :label="$t('parameters.userInterfaceForm.fields.infoPopupTimeout')"
    )
    popup-timeout-field(
      v-model="form.popup_timeout.error",
      :label="$t('parameters.userInterfaceForm.fields.errorPopupTimeout')"
    )
    v-layout(row)
      v-select(
        v-model="form.language",
        :items="languages",
        :label="$t('parameters.userInterfaceForm.fields.language')"
      )
    v-layout(row)
      v-text-field(
        v-model.number="form.max_matched_items",
        v-validate="'numeric|min_value:1'",
        :label="$t('parameters.userInterfaceForm.fields.maxMatchedItems')",
        :error-messages="errors.collect('max_matched_items')",
        type="number",
        name="max_matched_items"
      )
        v-tooltip(slot="append", left)
          v-icon(slot="activator") help
          div(v-html="$t('parameters.userInterfaceForm.tooltips.maxMatchedItems')")
    v-layout(row)
      v-text-field(
        v-model.number="form.check_count_request_timeout",
        v-validate="'numeric|min_value:1'",
        :label="$t('parameters.userInterfaceForm.fields.checkCountRequestTimeout')",
        :error-messages="errors.collect('check_count_request_timeout')",
        type="number",
        name="check_count_request_timeout"
      )
        v-tooltip(slot="append", left)
          v-icon(slot="activator") help
          div(v-html="$t('parameters.userInterfaceForm.tooltips.checkCountRequestTimeout')")
    v-layout(row)
      timezone-field(v-model="form.timezone", disabled)
    v-layout(row)
      v-switch(
        v-model="form.allow_change_severity_to_info",
        :label="$t('parameters.userInterfaceForm.fields.allowChangeSeverityToInfo')",
        color="primary"
      )
    v-layout(row)
      v-flex
        text-editor-field(
          v-model="form.footer",
          :label="$t('parameters.userInterfaceForm.fields.footer')",
          :config="textEditorConfig"
        )
    v-layout.mt-3(row)
      v-flex
        text-editor-field(
          v-model="form.login_page_description",
          :label="$t('parameters.userInterfaceForm.fields.description')",
          :config="textEditorConfig"
        )
    v-layout.mt-3(row)
      v-flex
        span.theme--light.v-label.file-selector__label {{ $t('parameters.userInterfaceForm.fields.logo') }}
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
    template(v-if="!disabled")
      v-divider.mt-3
      v-layout.mt-3(row, justify-end)
        v-btn(
          flat,
          @click="reset"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="submitting",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { getFileDataUrlContent } from '@/helpers/file/file-select';
import { formToUserInterface, userInterfaceToForm } from '@/helpers/forms/user-interface';

import entitiesInfoMixin from '@/mixins/entities/info';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import PopupTimeoutField from '@/components/forms/fields/popup-timeout.vue';
import TimezoneField from '@/components/forms/fields/timezone-field.vue';
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    TimezoneField,
    PopupTimeoutField,
    FileSelector,
    TextEditorField,
  },
  mixins: [entitiesInfoMixin],
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      submitting: false,
      logoFile: null,
      form: userInterfaceToForm(),
    };
  },
  computed: {
    languages() {
      return Object.keys(this.$i18n.messages);
    },

    textEditorConfig() {
      return { disabled: this.disabled };
    },
  },
  async mounted() {
    await this.fetchAllInfos();

    this.initForm();
  },
  methods: {
    initForm() {
      /**
       * TODO: update this preparing
       */
      const userInterface = {
        app_title: this.appTitle,
        language: this.language,
        footer: this.footer,
        login_page_description: this.description,
        popup_timeout: this.popupTimeout,
        allow_change_severity_to_info: this.allowChangeSeverityToInfo,
        timezone: this.timezone,
        max_matched_items: this.maxMatchedItems,
        check_count_request_timeout: this.checkCountRequestTimeout,
      };

      this.form = userInterfaceToForm(userInterface);
    },

    changeLogoFile(files = []) {
      const [file] = files;

      this.logoFile = file || null;
    },

    async submit() {
      try {
        this.submitting = true;

        const isValid = await this.$validator.validateAll();

        if (isValid) {
          const data = formToUserInterface(this.form);

          if (this.logoFile) {
            data.logo = await getFileDataUrlContent(this.logoFile);
          }

          await this.updateUserInterface({ data });
          await this.fetchAllInfos();

          this.setTitle();

          this.$popups.success({ text: this.$t('success.default') });
          this.reset();
        }
      } catch (err) {
        console.warn(err);
      } finally {
        this.submitting = false;
      }
    },

    reset() {
      this.logoFile = null;
      this.initForm();
      this.$refs.fileSelector.clear();
    },

    fetchAllInfos() {
      return Promise.all([
        this.fetchAppInfos(),
        this.fetchLoginInfos(),
      ]);
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
