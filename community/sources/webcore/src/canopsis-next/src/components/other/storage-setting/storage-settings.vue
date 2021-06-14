<template lang="pug">
  v-layout.my-2(v-if="!form", justify-center)
    v-progress-circular(indeterminate, color="primary")
  v-flex(v-else, offset-xs1, md10)
    v-form(@submit.prevent="submit")
      storage-settings-form(v-model="form", :history="history")
      v-layout.mt-3(row, justify-end)
        v-btn.primary.mr-0(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { formToDataStorageSettings, dataStorageSettingsToForm } from '@/helpers/forms/data-storage';

import { submittableMixin } from '@/mixins/submittable';
import { validationErrorsMixin } from '@/mixins/form/validation-errors';
import { entitiesDataStorageSettingsMixin } from '@/mixins/entities/data-storage';

import StorageSettingsForm from '@/components/other/storage-setting/form/storage-settings-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { StorageSettingsForm },
  mixins: [
    submittableMixin(),
    validationErrorsMixin(),
    entitiesDataStorageSettingsMixin,
  ],
  data() {
    return {
      form: null,
      history: null,
    };
  },
  async mounted() {
    const dataStorageSettings = await this.fetchDataStorageSettingsWithoutStore();

    this.form = dataStorageSettingsToForm(dataStorageSettings.config);
    this.history = dataStorageSettings.history;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          await this.updateDataStorageSettings({ data: formToDataStorageSettings(this.form) });

          this.$popups.success({ text: this.$t('success.default') });
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>
