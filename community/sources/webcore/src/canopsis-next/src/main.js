import { createPinia } from 'pinia';
import { createApp } from 'vue';

import { FEATURES_SERVICE_KEYS, MODALS } from '@/constants';
import * as constants from '@/constants';

import { featuresService } from '@/services/features';

import * as ActionsPanelComponents from '@/components/common/actions-panel';
import CCircleBadge from '@/components/common/badge/c-circle-badge.vue';
import * as BlockComponents from '@/components/common/block';
import * as ButtonComponents from '@/components/common/buttons';
import * as CalendarComponents from '@/components/common/calendar';
import * as ChartComponents from '@/components/common/chart';
import * as ChipComponents from '@/components/common/chips';
import * as DragDropComponents from '@/components/common/drag-drop';
import * as IconComponents from '@/components/common/icon';
import * as InputComponents from '@/components/common/inputs';
import * as LayoutComponents from '@/components/common/layout';
import CLayout from '@/components/common/layout/c-layout.vue';
import * as ListComponents from '@/components/common/list';
import * as ModalComponents from '@/components/common/modal';
import * as OverlayComponents from '@/components/common/overlay';
import * as PageComponents from '@/components/common/page';
import * as PaginationComponents from '@/components/common/pagination';
import * as PlaceholderComponents from '@/components/common/placeholder';
import * as PopupComponents from '@/components/common/popups';
import CProgressOverlay from '@/components/common/progress/c-progress-overlay.vue';
import * as RequestComponents from '@/components/common/request';
import * as ResponsiveListComponents from '@/components/common/responsive-list';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CCopyWrapper from '@/components/common/runtime-template/c-copy-wrapper.vue';
import CNow from '@/components/common/runtime-template/c-now.vue';
import CRequestHelper from '@/components/common/runtime-template/c-request-helper.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import * as ScrollbarComponents from '@/components/common/scrollbar';
import * as SearchComponents from '@/components/common/search';
import * as TableComponents from '@/components/common/table';
import * as TooltipComponents from '@/components/common/tooltip';
import * as TreeviewComponents from '@/components/common/treeview';
import * as VuetifyCommonComponents from '@/components/common/vuetify';
import * as FormFieldComponents from '@/components/forms/fields';
import RegisteredModals from '@/components/modals';

import App from './app.vue';
import { i18n } from './i18n';
import modalsPlugin from './plugins/modals';
import popupsPlugin from './plugins/popups';
import socketPlugin from './plugins/socket';
import veeValidatePlugin from './plugins/vee-validate';
import vuetify, { registerVuetifyComponentAliases } from './plugins/vuetify';
import { customIconsRegistry } from './plugins/vuetify/icons';
import router from './router';
import { registerVitePreloadErrorHandler } from './services/vite-preload-error';
import { useSidebarStore } from './stores/modules/sidebar';

import 'vuetify/styles';

registerVitePreloadErrorHandler();

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(vuetify);
const vuetifyGlobalProperties = app.config.globalProperties.$vuetify ?? {};
vuetifyGlobalProperties.icons = {
  values: customIconsRegistry,
};
app.config.globalProperties.$vuetify = vuetifyGlobalProperties;
app.config.globalProperties.$constants = constants;
app.config.globalProperties.$sidebar = {
  show: payload => useSidebarStore().show(payload),
  hide: payload => useSidebarStore().hide(payload),
  updateConfig: payload => useSidebarStore().updateConfig(payload),
};
registerVuetifyComponentAliases(app);
app.use(i18n);
app.use(veeValidatePlugin, { i18n });
app.use(socketPlugin);

app.component('c-runtime-template', CRuntimeTemplate);
app.component('c-compiled-template', CCompiledTemplate);
app.component('c-request-helper', CRequestHelper);
app.component('c-now', CNow);
app.component('c-copy-wrapper', CCopyWrapper);
app.component('c-circle-badge', CCircleBadge);
app.component('c-progress-overlay', CProgressOverlay);
app.component('c-layout', CLayout);

