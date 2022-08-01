<template lang="pug">
  v-container
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(fixed-tabs, slider-color="primary")
            v-tab {{ $t('parameters.tabs.parameters') }}
            v-tab-item
              v-card-text
                user-interface(:disabled="!hasUpdateParametersAccess")
            v-tab {{ $t('parameters.tabs.importExportViews') }}
            v-tab-item
              views-import-export
            v-tab {{ $t('parameters.tabs.stateSettings') }}
            v-tab-item(lazy)
              v-card-text
                state-settings
            template(v-if="isProVersion")
              v-tab {{ $t('parameters.tabs.notificationsSettings') }}
              v-tab-item(lazy)
                v-card-text
                  notifications-settings
              v-tab {{ $t('parameters.tabs.storageSettings') }}
              v-tab-item(lazy)
                v-card-text
                  storage-settings
</template>

<script>
import { entitiesInfoMixin } from '@/mixins/entities/info';
import { permissionsTechnicalParametersMixin } from '@/mixins/permissions/technical/parameters';

import UserInterface from '@/components/other/user-interface/user-interface.vue';
import ViewsImportExport from '@/components/other/view/views-import-export.vue';
import StateSettings from '@/components/other/state-setting/state-settings.vue';
import NotificationsSettings from '@/components/other/notification/notifications-settings.vue';
import StorageSettings from '@/components/other/storage-setting/storage-settings.vue';

export default {
  components: {
    UserInterface,
    ViewsImportExport,
    StateSettings,
    NotificationsSettings,
    StorageSettings,
  },
  mixins: [entitiesInfoMixin, permissionsTechnicalParametersMixin],
};
</script>
