<template>
  <div>
    <v-tooltip left>
      <template #activator="{ on }">
        <div v-on="on">
          {{ version }}
        </div>
      </template>
      <c-compiled-template :template="template" :context="context" />
    </v-tooltip>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useInfo } from '@/hooks/store/modules/info';
import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();
    const {
      version,
      edition,
      serialName,
      versionUpdated,
      versionDescription,
    } = useInfo();

    const template = computed(() => (versionDescription.value || t('userInterface.defaultVersionDescription')));
    const context = computed(() => ({
      version: version.value,
      edition: edition.value,
      serialName: serialName.value,
      versionUpdated: versionUpdated.value,
    }));

    return {
      version,
      template,
      context,
    };
  },
};
</script>