Object.entries(BlockComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(ActionsPanelComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(ButtonComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(ChipComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(CalendarComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(ChartComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(DragDropComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(IconComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(InputComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(LayoutComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(ListComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(ModalComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(OverlayComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(PageComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(PaginationComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(PlaceholderComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(PopupComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(ResponsiveListComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(RequestComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(ScrollbarComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(SearchComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(TableComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(TooltipComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(TreeviewComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(VuetifyCommonComponents).forEach(([name, component]) => {
  app.component(name, component);
});

Object.entries(FormFieldComponents).forEach(([name, component]) => {
  app.component(name, component);
});

app.use(modalsPlugin, {
  components: RegisteredModals,
  dialogPropsMap: {
    [MODALS.addDynamicInfoInfosFromTemplate]: { maxWidth: 500, autoHeight: true },
    [MODALS.addInfoPopup]: { maxWidth: 700, persistent: true, autoHeight: true },
    [MODALS.aiChatHistory]: { maxWidth: 700, autoHeight: true },
    [MODALS.alarmsList]: { maxWidth: '95%', autoHeight: true },
    [MODALS.anomalyMonitoredConnectorHistory]: { maxWidth: 1400, autoHeight: true },
    [MODALS.applyEventFilter]: { maxWidth: 960 },
    [MODALS.clickOutsideConfirmation]: { autoHeight: true },
    [MODALS.colorPicker]: { autoHeight: true },
    [MODALS.confirmAckWithTicket]: { autoHeight: true },
    [MODALS.confirmation]: { autoHeight: true },
    [MODALS.confirmationPhrase]: { autoHeight: true },
    [MODALS.createAckEvent]: { autoHeight: true },
    [MODALS.createEvent]: { autoHeight: true },
    [MODALS.createAlarmChart]: { maxWidth: 500, autoHeight: true },
    [MODALS.createAlarmStatusRule]: { maxWidth: 1280 },
    [MODALS.createAnomalyMonitoredConnector]: { autoHeight: true },
    [MODALS.createAssociateTicketEvent]: { autoHeight: true },
    [MODALS.createChangeStateEvent]: { autoHeight: true },
    [MODALS.createCommentEvent]: { autoHeight: true },
    [MODALS.createCommentTemplate]: { autoHeight: true },
    [MODALS.createDeclareTicketEvent]: { maxWidth: 1280, autoHeight: true },
    [MODALS.createDeclareTicketRule]: { maxWidth: 1280 },
    [MODALS.createDynamicInfo]: { maxWidth: 1000 },
    [MODALS.createDynamicInfoTemplate]: { maxWidth: 600, autoHeight: true },
    [MODALS.createEntity]: { maxWidth: 1000 },
    [MODALS.createEntityInfo]: { autoHeight: true },
    [MODALS.createEntityInfoProperty]: { autoHeight: true },
    [MODALS.createEvent]: { autoHeight: true },
    [MODALS.createEventFilter]: { maxWidth: 1280 },
    [MODALS.createExternalAuthToken]: { maxWidth: 1100 },
    [MODALS.createExternalDataTable]: { autoHeight: true },
    [MODALS.createExternalDataTableRecord]: { autoHeight: true },
    [MODALS.createFilter]: { maxWidth: 1100 },
    [MODALS.createFlowchartMap]: { maxWidth: 1600, autoHeight: true },
    [MODALS.createGeoMap]: { maxWidth: 1280, autoHeight: true },
    [MODALS.createGroup]: { autoHeight: true },
    [MODALS.createIcon]: { maxWidth: 400, autoHeight: true },
    [MODALS.createIdleRule]: { maxWidth: 1280 },
    [MODALS.createJunitStateSetting]: { autoHeight: true },
    [MODALS.createKpiFilter]: { maxWidth: 1280, autoHeight: true },
    [MODALS.createLinkRule]: { maxWidth: 1000 },
    [MODALS.createLlm]: { maxWidth: 920, autoHeight: true },
    [MODALS.llmPromptsHistory]: { maxWidth: 1200 },
    [MODALS.createMaintenance]: { autoHeight: true },
    [MODALS.createMap]: { maxWidth: 500, autoHeight: true },
    [MODALS.createMermaidMap]: { maxWidth: 1600, autoHeight: true },
    [MODALS.createMetaAlarmRule]: { maxWidth: 1000 },
    [MODALS.createPattern]: { maxWidth: 1280 },
    [MODALS.createPbehaviorException]: { maxWidth: 1000, autoHeight: true },
    [MODALS.createPbehavior]: { maxWidth: 1280, autoHeight: true },
    [MODALS.createPbehaviorReason]: { autoHeight: true },
    [MODALS.createPbehaviorType]: { autoHeight: true },
    [MODALS.createPlaylist]: { maxWidth: 920, autoHeight: true },
    [MODALS.createRecurrenceRule]: { maxWidth: 1000, autoHeight: true },
    [MODALS.createRemediationConfiguration]: { autoHeight: true },
    [MODALS.createRemediationInstruction]: { maxWidth: 1200 },
    [MODALS.createRole]: { autoHeight: true },
    [MODALS.createScenario]: { maxWidth: 1280 },
    [MODALS.createService]: { maxWidth: 1000 },
    [MODALS.createShareToken]: { maxWidth: 1000, autoHeight: true },
    [MODALS.createApiToken]: { maxWidth: 1000, autoHeight: true },
    [MODALS.copyValue]: { maxWidth: 1000, autoHeight: true },
    [MODALS.createSnmpRule]: { maxWidth: 1000, autoHeight: true },
    [MODALS.createSnoozeEvent]: { autoHeight: true },
    [MODALS.createStateSetting]: { maxWidth: 1000 },
    [MODALS.createTag]: { maxWidth: 920 },
    [MODALS.createTemplateTestingData]: { autoHeight: true },
    [MODALS.createTemplateTestingTest]: { autoHeight: true },
    [MODALS.createTheme]: { maxWidth: 900, autoHeight: true },
    [MODALS.createTicketStatusJob]: { maxWidth: 1200, autoHeight: true },
    [MODALS.stopJobs]: { maxWidth: 720, autoHeight: true },
    [MODALS.createTreeOfDependenciesMap]: { maxWidth: 1334, autoHeight: true },
    [MODALS.createUser]: { autoHeight: true },
    [MODALS.createView]: { autoHeight: true },
    [MODALS.createWidget]: { maxWidth: 500, autoHeight: true },
    [MODALS.createWidgetTemplate]: { maxWidth: 920, autoHeight: true },
    [MODALS.duration]: { autoHeight: true },
    [MODALS.editLiveReporting]: { autoHeight: true },
    [MODALS.entitiesComparison]: { maxWidth: 1100, autoHeight: true },
    [MODALS.entitiesList]: { maxWidth: '95%', autoHeight: true },
    [MODALS.entitiesRootCauseDiagram]: { maxWidth: 1600, autoHeight: true },
    [MODALS.entityDependenciesList]: { maxWidth: 1600, autoHeight: true },
    [MODALS.entityUpstream]: { maxWidth: 1600, autoHeight: true },
    [MODALS.eventsRecord]: { maxWidth: 1600, persistent: true },
    [MODALS.executeDeclareTickets]: { maxWidth: 920, autoHeight: true },
    [MODALS.executeRemediationInstruction]: { maxWidth: 960, autoHeight: true },
    [MODALS.executeRemediationSimpleInstruction]: { maxWidth: 960, autoHeight: true },
    [MODALS.filtersList]: { autoHeight: true },
    [MODALS.healthcheckEngine]: { autoHeight: true },
    [MODALS.healthcheckEnginesChainReference]: { autoHeight: true },
    [MODALS.imageViewer]: {
      maxWidth: '90%',
      contentClass: 'v-dialog__image-viewer',
      autoHeight: true,
    },
    [MODALS.imagesViewer]: {
      maxWidth: '100%',
      contentClass: 'v-dialog__images-viewer',
      autoHeight: true,
    },
    [MODALS.importExportViews]: { maxWidth: 920, persistent: true, autoHeight: true },
    [MODALS.importExternalDataTableRecords]: { maxWidth: 1200, persistent: true, autoHeight: true },
    [MODALS.importPbehaviorException]: { autoHeight: true },
    [MODALS.info]: { autoHeight: true },
    [MODALS.infoPopupSetting]: { autoHeight: true },
    [MODALS.linkToMetaAlarm]: { maxWidth: 920, autoHeight: true },
    [MODALS.managePlaylistTabs]: { autoHeight: true },
    [MODALS.payloadTextareaEditor]: { autoHeight: true },
    [MODALS.pbehaviorPlanning]: { maxWidth: '95%', persistent: true, autoHeight: true },
    [MODALS.pbehaviorRecurrenceRule]: { maxWidth: '95%', persistent: true, autoHeight: true },
    [MODALS.pbehaviorRecurrentChangesConfirmation]: {
      maxWidth: 400,
      persistent: true,
      autoHeight: true,
    },
    [MODALS.pbehaviorsCalendar]: { maxWidth: '95%', persistent: true, autoHeight: true },
    [MODALS.rate]: { maxWidth: 500, autoHeight: true },
    [MODALS.remediationInstructionApproval]: { maxWidth: 960 },
    [MODALS.removeAlarmsFromMetaAlarm]: { autoHeight: true },
    [MODALS.removeAssociatedTicketEvent]: { autoHeight: true },
    [MODALS.selectView]: { autoHeight: true },
    [MODALS.selectViewTab]: { autoHeight: true },
    [MODALS.serviceDependencies]: { maxWidth: 1100, autoHeight: true },
    [MODALS.serviceEntities]: { maxWidth: 920 },
    [MODALS.stateSettingInheritedEntityPattern]: { maxWidth: 960, autoHeight: true },
    [MODALS.testSuite]: { maxWidth: 920 },
    [MODALS.textEditor]: { maxWidth: 900, persistent: true, autoHeight: true },
    [MODALS.textEditorPaste]: { autoHeight: true },
    [MODALS.textEditorWithTemplate]: { maxWidth: 900, persistent: true, autoHeight: true },
    [MODALS.textFieldEditor]: { autoHeight: true },
    [MODALS.userInterface]: { maxWidth: 1000 },
    [MODALS.variablesHelp]: { autoHeight: true },

    ...featuresService.get(FEATURES_SERVICE_KEYS.components.modals.dialogPropsMap, {}),
  },
});
app.use(popupsPlugin);

app.mount('#app');
