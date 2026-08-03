<template>
  <c-layout
    v-if="!loaded"
    class="my-2"
    justify-center
  >
    <v-progress-circular
      color="primary"
      indeterminate
    />
  </c-layout>
  <c-layout v-else gap="4" column>
    <c-layout class="storage-settings__top-actions" gap="2" justify="end">
      <storage-settings-history-message-btn
        :history="history.entity_disabled"
        color="primary"
        variant="outlined"
        @click="archiveDisabledEntities"
      >
        {{ $t('storageSetting.entity.archiveDisabled') }}
      </storage-settings-history-message-btn>

      <storage-settings-history-message-btn
        :history="history.entity_unlinked"
        color="primary"
        variant="outlined"
        @click="archiveUnlinkedEntities"
      >
        {{ $t('storageSetting.entityUnlinked.archiveUnlinked') }}
      </storage-settings-history-message-btn>

      <storage-settings-history-message-btn
        :history="history.entity_cleaned"
        color="error"
        variant="outlined"
        @click="cleanArchivedEntities"
      >
        {{ $t('storageSetting.entityArchived.cleanArchive') }}
      </storage-settings-history-message-btn>
    </c-layout>

    <v-row justify="center">
      <v-col md="10" cols="12">
        <v-form @submit.prevent="submit">
          <storage-settings-form
            :history="history"
            @archive:disabled="archiveDisabledEntities"
            @archive:unlinked="archiveUnlinkedEntities"
            @clean:archive="cleanArchivedEntities"
          />
          <v-divider class="mt-3" />
        </v-form>

        <c-layout class="sticky-bottom-buttons mt-4" gap="2" justify="end">
          <v-btn
            :disabled="!hasChanges"
            color="primary"
            variant="outlined"
            @click="cancel"
          >
            {{ $t('common.cancel') }}
          </v-btn>
          <v-btn
            :loading="submitting"
            :disabled="submitting || !hasChanges"
            color="primary"
            @click="submit"
          >
            {{ $t('common.submit') }}
          </v-btn>
        </c-layout>
      </v-col>
    </v-row>
  </c-layout>
</template>

<script>
import { isEqual, cloneDeep } from 'lodash';
import { ref, onMounted, computed } from 'vue';

import { MODALS, TIME_UNITS } from '@/constants';

import { dataStorageSettingsToForm } from '@/helpers/entities/data-storage/form';

import { useForm } from '@/hooks/form';
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
  components: {
    StorageSettingsForm,
    StorageSettingsHistoryMessageBtn,
  },
  setup() {
    let dataStorageSettings = null;

    const form = useForm();
    const originalForm = ref(null);
    const history = ref({});
    const loaded = ref(false);

    const { t, te } = useI18n();
    const modals = useModals();
    const popups = usePopups();

    const {
      fetchDataStorageWithoutStore,
      updateDataStorage,
    } = useDataStorage();

    const {
      archiveDisabledEntitiesData,
      archiveUnlinkedEntitiesData,
      cleanArchivedEntitiesData,
    } = useEntity();
    const { setFormErrors } = useValidationFormErrors(form.values);

    const hasChanges = computed(() => !isEqual(form.values.value, originalForm.value));

    const getConfirmationPhraseMessages = key => ['title', 'text', 'phraseText', 'phrase']
      .reduce((acc, field) => {
        const messageKey = `${key}.${field}`;

        if (te(messageKey)) {
          acc[field] = t(messageKey);
        }

        return acc;
      }, {});

    /**
     * Fetches the latest history data from the data storage settings
     */
    const fetchHistory = async () => {
      const { history: newHistory } = await fetchDataStorageWithoutStore();

      history.value = newHistory || {};
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
      ...getConfirmationPhraseMessages('modals.confirmationPhrase.archiveUnlinkedEntities'),
      ...getConfirmationPhraseMessages('modals.confirmationPhrase.archiveDisabledEntities'),

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
      ...getConfirmationPhraseMessages('modals.confirmationPhrase.archiveUnlinkedEntities'),

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
      ...getConfirmationPhraseMessages('modals.confirmationPhrase.cleanStorage'),

      action: () => cleanArchivedEntitiesData(),
    });

    /**
     * Resets the form to its original state from data storage settings
     */
    const resetForm = async () => {
      form.reset(dataStorageSettingsToForm(dataStorageSettings.config));

      /**
       * We need to save the original form to compare it with the new one
       */
      originalForm.value = cloneDeep(form.values.value);

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

    const { submit, submitting } = useSubmittableForm({
      method: form.handleSubmit(() => modals.show({
        name: MODALS.confirmationPhrase,
        config: {
          ...getConfirmationPhraseMessages('modals.confirmationPhrase.updateStorageSettings'),

          action: async () => {
            try {
              dataStorageSettings = await updateDataStorage({ data: form.values.value });

              popups.success({ text: t('success.default') });

              resetForm();
            } catch (err) {
              setFormErrors(err);
            }
          },
        },
      })),
    });

    onMounted(async () => {
      dataStorageSettings = await fetchDataStorageWithoutStore();

      resetForm();
      loaded.value = true;
    });

    return {
      loaded,
      history,
      hasChanges,
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
