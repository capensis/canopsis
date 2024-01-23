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
            <v-tab
              v-if="isProVersion"
              :href="`#${$constants.PARAMETERS_TABS.notificationSettings}`"
            >
              {{ $t('parameters.tabs.notificationsSettings') }}
            </v-tab>
            <v-tab
              v-if="hasReadAnyWidgetTemplateAccess"
              :href="`#${$constants.PARAMETERS_TABS.widgetTemplates}`"
            >
              {{ $t('parameters.tabs.widgetTemplates') }}
            </v-tab>
            <v-tab
              v-if="hasReadAnyIconAccess"
              :href="`#${$constants.PARAMETERS_TABS.icons}`"
            >
              {{ $tc('common.icon', 2) }}
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
              <v-tab-item
                v-if="isProVersion"
                :value="$constants.PARAMETERS_TABS.notificationSettings"
              >
                <v-card-text>
                  <notifications-settings />
                </v-card-text>
              </v-tab-item>
              <v-tab-item
                v-if="hasReadAnyWidgetTemplateAccess"
                :value="$constants.PARAMETERS_TABS.widgetTemplates"
              >
                <v-card-text>
                  <widget-templates />
                </v-card-text>
              </v-tab-item>
              <v-tab-item
                v-if="hasReadAnyIconAccess"
                :value="$constants.PARAMETERS_TABS.icons"
              >
                <v-card-text>
                  <icons />
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
        <span>{{ fabData.text }}</span>
      </c-fab-btn>
    </v-fade-transition>
  </v-container>
</template>

<script>
import { MODALS, PARAMETERS_TABS } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesWidgetTemplatesMixin } from '@/mixins/entities/widget-template';
import { entitiesIconMixin } from '@/mixins/entities/icon';
import { permissionsTechnicalParametersMixin } from '@/mixins/permissions/technical/parameters';
import { permissionsTechnicalWidgetTemplateMixin } from '@/mixins/permissions/technical/widget-templates';
import { permissionsTechnicalIconMixin } from '@/mixins/permissions/technical/icon';

import UserInterface from '@/components/other/user-interface/user-interface.vue';
import ViewsImportExport from '@/components/other/view/views-import-export.vue';
import NotificationsSettings from '@/components/other/notification/notifications-settings.vue';
import WidgetTemplates from '@/components/other/widget-template/widget-templates.vue';
import Icons from '@/components/other/icons/icons.vue';

export default {
  components: {
    UserInterface,
    ViewsImportExport,
    NotificationsSettings,
    WidgetTemplates,
    Icons,
  },
  mixins: [
    entitiesInfoMixin,
    entitiesIconMixin,
    entitiesWidgetTemplatesMixin,
    permissionsTechnicalParametersMixin,
    permissionsTechnicalWidgetTemplateMixin,
    permissionsTechnicalIconMixin,
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
          text: this.$t('modals.createWidgetTemplate.create.title'),
          on: {
            refresh: this.fetchWidgetTemplatesListWithPreviousParams,
            create: this.showSelectWidgetTemplateTypeModal,
          },
        },
        [PARAMETERS_TABS.icons]: {
          text: this.$t('modals.createIcon.create.title'),
          on: {
            refresh: this.fetchIconsListWithPreviousParams,
            create: this.showCreateIconModal,
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

            this.fabData?.on?.refresh();
          },
        },
      });
    },

    showCreateIconModal() {
      this.$modals.show({
        name: MODALS.createIcon,
        config: {
          title: this.$t('modals.createIcon.create.title'),
          action: async (newIcon) => {
            await this.createIcon({ data: newIcon });

            this.$popups.success({ text: this.$t('modals.createIcon.create.success') });
            this.fabData?.on?.refresh();
          },
        },
      });
    },
  },
};
</script>
