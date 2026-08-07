import { createRouter, createWebHistory } from 'vue-router';

import {
  CRUD_ACTIONS,
  GROUPED_USER_PERMISSIONS,
  OLD_ROUTES,
  ROUTES,
  ROUTES_NAMES,
  USER_PERMISSIONS,
} from '@/constants';
import { ROUTER_ACCESS_TOKEN_KEY } from '@/config';

import { usePopupsStore } from '@/plugins/popups/stores';
import { useAuthStore } from '@/stores/modules/auth';

import { checkAppInfoAccessForRoute, checkUserAccessForRoute } from '@/helpers/router';

const PlaceholderPage = () => import('@/views/placeholder-page.vue');
const AdminRights = () => import('@/views/admin/rights-page.vue');
const AdminUsers = () => import('@/views/admin/users-page.vue');
const AdminRoles = () => import('@/views/admin/roles-page.vue');
const AdminBroadcastMessages = () => import('@/views/admin/broadcast-messages-page.vue');
const AdminPlaylists = () => import('@/views/admin/playlists-page.vue');
const AdminPlanning = () => import('@/views/admin/planning-page.vue');
const AdminHealthcheck = () => import('@/views/admin/healthcheck-page.vue');
const AdminKPI = () => import('@/views/admin/kpi-page.vue');
const AdminEventsRecords = () => import('@/views/admin/events-records-page.vue');
const AdminAnomalyMonitoredConnectors = () =>
  import('@/views/admin/anomaly-monitored-connectors-page.vue');
const AdminEntityInfosProperties = () =>
  import('@/views/admin/custom-objects/entity-infos-properties-page.vue');
const AdminIcons = () => import('@/views/admin/custom-objects/icons-page.vue');
const AdminMaps = () => import('@/views/admin/custom-objects/maps-page.vue');
const AdminTags = () => import('@/views/admin/custom-objects/tags-page.vue');
const AdminLlms = () => import('@/views/admin/custom-objects/llms-page.vue');
const AdminViewsImportExport = () => import('@/views/admin/settings/views-import-export-page.vue');
const AdminSettingsNotifications = () =>
  import('@/views/admin/settings/notifications-settings-page.vue');
const AdminSettingsCommentTemplates = () =>
  import('@/views/admin/settings/comment-templates-page.vue');
const AdminSettingsStorageSettings = () =>
  import('@/views/admin/settings/storage-settings-page.vue');
const AdminSettingsStateSettings = () => import('@/views/admin/settings/state-settings-page.vue');
const DynamicInfosPage = () => import('@/views/exploitation/dynamic-infos-page.vue');
const DeclareTicketRulesPage = () => import('@/views/exploitation/declare-ticket-rules-page.vue');
const EventFiltersPage = () => import('@/views/exploitation/event-filters-page.vue');
const FlappingRulesPage = () => import('@/views/exploitation/flapping-rules-page.vue');
const IdleRulesPage = () => import('@/views/exploitation/idle-rules-page.vue');
const LinkRulesPage = () => import('@/views/exploitation/link-rules-page.vue');
const MetaAlarmRulesPage = () => import('@/views/exploitation/meta-alarm-rules-page.vue');
const PbehaviorsPage = () => import('@/views/exploitation/pbehaviors-page.vue');
const ResolveRulesPage = () => import('@/views/exploitation/resolve-rules-page.vue');
const RemediationPage = () => import('@/views/exploitation/remediation-page.vue');
const ScenariosPage = () => import('@/views/exploitation/scenarios-page.vue');
const SnmpRulesPage = () => import('@/views/exploitation/snmp-rules-page.vue');
const ProfilePatternsPage = () => import('@/views/profile/patterns-page.vue');
const ProfileThemesPage = () => import('@/views/profile/themes-page.vue');
const NotificationsPage = () => import('@/views/notification/notifications-page.vue');
const ViewPage = () => import('@/views/view-page.vue');

const protectedRoute = (meta = {}) => ({
  requiresLogin: true,
  ...meta,
});

