import VueMq from 'vue-mq';
import VueFullScreen from 'vue-fullscreen';
import VueClipboard from 'vue-clipboard2';
import PortalVue from 'portal-vue';

import { MODALS } from '@/constants';
import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';

import ValidatorPlugin from '@/plugins/validator';
import ModalsPlugin from '@/plugins/modals';
import PopupsPlugin from '@/plugins/popups';
import SidebarPlugin from '@/plugins/sidebar';
import SetSeveralPlugin from '@/plugins/set-several';
import SetOnlyDiffPlugin from '@/plugins/set-only-diff';
import UpdateFieldPlugin from '@/plugins/update-field';
import SocketPlugin from '@/plugins/socket';

import { featuresService } from '@/services/features';

import store from '@/store';
import i18n from '@/i18n';
import Filters from '@/filters';

import * as modalsComponents from '@/components/modals';
import * as sidebarsComponents from '@/components/sidebars';

/**
 * @param {import('vue').VueConstructor | import('vue').Vue} Vue
 * @returns {*}
 */
export const bootstrapApplicationPlugins = (Vue) => {
  Vue.use(PortalVue);
  Vue.use(Filters);
  Vue.use(VueFullScreen);

  Vue.use(VueMq, {
    breakpoints: MEDIA_QUERIES_BREAKPOINTS,
  });

  VueClipboard.config.autoSetContainer = true;
  Vue.use(VueClipboard);

  Vue.use(ValidatorPlugin, { i18n });

  Vue.use(ModalsPlugin, {
    store,

    components: {
      ...modalsComponents,
      ...featuresService.get('components.modals.components'),
    },

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
      [MODALS.createMaintenance]: { autoHeight: true },
      [MODALS.createMap]: { maxWidth: 500, autoHeight: true },
      [MODALS.createMermaidMap]: { maxWidth: 1600, autoHeight: true },
      [MODALS.createMetaAlarmRule]: { maxWidth: 1000 },
      [MODALS.createPattern]: { maxWidth: 1280 },
      [MODALS.createPbehaviorException]: { autoHeight: true },
      [MODALS.createPbehaviorReason]: { autoHeight: true },
      [MODALS.createPbehaviorType]: { autoHeight: true },
      [MODALS.createPlaylist]: { maxWidth: 920, autoHeight: true },
      [MODALS.createRecurrenceRule]: { maxWidth: 1000, autoHeight: true },
      [MODALS.createRemediationConfiguration]: { autoHeight: true },
      [MODALS.createRemediationInstruction]: { maxWidth: 1200 },
      [MODALS.createRole]: { autoHeight: true },
      [MODALS.createScenario]: { maxWidth: 1280 },
      [MODALS.createService]: { maxWidth: 1280 },
      [MODALS.createServicePauseEvent]: { autoHeight: true },
      [MODALS.createShareToken]: { autoHeight: true },
      [MODALS.createSnmpRule]: { maxWidth: 1000, autoHeight: true },
      [MODALS.createSnoozeEvent]: { autoHeight: true },
      [MODALS.createStateSetting]: { maxWidth: 1000 },
      [MODALS.createTag]: { maxWidth: 920 },
      [MODALS.createTemplateTestingData]: { autoHeight: true },
      [MODALS.createTemplateTestingTest]: { autoHeight: true },
      [MODALS.createTheme]: { maxWidth: 500, autoHeight: true },
      [MODALS.createTicketStatusJob]: { maxWidth: 1200, autoHeight: true },
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
      [MODALS.filtersList]: { autoHeight: true },
      [MODALS.healthcheckEngine]: { autoHeight: true },
      [MODALS.healthcheckEnginesChainReference]: { autoHeight: true },
      [MODALS.imageViewer]: { maxWidth: '90%', contentClass: 'v-dialog__image-viewer', autoHeight: true },
      [MODALS.imagesViewer]: { maxWidth: '100%', contentClass: 'v-dialog__images-viewer', autoHeight: true },
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
      [MODALS.pbehaviorRecurrentChangesConfirmation]: { maxWidth: 400, persistent: true, autoHeight: true },
      [MODALS.pbehaviorsCalendar]: { maxWidth: '95%', persistent: true, autoHeight: true },
      [MODALS.rate]: { maxWidth: 500, autoHeight: true },
      [MODALS.remediationInstructionApproval]: { maxWidth: 960 },
      [MODALS.removeAlarmsFromMetaAlarm]: { autoHeight: true },
      [MODALS.removeAssociatedTicketEvent]: { autoHeight: true },
      [MODALS.selectView]: { autoHeight: true },
      [MODALS.selectViewTab]: { autoHeight: true },
      [MODALS.serviceDependencies]: { maxWidth: 1100, autoHeight: true },
      [MODALS.serviceEntities]: { maxWidth: 920 },
      [MODALS.shareView]: { autoHeight: true },
      [MODALS.stateSettingInheritedEntityPattern]: { maxWidth: 960, autoHeight: true },
      [MODALS.testSuite]: { maxWidth: 920 },
      [MODALS.textEditor]: { maxWidth: 900, persistent: true, autoHeight: true },
      [MODALS.textEditorWithTemplate]: { maxWidth: 900, persistent: true, autoHeight: true },
      [MODALS.textFieldEditor]: { autoHeight: true },
      [MODALS.userInterface]: { maxWidth: 1000 },
      [MODALS.variablesHelp]: { autoHeight: true },

      ...featuresService.get('components.modals.dialogPropsMap'),
    },
  });

  Vue.use(PopupsPlugin, { store });
  Vue.use(SidebarPlugin, {
    store,

    components: {
      ...sidebarsComponents,
      ...featuresService.get('components.sidebars.components'),
    },
  });
  Vue.use(SetSeveralPlugin);
  Vue.use(SetOnlyDiffPlugin);
  Vue.use(UpdateFieldPlugin);
  Vue.use(SocketPlugin);
};
