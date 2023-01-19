<template lang="pug">
  v-container
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(v-model="activeTab", slider-color="primary", fixed-tabs)
            v-tab(:href="`#${$constants.PARAMETERS_TABS.parameters}`") {{ $t('parameters.tabs.parameters') }}
            v-tab-item(:value="$constants.PARAMETERS_TABS.parameters")
              v-card-text
                user-interface(:disabled="!hasUpdateParametersAccess")
            v-tab(:href="`#${$constants.PARAMETERS_TABS.viewExportImport}`")
              | {{ $t('parameters.tabs.importExportViews') }}
            v-tab-item(:value="$constants.PARAMETERS_TABS.viewExportImport")
              views-import-export
            v-tab(:href="`#${$constants.PARAMETERS_TABS.stateSettings}`") {{ $t('parameters.tabs.stateSettings') }}
            v-tab-item(:value="$constants.PARAMETERS_TABS.stateSettings", lazy)
              v-card-text
                state-settings
            v-tab(:href="`#${$constants.PARAMETERS_TABS.widgetTemplates}`") {{ $t('parameters.tabs.widgetTemplates') }}
            v-tab-item(:value="$constants.PARAMETERS_TABS.widgetTemplates", lazy)
              v-card-text
                widget-templates
            template(v-if="isProVersion")
              v-tab(:href="`#${$constants.PARAMETERS_TABS.notificationSettings}`")
                | {{ $t('parameters.tabs.notificationsSettings') }}
              v-tab-item(:value="$constants.PARAMETERS_TABS.notificationSettings", lazy)
                v-card-text
                  notifications-settings
              v-tab(:href="`#${$constants.PARAMETERS_TABS.storageSettings}`")
                | {{ $t('parameters.tabs.storageSettings') }}
              v-tab-item(:value="$constants.PARAMETERS_TABS.storageSettings", lazy)
                v-card-text
                  storage-settings
    v-fade-transition
      c-fab-btn(
        v-if="hasFabButton",
        :has-access="hasCreateAnyWidgetTemplateAccess",
        @refresh="refresh",
        @create="create"
      )
        span {{ $t('modals.createWidgetTemplate.create.title') }}
</template>

<script>
import { MODALS, PARAMETERS_TABS } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesWidgetTemplatesMixin } from '@/mixins/entities/view/widget/template';
import { permissionsTechnicalParametersMixin } from '@/mixins/permissions/technical/parameters';
import { permissionsTechnicalWidgetTemplateMixin } from '@/mixins/permissions/technical/widget-templates';

import UserInterface from '@/components/other/user-interface/user-interface.vue';
import ViewsImportExport from '@/components/other/view/views-import-export.vue';
import StateSettings from '@/components/other/state-setting/state-settings.vue';
import NotificationsSettings from '@/components/other/notification/notifications-settings.vue';
import StorageSettings from '@/components/other/storage-setting/storage-settings.vue';
import WidgetTemplates from '@/components/other/widget-template/widget-templates.vue';

export default {
  components: {
    UserInterface,
    ViewsImportExport,
    StateSettings,
    NotificationsSettings,
    StorageSettings,
    WidgetTemplates,
  },
  mixins: [
    entitiesInfoMixin,
    entitiesWidgetTemplatesMixin,
    permissionsTechnicalParametersMixin,
    permissionsTechnicalWidgetTemplateMixin,
  ],
  data() {
    return {
      activeTab: PARAMETERS_TABS.parameters,
    };
  },
  computed: {
    hasFabButton() {
      return this.activeTab === PARAMETERS_TABS.widgetTemplates;
    },
  },
  methods: {
    refresh() {
      this.fetchWidgetTemplatesListWithPreviousParams();
    },

    create() {
      this.$modals.show({
        name: MODALS.createWidgetTemplate,
        config: {
          action: async (newWidgetTemplate) => {
            await this.createWidgetTemplate({ data: newWidgetTemplate });

            return this.refresh();
          },
        },
      });
    },
  },
};
</script>
