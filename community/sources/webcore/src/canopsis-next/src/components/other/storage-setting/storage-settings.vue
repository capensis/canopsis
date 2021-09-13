<template lang="pug">
  v-layout.my-2(v-if="!form", justify-center)
    v-progress-circular(indeterminate, color="primary")
  v-flex(v-else, offset-xs1, md10)
    v-form(@submit.prevent="submit")
      storage-settings-form(
        v-model="form",
        :history="history",
        @clean-entities="cleanEntities"
      )
      v-divider.mt-3
      v-layout.mt-3(row, justify-end)
        v-btn.primary.mr-0(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { formToDataStorageSettings, dataStorageSettingsToForm } from '@/helpers/forms/data-storage';

import { submittableMixinCreator } from '@/mixins/submittable';
import { validationErrorsMixinCreator } from '@/mixins/form/validation-errors';
import { entitiesDataStorageSettingsMixin } from '@/mixins/entities/data-storage';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';

import StorageSettingsForm from '@/components/other/storage-setting/form/storage-settings-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { StorageSettingsForm },
  mixins: [
    submittableMixinCreator(),
    validationErrorsMixinCreator(),
    entitiesDataStorageSettingsMixin,
    entitiesContextEntityMixin,
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
    cleanEntities() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          title: this.$t('storageSetting.entity.confirmation.title'),
          text: this.form.entity.archive
            ? this.$t('storageSetting.entity.confirmation.archive')
            : this.$t('storageSetting.entity.confirmation.delete'),
          action: async () => {
            await this.cleanEntitiesData({ data: this.form.entity });

            this.$popups.success({ text: this.$t('success.default') });

            const { history } = await this.fetchDataStorageSettingsWithoutStore();

            this.history = history;
          },
        },
      });
    },

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
