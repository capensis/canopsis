<template>
  <v-container>
    <v-layout wrap>
      <v-flex xs12>
        <v-card class="ma-2">
          <v-tabs
            v-model="activeTab"
            slider-color="primary"
            transition="slide-y"
            centered
          >
            <v-tab :href="`#${$constants.PARAMETERS_TABS.parameters}`">
              {{ $t('parameters.tabs.parameters') }}
            </v-tab>
            <v-tab :href="`#${$constants.PARAMETERS_TABS.viewExportImport}`">
              {{ $t('parameters.tabs.importExportViews') }}
            </v-tab>
            <v-tab :href="`#${$constants.PARAMETERS_TABS.stateSettings}`">
              {{ $t('parameters.tabs.stateSettings') }}
            </v-tab>
            <template v-if="isProVersion">
              <v-tab :href="`#${$constants.PARAMETERS_TABS.notificationSettings}`">
                {{ $t('parameters.tabs.notificationsSettings') }}
              </v-tab>
              <v-tab :href="`#${$constants.PARAMETERS_TABS.storageSettings}`">
                {{ $t('parameters.tabs.storageSettings') }}
              </v-tab>
            </template>
            <v-tab
              v-if="hasReadAnyWidgetTemplateAccess"
              :href="`#${$constants.PARAMETERS_TABS.widgetTemplates}`"
            >
              {{ $t('parameters.tabs.widgetTemplates') }}
            </v-tab>
            <v-tabs-items v-model="activeTab">
              <v-tab-item :value="$constants.PARAMETERS_TABS.parameters">
                <v-card-text>
                  <user-interface :disabled="!hasUpdateParametersAccess" />
                </v-card-text>
              </v-tab-item>
              <v-tab-item :value="$constants.PARAMETERS_TABS.viewExportImport">
                <views-import-export />
              </v-tab-item>
              <v-tab-item :value="$constants.PARAMETERS_TABS.stateSettings">
                <v-card-text>
                  <state-settings />
                </v-card-text>
              </v-tab-item>
              <template v-if="isProVersion">
                <v-tab-item :value="$constants.PARAMETERS_TABS.notificationSettings">
                  <v-card-text>
                    <notifications-settings />
                  </v-card-text>
                </v-tab-item>
                <v-tab-item :value="$constants.PARAMETERS_TABS.storageSettings">
                  <v-card-text>
                    <storage-settings />
                  </v-card-text>
                </v-tab-item>
              </template>
              <v-tab-item
                v-if="hasReadAnyWidgetTemplateAccess"
                :value="$constants.PARAMETERS_TABS.widgetTemplates"
              >
                <v-card-text>
                  <widget-templates />
                </v-card-text>
              </v-tab-item>
            </v-tabs-items>
          </v-tabs>
        </v-card>
      </v-flex>
    </v-layout>
    <v-fade-transition>
      <c-fab-btn
        v-if="fabData"
        v-on="fabData.on"
      >
        <span>{{ fabData.label }}</span>
      </c-fab-btn>
    </v-fade-transition>
  </v-container>
</template>

<script>
import { MODALS, PARAMETERS_TABS } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesStateSettingMixin } from '@/mixins/entities/state-setting';
import { entitiesWidgetTemplatesMixin } from '@/mixins/entities/widget-template';
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
    entitiesStateSettingMixin,
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
    fabData() {
      return {
        [PARAMETERS_TABS.widgetTemplates]: {
          label: this.$t('modals.createWidgetTemplate.create.title'),
          on: {
            refresh: this.fetchWidgetTemplatesListWithPreviousParams,
            create: this.showSelectWidgetTemplateTypeModal,
          },
        },
        [PARAMETERS_TABS.stateSettings]: {
          label: this.$t('modals.createStateSetting.create.title'),
          on: {
            refresh: this.fetchStateSettingsListWithPreviousParams,
            create: this.showCreateStateSettingModal,
          },
        },
      }[this.activeTab];
    },
  },
  methods: {
    showSelectWidgetTemplateTypeModal() {
      this.$modals.show({
        name: MODALS.selectWidgetTemplateType,
        config: {
          title: this.$t('modals.createWidgetTemplate.create.title'),
          action: async (newWidgetTemplate) => {
            await this.createWidgetTemplate({ data: newWidgetTemplate });

            this.fabData?.on.refresh();
          },
        },
      });
    },

    showCreateStateSettingModal() {
      this.$modals.show({
        name: MODALS.createStateSetting,
        config: {
          title: this.$t('modals.createStateSetting.create.title'),
          action: async (newStateSetting) => {
            await this.createStateSetting({ data: newStateSetting });

            this.$popups.success({ text: this.$t('modals.createStateSetting.create.success') });
            this.fabData?.on.refresh();
          },
        },
      });
    },
  },
};
</script>
