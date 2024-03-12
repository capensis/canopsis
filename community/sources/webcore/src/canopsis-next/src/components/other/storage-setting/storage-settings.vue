<template>
  <v-layout
    v-if="!form"
    class="my-2"
    justify-center
  >
    <v-progress-circular
      color="primary"
      indeterminate
    />
  </v-layout>
  <v-flex
    v-else
    offset-xs1
    md10
  >
    <v-form @submit.prevent="submit">
      <storage-settings-form
        v-model="form"
        :history="history"
        @archive:disabled="archiveDisabledEntities"
        @archive:unlinked="archiveUnlinkedEntities"
        @clean:archive="cleanArchivedEntities"
      />
      <v-divider class="mt-3" />
      <v-layout
        class="mt-3"
        justify-end
      >
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary mr-0"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </v-layout>
    </v-form>
  </v-flex>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import { submittableMixinCreator } from '@/mixins/submittable';
import { entitiesDataStorageSettingsMixin } from '@/mixins/entities/data-storage';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';

import StorageSettingsForm from '@/components/other/storage-setting/form/storage-settings-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { StorageSettingsForm },
  mixins: [
    entitiesDataStorageSettingsMixin,
    entitiesContextEntityMixin,
    submittableMixinCreator(),
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
    showConfirmationPhraseModal({ action, ...config }) {
      this.$modals.show({
        name: MODALS.confirmationPhrase,
        config: {
          title: this.$t('modals.confirmationPhrase.cleanStorage.title'),
          text: this.$t('modals.confirmationPhrase.cleanStorage.text'),
          phraseText: this.$t('modals.confirmationPhrase.cleanStorage.phraseText'),
          phrase: this.$t('modals.confirmationPhrase.cleanStorage.phrase'),
          ...config,
          action: async () => {
            await action();

            this.$popups.success({ text: this.$t('success.default') });

            const { history } = await this.fetchDataStorageSettingsWithoutStore();

            this.history = history;
          },
        },
      });
    },

    archiveDisabledEntities() {
      this.$modals.show({
        name: MODALS.archiveDisabledEntities,
        config: {
          action: data => this.archiveDisabledEntitiesData({ data }),
        },
      });
    },

    archiveUnlinkedEntities() {
      this.showConfirmationPhraseModal({
        action: () => this.archiveUnlinkedEntitiesData({ data: this.form.entity_unlinked }),
      });
    },

    cleanArchivedEntities() {
      this.showConfirmationPhraseModal({
        action: () => this.cleanArchivedEntitiesData(),
      });
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.$modals.show({
          name: MODALS.confirmationPhrase,
          config: {
            title: this.$t('modals.confirmationPhrase.updateStorageSettings.title'),
            text: this.$t('modals.confirmationPhrase.updateStorageSettings.text'),
            phraseText: this.$t('modals.confirmationPhrase.updateStorageSettings.phraseText'),
            phrase: this.$t('modals.confirmationPhrase.updateStorageSettings.phrase'),
            action: async () => {
              try {
                await this.updateDataStorageSettings({ data: this.form });

                this.$popups.success({ text: this.$t('success.default') });
              } catch (err) {
                this.setFormErrors(err);
              }
            },
          },
        });
      }
    },
  },
};
</script>
