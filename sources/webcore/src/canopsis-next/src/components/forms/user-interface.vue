<template lang="pug">
  v-form(
    data-test="userInterfaceForm",
    @submit.prevent="submit"
  )
    v-layout(row)
      v-flex.text-xs-center
        .title {{ $t('parameters.userInterfaceForm.title') }}
    v-layout(row)
      v-flex
        v-text-field(
          data-test="appTitle",
          v-model="form.appTitle",
          :disabled="disabled",
          :label="$t('parameters.userInterfaceForm.fields.appTitle')"
        )
    popup-timeout-field(
      v-model="form.popupTimeout.info",
      :label="$t('parameters.userInterfaceForm.fields.infoPopupTimeout')"
    )
    popup-timeout-field(
      v-model="form.popupTimeout.error",
      :label="$t('parameters.userInterfaceForm.fields.errorPopupTimeout')"
    )
    v-layout(
      data-test="languageLayout",
      row
    )
      v-flex
        v-select(
          v-model="form.language",
          :items="languages",
          :label="$t('parameters.userInterfaceForm.fields.language')"
        )
    v-layout(row)
      v-flex
        timezone-field(v-model="form.timezone", disabled)
    v-layout(row)
      v-switch(
        v-model="form.allowChangeSeverityToInfo",
        :label="$t('parameters.userInterfaceForm.fields.allowChangeSeverityToInfo')",
        color="primary"
      )
    v-layout(
      data-test="footerLayout",
      row
    )
      v-flex
        text-editor-field(
          v-model="form.footer",
          :label="$t('parameters.userInterfaceForm.fields.footer')",
          :config="textEditorConfig"
        )
    v-layout.mt-3(
      data-test="descriptionLayout",
      row
    )
      v-flex
        text-editor-field(
          v-model="form.description",
          :label="$t('parameters.userInterfaceForm.fields.description')",
          :config="textEditorConfig"
        )
    v-layout.mt-3(row)
      v-flex
        span.theme--light.v-label.file-selector__label {{ $t('parameters.userInterfaceForm.fields.logo') }}
        v-layout(row)
          file-selector.mt-1(
            data-test="fileSelector",
            ref="fileSelector",
            v-validate="`image|size:${$config.MAX_LOGO_SIZE_IN_KB}`",
            :error-messages="errors.collect('logo')",
            :disabled="disabled",
            accept="image/*",
            name="logo",
            withFilesList,
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
          data-test="submitButton",
          :disabled="submitting",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { DEFAULT_APP_TITLE, DEFAULT_LOCALE } from '@/config';

import { getFileDataUrlContent } from '@/helpers/file/file-select';

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
      form: {
        appTitle: DEFAULT_APP_TITLE,
        language: DEFAULT_LOCALE,
        footer: '',
        description: '',
        timezone: '',
        popupTimeout: {},
        allowChangeSeverityToInfo: false,
      },
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
      this.form = {
        appTitle: this.appTitle || DEFAULT_APP_TITLE,
        language: this.language || DEFAULT_LOCALE,
        footer: this.footer || '',
        description: this.description || '',
        popupTimeout: this.popupTimeout ? { ...this.popupTimeout } : {},
        allowChangeSeverityToInfo: this.allowChangeSeverityToInfo || false,
        timezone: this.timezone,
      };
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
          const data = {
            app_title: this.form.appTitle,
            footer: this.form.footer,
            language: this.form.language,
            login_page_description: this.form.description,
            popup_timeout: this.form.popupTimeout,
            allow_change_severity_to_info: this.form.allowChangeSeverityToInfo,
          };

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
