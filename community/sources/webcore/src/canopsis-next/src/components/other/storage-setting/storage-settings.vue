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
  <v-layout v-else class="gap-4" column>
    <v-layout class="gap-2" justify-end>
      <storage-settings-history-message-btn
        :history="history.entity_disabled"
        color="primary"
        outlined
        @click="archiveDisabledEntities"
      >
        {{ $t('storageSetting.entity.archiveDisabled') }}
      </storage-settings-history-message-btn>

      <storage-settings-history-message-btn
        :history="history.entity_unlinked"
        color="primary"
        outlined
        @click="archiveUnlinkedEntities"
      >
        {{ $t('storageSetting.entityUnlinked.archiveUnlinked') }}
      </storage-settings-history-message-btn>

      <storage-settings-history-message-btn
        :history="history.entity_cleaned"
        color="error"
        outlined
        @click="cleanArchivedEntities"
      >
        {{ $t('storageSetting.entityArchived.cleanArchive') }}
      </storage-settings-history-message-btn>
    </v-layout>
    <v-flex
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
      </v-form>
    </v-flex>
    <v-layout class="sticky-bottom-buttons gap-2" justify-end>
      <v-btn
        :disabled="!hasChanges"
        color="primary"
        outlined
        @click="cancel"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
        :loading="submitting"
        :disabled="isDisabled || !hasChanges"
        color="primary"
        @click="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </v-layout>
  </v-layout>
</template>

<script>
import { isEqual, cloneDeep } from 'lodash';
import { ref, onMounted, computed } from 'vue';

import { MODALS, VALIDATION_DELAY, TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useDataStorage } from '@/hooks/store/modules/data-storage';
import { useEntity } from '@/hooks/store/modules/entity';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';
import { useSubmittableForm } from '@/hooks/submittable-form';

import StorageSettingsForm from '@/components/other/storage-setting/form/storage-settings-form.vue';
import StorageSettingsHistoryMessageBtn from '@/components/other/storage-setting/partials/storage-settings-history-message-btn.vue';
import StorageSettingsArchiveEntityUnlinkedForm from '@/components/other/storage-setting/form/storage-settings-archive-entity-unlinked-form.vue';
import StorageSettingsArchiveEntityDisabledForm from '@/components/other/storage-setting/form/storage-settings-archive-entity-disabled-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    StorageSettingsForm,
    StorageSettingsHistoryMessageBtn,
  },
  setup() {
    let dataStorageSettings = null;

    const form = ref(null);
    const originalForm = ref(null);
    const history = ref(null);

    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();

    const {
      fetchDataStorageSettingsWithoutStore,
      updateDataStorageSettings,
    } = useDataStorage();

    const {
      archiveDisabledEntitiesData,
      archiveUnlinkedEntitiesData,
      cleanArchivedEntitiesData,
    } = useEntity();
    const { setFormErrors } = useValidationFormErrors(form);

    const hasChanges = computed(() => !isEqual(form.value, originalForm.value));

    /**
     * Fetches the latest history data from the data storage settings
     */
    const fetchHistory = async () => {
      const { history: newHistory } = await fetchDataStorageSettingsWithoutStore();

      history.value = newHistory;
    };

    /**
     * Shows a confirmation phrase modal with custom action and configuration
     *
     * @param {Object} params - The modal parameters
     * @param {Function} params.action - The action to execute after confirmation
     * @param {Object} params.config - Additional modal configuration
     */
    const showConfirmationPhraseModal = ({ action, ...config }) => modals.show({
      name: MODALS.confirmationPhrase,
      config: {
        ...config,
        action: async (data) => {
          await action(data);

          popups.success({ text: t('success.default') });

          return fetchHistory();
        },
      },
    });

    /**
     * Initiates the archiving process for disabled entities
     */
    const archiveDisabledEntities = () => showConfirmationPhraseModal({
      ...t('modals.confirmationPhrase.archiveUnlinkedEntities'),
      ...t('modals.confirmationPhrase.archiveDisabledEntities'),

      additionalForm: {
        component: StorageSettingsArchiveEntityDisabledForm,
        form: { with_dependencies: false },
      },
      action: data => archiveDisabledEntitiesData({ data }),
    });

    /**
     * Initiates the archiving process for unlinked entities with duration configuration
     */
    const archiveUnlinkedEntities = () => showConfirmationPhraseModal({
      ...t('modals.confirmationPhrase.archiveUnlinkedEntities'),

      additionalForm: {
        component: StorageSettingsArchiveEntityUnlinkedForm,
        form: { value: 60, unit: TIME_UNITS.day },
      },
      action: duration => archiveUnlinkedEntitiesData({ data: { archive_before: duration } }),
    });

    /**
     * Initiates the cleaning process for archived entities
     */
    const cleanArchivedEntities = () => showConfirmationPhraseModal({
      ...t('modals.confirmationPhrase.cleanStorage'),

      action: () => cleanArchivedEntitiesData(),
    });

    /**
     * Resets the form to its original state from data storage settings
     */
    const resetForm = () => {
      form.value = dataStorageSettingsToForm(dataStorageSettings.config);

      /**
       * We need to save the original form to compare it with the new one
       */
      originalForm.value = cloneDeep(form.value);

      history.value = dataStorageSettings.history ?? {};
    };

    /**
     * Shows a confirmation modal before canceling changes and resetting the form
     */
    const cancel = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: resetForm,
      },
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: () => modals.show({
        name: MODALS.confirmationPhrase,
        config: {
          action: async () => {
            try {
              await updateDataStorageSettings({ data: form.value });

              popups.success({ text: t('success.default') });
            } catch (err) {
              setFormErrors(err);
            }
          },
        },
      }),
    });

    onMounted(async () => {
      dataStorageSettings = await fetchDataStorageSettingsWithoutStore();

      resetForm();
    });

    return {
      form,
      history,
      hasChanges,
      isDisabled,
      submitting,
      archiveDisabledEntities,
      archiveUnlinkedEntities,
      cleanArchivedEntities,
      submit,
      cancel,
    };
  },
};
</script>
