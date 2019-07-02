<template lang="pug">
  v-form(@submit.prevent="submit")
    v-layout(row)
      v-flex.text-xs-center
        .title {{ $t('parameters.userInterfaceForm.title') }}
    v-layout(row)
      v-flex
        v-text-field(
        v-model="form.appTitle",
        :label="$t('parameters.userInterfaceForm.fields.appTitle')"
        )
    v-layout(row)
      v-flex
        v-text-field(
        v-model="form.footer",
        :label="$t('parameters.userInterfaceForm.fields.footer')"
        )
    v-layout(row)
      v-flex
        span.theme--light.v-label.file-selector__label {{ $t('parameters.userInterfaceForm.fields.logo') }}
        v-layout(row)
          file-selector.mt-1(
          ref="fileSelector",
          v-validate="`image|size:${$config.MAX_LOGO_SIZE_IN_KB}`",
          :error-messages="errors.collect('logo')",
          accept="image/*",
          name="logo",
          withFilesList,
          @change="changeLogoFile"
          )
    v-divider.mt-3
    v-layout.mt-3(row, justify-end)
      v-btn(
      flat,
      @click="reset"
      ) {{ $t('common.cancel') }}
      v-btn.primary(
      type="submit",
      :loading="submitting",
      :disabled="submitting"
      ) {{ $t('common.submit') }}
</template>

<script>
import { DEFAULT_APP_TITLE } from '@/config';

import { getFileDataUrlContent } from '@/helpers/file-select';

import FileSelector from '@/components/forms/fields/file-selector.vue';

import entitiesInfoMixin from '@/mixins/entities/info';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { FileSelector },
  mixins: [entitiesInfoMixin],
  data() {
    return {
      submitting: false,
      logoFile: null,
      form: {
        appTitle: DEFAULT_APP_TITLE,
        footer: '',
      },
    };
  },
  async mounted() {
    await this.fetchAllInfos();

    this.initForm();
  },
  methods: {
    initForm() {
      this.form = {
        appTitle: this.appTitle || DEFAULT_APP_TITLE,
        footer: this.footer,
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
          };

          if (this.logoFile) {
            data.logo = await getFileDataUrlContent(this.logoFile);
          }

          await this.updateUserInterface({ data });
          await this.fetchAllInfos();

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
    }
  }
</style>