const routes = [
  {
    path: ROUTES.login,
    name: ROUTES_NAMES.login,
    meta: {
      hideHeader: true,
      requiresLogin: false,
    },
    component: () => import('@/views/login-page.vue'),
  },
  {
    path: ROUTES.home,
    name: ROUTES_NAMES.home,
    meta: protectedRoute(),
    component: () => import('@/views/home-page.vue'),
  },
  {
    path: ROUTES.view,
    name: ROUTES_NAMES.view,
    meta: protectedRoute(),
    component: ViewPage,
    props: route => ({ id: route.params.id }),
  },
  {
    path: ROUTES.viewKiosk,
    name: ROUTES_NAMES.viewKiosk,
    meta: protectedRoute({
      simpleNavigation: true,
      requiresPermission: {
        id: route => route.params.id,
      },
    }),
    component: PlaceholderPage,
    props: route => ({ id: route.params.id, tabId: route.params.tabId }),
  },
  {
    path: ROUTES.alarms,
    name: ROUTES_NAMES.alarms,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.view,
      },
    }),
    component: PlaceholderPage,
    props: route => ({ id: route.params.id }),
  },
  {
    path: ROUTES.adminRights,
    name: ROUTES_NAMES.adminRights,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.permission,
      },
    }),
    component: AdminRights,
  },
  {
    path: ROUTES.adminUsers,
    name: ROUTES_NAMES.adminUsers,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.user,
      },
    }),
    component: AdminUsers,
  },
  {
    path: ROUTES.adminRoles,
    name: ROUTES_NAMES.adminRoles,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.role,
      },
    }),
    component: AdminRoles,
  },
  {
    path: ROUTES.adminBroadcastMessages,
    name: ROUTES_NAMES.adminBroadcastMessages,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.broadcastMessage,
      },
    }),
    component: AdminBroadcastMessages,
  },
  {
    path: ROUTES.adminPlaylists,
    name: ROUTES_NAMES.adminPlaylists,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.playlist,
      },
    }),
    component: AdminPlaylists,
  },
  {
    path: ROUTES.adminPlanning,
    name: ROUTES_NAMES.adminPlanning,
    meta: protectedRoute({
      requiresPermission: {
        id: GROUPED_USER_PERMISSIONS.planning,
      },
    }),
    component: AdminPlanning,
  },
  {
    path: ROUTES.adminHealthcheck,
    name: ROUTES_NAMES.adminHealthcheck,
    meta: protectedRoute({
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.healthcheck,
      },
    }),
    component: AdminHealthcheck,
  },
  {
    path: ROUTES.adminKPI,
    name: ROUTES_NAMES.adminKPI,
    meta: protectedRoute({
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.kpi,
      },
    }),
    component: AdminKPI,
  },
  {
    path: ROUTES.adminEventsRecords,
    name: ROUTES_NAMES.adminEventsRecords,
    meta: protectedRoute({
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.eventsRecord,
      },
    }),
    component: AdminEventsRecords,
  },
  {
    path: ROUTES.adminAnomalyMonitoredConnectors,
    name: ROUTES_NAMES.adminAnomalyMonitoredConnectors,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.anomalyMonitoredConnector,
      },
    }),
    component: AdminAnomalyMonitoredConnectors,
  },
  {
    path: ROUTES.adminTemplateTesting,
    name: ROUTES_NAMES.adminTemplateTesting,
    meta: protectedRoute({
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.templateTesting,
      },
    }),
    component: () => import('@/views/admin/template-testing-page.vue'),
  },
  {
    path: ROUTES.adminJobsManagement,
    name: ROUTES_NAMES.adminJobsManagement,
    meta: protectedRoute({
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.jobManagement,
      },
    }),
    component: () => import('@/views/admin/jobs-page.vue'),
  },
  {
    path: ROUTES.adminCustomObjectsExternalAuthTokens,
    name: ROUTES_NAMES.adminCustomObjectsExternalAuthTokens,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.externalAuthTokens,
      },
    }),
    component: () => import('@/views/admin/custom-objects/external-auth-tokens-page.vue'),
  },
  {
    path: ROUTES.adminCustomObjectsExternalDataTables,
    name: ROUTES_NAMES.adminCustomObjectsExternalDataTables,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.externalDataTable,
      },
    }),
    component: () => import('@/views/admin/custom-objects/external-data-tables-page.vue'),
  },
  {
    path: ROUTES.adminCustomObjectsEntityInfosProperties,
    name: ROUTES_NAMES.adminCustomObjectsEntityInfosProperties,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.entityInfoProperty,
      },
    }),
    component: AdminEntityInfosProperties,
  },
  {
    path: ROUTES.adminCustomObjectsIcons,
    name: ROUTES_NAMES.adminCustomObjectsIcons,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.icon,
      },
    }),
    component: AdminIcons,
  },
  {
    path: ROUTES.adminCustomObjectsMaps,
    name: ROUTES_NAMES.adminCustomObjectsMaps,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.map,
      },
    }),
    component: AdminMaps,
  },
  {
    path: ROUTES.adminCustomObjectsTags,
    name: ROUTES_NAMES.adminCustomObjectsTags,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.tag,
      },
    }),
    component: AdminTags,
  },
  {
    path: ROUTES.adminCustomObjectsLlms,
    name: ROUTES_NAMES.adminCustomObjectsLlms,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.llm,
      },
    }),
    component: AdminLlms,
  },
  {
    path: ROUTES.adminSettingsViewsImportExport,
    name: ROUTES_NAMES.adminSettingsViewsImportExport,
    meta: protectedRoute({
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.viewImportExport,
      },
    }),
    component: AdminViewsImportExport,
  },
  {
    path: ROUTES.adminSettingsNotifications,
    name: ROUTES_NAMES.adminSettingsNotifications,
    meta: protectedRoute({
      requiresPermission: {
        action: CRUD_ACTIONS.can,
        id: USER_PERMISSIONS.technical.notification.common,
      },
    }),
    component: AdminSettingsNotifications,
  },
  {
    path: ROUTES.adminSettingsCommentTemplates,
    name: ROUTES_NAMES.adminSettingsCommentTemplates,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.commentTemplate,
      },
    }),
    component: AdminSettingsCommentTemplates,
  },
  {
    path: ROUTES.adminSettingsWidgetTemplates,
    name: ROUTES_NAMES.adminSettingsWidgetTemplates,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.widgetTemplate,
      },
    }),
    component: DeclareTicketRulesPage,
  },
  {
    path: ROUTES.adminSettingsStorageSettings,
    name: ROUTES_NAMES.adminSettingsStorageSettings,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.storageSettings,
      },
    }),
    component: AdminSettingsStorageSettings,
  },
  {
    path: ROUTES.adminSettingsStateSettings,
    name: ROUTES_NAMES.adminSettingsStateSettings,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.stateSetting,
      },
    }),
    component: AdminSettingsStateSettings,
  },
  {
    path: ROUTES.exploitationPbehaviors,
    name: ROUTES_NAMES.exploitationPbehaviors,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.pbehavior,
      },
    }),
    component: PbehaviorsPage,
  },
  {
    path: ROUTES.exploitationEventFilters,
    name: ROUTES_NAMES.exploitationEventFilters,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.eventFilter,
      },
    }),
    component: EventFiltersPage,
  },
  {
    path: ROUTES.exploitationSnmpRules,
    name: ROUTES_NAMES.exploitationSnmpRules,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.snmpRule,
      },
    }),
    component: SnmpRulesPage,
  },
  {
    path: ROUTES.exploitationDynamicInfos,
    name: ROUTES_NAMES.exploitationDynamicInfos,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.dynamicInfo,
      },
    }),
    component: DynamicInfosPage,
  },
  {
    path: ROUTES.exploitationMetaAlarmRules,
    name: ROUTES_NAMES.exploitationMetaAlarmRules,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.metaAlarmRule,
      },
    }),
    component: MetaAlarmRulesPage,
  },
  {
    path: ROUTES.exploitationScenarios,
    name: ROUTES_NAMES.exploitationScenarios,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.scenario,
      },
    }),
    component: ScenariosPage,
  },
  {
    path: ROUTES.exploitationIdleRules,
    name: ROUTES_NAMES.exploitationIdleRules,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.idleRules,
      },
    }),
    component: IdleRulesPage,
  },
  {
    path: ROUTES.exploitationFlappingRules,
    name: ROUTES_NAMES.exploitationFlappingRules,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.flappingRules,
      },
    }),
    component: FlappingRulesPage,
  },
  {
    path: ROUTES.exploitationResolveRules,
    name: ROUTES_NAMES.exploitationResolveRules,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.resolveRules,
      },
    }),
    component: ResolveRulesPage,
  },
  {
    path: ROUTES.exploitationDeclareTicketRules,
    name: ROUTES_NAMES.exploitationDeclareTicketRules,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.declareTicketRule,
      },
    }),
    component: DeclareTicketRulesPage,
  },
  {
    path: ROUTES.exploitationLinkRules,
    name: ROUTES_NAMES.exploitationLinkRules,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.exploitation.linkRule,
      },
    }),
    component: LinkRulesPage,
  },
  {
    path: ROUTES.exploitationRemediation,
    name: ROUTES_NAMES.exploitationRemediation,
    meta: protectedRoute({
      requiresPermission: {
        id: GROUPED_USER_PERMISSIONS.remediation,
      },
    }),
    component: RemediationPage,
  },
  {
    path: ROUTES.playlist,
    name: ROUTES_NAMES.playlist,
    meta: protectedRoute({
      requiresPermission: {
        action: CRUD_ACTIONS.read,
        id: route => route.params.id,
      },
    }),
    component: PlaceholderPage,
    props: route => ({ id: route.params.id, autoplay: String(route.query.autoplay) === 'true' }),
  },
  {
    path: ROUTES.profilePatterns,
    name: ROUTES_NAMES.profilePatterns,
    meta: protectedRoute(),
    component: ProfilePatternsPage,
  },
  {
    path: ROUTES.profileThemes,
    name: ROUTES_NAMES.profileThemes,
    meta: protectedRoute({
      requiresPermission: {
        id: USER_PERMISSIONS.technical.profile.theme,
      },
    }),
    component: ProfileThemesPage,
  },
  {
    path: ROUTES.notifications,
    name: ROUTES_NAMES.notifications,
    meta: protectedRoute(),
    component: NotificationsPage,
    props: route => ({ tabId: route.params.tabId, activeId: route.query.id }),
  },
  {
    path: ROUTES.error,
    name: ROUTES_NAMES.error,
    meta: {
      hideHeader: true,
    },
    component: () => import('@/views/error-page.vue'),
    props: route => ({
      message: route.query.message,
      redirect: route.query.redirect,
    }),
  },
  {
    path: OLD_ROUTES.remediation,
    redirect: { name: ROUTES_NAMES.exploitationRemediation },
  },
  {
    path: OLD_ROUTES.externalDataTables,
    redirect: { name: ROUTES_NAMES.adminCustomObjectsExternalDataTables },
  },
  {
    path: OLD_ROUTES.entityInfosProperties,
    redirect: { name: ROUTES_NAMES.adminCustomObjectsEntityInfosProperties },
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: { name: ROUTES_NAMES.home },
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach(async (to) => {
  const authStore = useAuthStore();
  const isRequiresAuth = to.matched.some(route => route.meta.requiresLogin);
  const isDontRequiresAuth = to.matched.some(route => route.meta.requiresLogin === false);
  const { [ROUTER_ACCESS_TOKEN_KEY]: accessToken, ...restQuery } = to.query;

  if (accessToken) {
    await authStore.applyAccessToken(accessToken);

    return {
      ...to,
      query: restQuery,
      replace: true,
    };
  }

  if (!authStore.isLoggedIn && isRequiresAuth) {
    return {
      name: ROUTES_NAMES.login,
      query: {
        redirect: to.fullPath,
        errorMessage: to.query.errorMessage,
      },
    };
  }

  if (authStore.isLoggedIn && isDontRequiresAuth) {
    return { name: ROUTES_NAMES.home };
  }

  return true;
});

router.beforeResolve(async (to) => {
  try {
    await checkAppInfoAccessForRoute(to);
    await checkUserAccessForRoute(to);

    return true;
  } catch (err) {
    console.error(err);

    return { name: ROUTES_NAMES.home };
  }
});

router.afterEach((to, from) => {
  if (to.path === from.path) {
    return;
  }

  // Entities sweep will be wired when the entity store is migrated.
});

router.isReady().then(() => {
  const { errorMessage } = router.currentRoute.value.query;

  if (errorMessage) {
    usePopupsStore().error({ text: errorMessage, autoClose: false });
  }
});

router.replaceAsync = route => router.replace(route);

export default router;
