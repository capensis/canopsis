<template>
  <c-layout gap="1">
    <c-action-btn
      v-if="removable && someOneRemovable && activeConfig.remove"
      :tooltip="removeTooltip"
      type="delete"
      @click="showRemoveModal"
    />
    <c-action-btn
      v-if="enablable && someOneDisable && activeConfig.enable"
      :tooltip="t(`${activeConfig.tooltipPrefix}.massEnable`)"
      icon="check_circle"
      color="primary"
      @click="showEnableModal"
    />
    <c-action-btn
      v-if="disablable && someOneEnable && activeConfig.disable"
      :tooltip="t(`${activeConfig.tooltipPrefix}.massDisable`)"
      icon="cancel"
      color="error"
      @click="showDisableModal"
    />
    <c-action-btn
      v-if="unhideable && someOneHidden && activeConfig.unhide"
      :tooltip="t(`${activeConfig.tooltipPrefix}.massUnhide`)"
      icon="check_circle"
      color="primary"
      @click="showUnhideModal"
    />
    <c-action-btn
      v-if="hideable && someOneVisible && activeConfig.hide"
      :tooltip="t(`${activeConfig.tooltipPrefix}.massHide`)"
      icon="cancel"
      color="error"
      @click="showHideModal"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { MODALS, RESPONSE_STATUSES } from '@/constants';
import { pickIds } from '@/helpers/array';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useAnomalyMonitoredConnector } from '@/hooks/store/modules/anomaly-monitored-connector';
import { useDeclareTicketRule } from '@/hooks/store/modules/declare-ticket-rule';
import { useDynamicInfo } from '@/hooks/store/modules/dynamic-info';
import { useDynamicInfoTemplates } from '@/hooks/store/modules/dynamic-info-templates';
import { useEventFilter } from '@/hooks/store/modules/event-filter';
import { useFlappingRules } from '@/hooks/store/modules/flapping-rules';
import { useIdleRules } from '@/hooks/store/modules/idle-rules';
import { useLinkRule } from '@/hooks/store/modules/link-rule';
import { useLlm } from '@/hooks/store/modules/llm';
import { useMap } from '@/hooks/store/modules/map';
import { useMetaAlarmRule } from '@/hooks/store/modules/meta-alarm-rule';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';
import { usePbehaviorException } from '@/hooks/store/modules/pbehavior-exception';
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';
import { usePlaylist } from '@/hooks/store/modules/playlist';
import { useRemediationConfiguration } from '@/hooks/store/modules/remediation-configuration';
import { useRemediationInstruction } from '@/hooks/store/modules/remediation-instruction';
import { useRemediationJob } from '@/hooks/store/modules/remediation-job';
import { useResolveRules } from '@/hooks/store/modules/resolve-rules';
import { useRole } from '@/hooks/store/modules/role';
import { useScenario } from '@/hooks/store/modules/scenario';
import { useSnmpRule } from '@/hooks/store/modules/snmp-rule';
import { useUser } from '@/hooks/store/modules/user';

export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    removable: {
      type: Boolean,
      default: false,
    },
    enablable: {
      type: Boolean,
      default: false,
    },
    disablable: {
      type: Boolean,
      default: false,
    },
    hideable: {
      type: Boolean,
      default: false,
    },
    unhideable: {
      type: Boolean,
      default: false,
    },
    itemId: {
      type: String,
      default: '_id',
    },
    dynamicInfo: {
      type: Boolean,
      default: false,
    },
    dynamicInfoTemplate: {
      type: Boolean,
      default: false,
    },
    scenario: {
      type: Boolean,
      default: false,
    },
    flappingRule: {
      type: Boolean,
      default: false,
    },
    resolveRule: {
      type: Boolean,
      default: false,
    },
    idleRule: {
      type: Boolean,
      default: false,
    },
    linkRule: {
      type: Boolean,
      default: false,
    },
    declareTicket: {
      type: Boolean,
      default: false,
    },
    llm: {
      type: Boolean,
      default: false,
    },
    map: {
      type: Boolean,
      default: false,
    },
    metaAlarmRule: {
      type: Boolean,
      default: false,
    },
    pbehavior: {
      type: Boolean,
      default: false,
    },
    eventFilter: {
      type: Boolean,
      default: false,
    },
    playlist: {
      type: Boolean,
      default: false,
    },
    pbehaviorType: {
      type: Boolean,
      default: false,
    },
    pbehaviorReason: {
      type: Boolean,
      default: false,
    },
    pbehaviorException: {
      type: Boolean,
      default: false,
    },
    user: {
      type: Boolean,
      default: false,
    },
    role: {
      type: Boolean,
      default: false,
    },
    instruction: {
      type: Boolean,
      default: false,
    },
    remediationConfiguration: {
      type: Boolean,
      default: false,
    },
    remediationJob: {
      type: Boolean,
      default: false,
    },
    anomalyMonitoredConnector: {
      type: Boolean,
      default: false,
    },
    snmpRule: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['clear:items', 'refresh'],

  setup(props, { emit }) {
    const { te, t } = useI18n();
    const modals = useModals();
    const popups = usePopups();

    const { bulkEnableDynamicInfos, bulkDisableDynamicInfos, bulkRemoveDynamicInfos } =
      useDynamicInfo();

    const { bulkRemoveDynamicInfoTemplates } = useDynamicInfoTemplates();

    const { bulkEnableScenarios, bulkDisableScenarios, bulkRemoveScenarios } = useScenario();
    const { bulkEnableEventFilters, bulkDisableEventFilters, bulkRemoveEventFilters } =
      useEventFilter();
    const { bulkEnableFlappingRules, bulkDisableFlappingRules, bulkRemoveFlappingRules } =
      useFlappingRules();
    const { bulkEnableResolveRules, bulkDisableResolveRules, bulkRemoveResolveRules } =
      useResolveRules();
    const { bulkEnableIdleRules, bulkDisableIdleRules, bulkRemoveIdleRules } = useIdleRules();
    const { bulkEnableLinkRules, bulkDisableLinkRules, bulkRemoveLinkRules } = useLinkRule();
    const {
      bulkEnableDeclareTicketRules,
      bulkDisableDeclareTicketRules,
      bulkRemoveDeclareTicketRules,
    } = useDeclareTicketRule();
    const { bulkRemoveLlms, bulkEnableLlms, bulkDisableLlms } = useLlm();
    const { bulkRemoveMaps } = useMap();
    const { bulkEnableMetaAlarmRules, bulkDisableMetaAlarmRules, bulkRemoveMetaAlarmRules } =
      useMetaAlarmRule();
    const {
      bulkEnableAnomalyMonitoredConnectors,
      bulkDisableAnomalyMonitoredConnectors,
      bulkRemoveAnomalyMonitoredConnectors,
    } = useAnomalyMonitoredConnector();
    const { bulkEnablePbehaviors, bulkDisablePbehaviors, bulkRemovePbehaviors } = usePbehavior();
    const { bulkEnablePlaylist, bulkDisablePlaylist, bulkRemovePlaylist } = usePlaylist();
    const { bulkHidePbehaviorTypes, bulkUnhidePbehaviorTypes, bulkRemovePbehaviorTypes } =
      usePbehaviorType();
    const { bulkHidePbehaviorReasons, bulkUnhidePbehaviorReasons, bulkRemovePbehaviorReasons } =
      usePbehaviorReason();
    const {
      bulkHidePbehaviorExceptions,
      bulkUnhidePbehaviorExceptions,
      bulkRemovePbehaviorExceptions,
    } = usePbehaviorException();
    const { bulkRemoveUsers, bulkEnableUsers, bulkDisableUsers } = useUser();
    const { bulkRemoveRoles } = useRole();
    const {
      bulkRemoveRemediationInstructions,
      bulkEnableRemediationInstructions,
      bulkDisableRemediationInstructions,
    } = useRemediationInstruction();
    const { bulkRemoveRemediationConfigurations } = useRemediationConfiguration();
    const { bulkRemoveRemediationJobs } = useRemediationJob();
    const { bulkEnableSnmpRules, bulkDisableSnmpRules, bulkRemoveSnmpRules } = useSnmpRule();

    const someOneEnable = computed(() => props.items.some(({ enabled }) => enabled));
    const someOneDisable = computed(() => props.items.some(({ enabled }) => !enabled));
    const someOneRemovable = computed(() => props.items.some(({ deletable = true }) => deletable));
    const someOneVisible = computed(() => props.items.some(({ hidden }) => !hidden));
    const someOneHidden = computed(() => props.items.some(({ hidden }) => hidden));

    const activeConfig = computed(() => {
      if (props.dynamicInfo) {
        return {
          remove: bulkRemoveDynamicInfos,
          enable: bulkEnableDynamicInfos,
          disable: bulkDisableDynamicInfos,
          tooltipPrefix: 'dynamicInfo',
        };
      }

      if (props.dynamicInfoTemplate) {
        return {
          remove: bulkRemoveDynamicInfoTemplates,
          tooltipPrefix: 'dynamicInfo.templatesList',
        };
      }

      if (props.scenario) {
        return {
          remove: bulkRemoveScenarios,
          enable: bulkEnableScenarios,
          disable: bulkDisableScenarios,
          tooltipPrefix: 'scenario',
        };
      }

      if (props.eventFilter) {
        return {
          remove: bulkRemoveEventFilters,
          enable: bulkEnableEventFilters,
          disable: bulkDisableEventFilters,
          tooltipPrefix: 'eventFilter',
        };
      }

      if (props.flappingRule) {
        return {
          remove: bulkRemoveFlappingRules,
          enable: bulkEnableFlappingRules,
          disable: bulkDisableFlappingRules,
          tooltipPrefix: 'flappingRule',
        };
      }

      if (props.resolveRule) {
        return {
          remove: bulkRemoveResolveRules,
          enable: bulkEnableResolveRules,
          disable: bulkDisableResolveRules,
          tooltipPrefix: 'resolveRule',
        };
      }

      if (props.idleRule) {
        return {
          remove: bulkRemoveIdleRules,
          enable: bulkEnableIdleRules,
          disable: bulkDisableIdleRules,
          tooltipPrefix: 'idleRules',
        };
      }

      if (props.linkRule) {
        return {
          remove: bulkRemoveLinkRules,
          enable: bulkEnableLinkRules,
          disable: bulkDisableLinkRules,
          tooltipPrefix: 'linkRule',
        };
      }

      if (props.declareTicket) {
        return {
          remove: bulkRemoveDeclareTicketRules,
          enable: bulkEnableDeclareTicketRules,
          disable: bulkDisableDeclareTicketRules,
          tooltipPrefix: 'declareTicket',
        };
      }

      if (props.llm) {
        return {
          remove: bulkRemoveLlms,
          enable: bulkEnableLlms,
          disable: bulkDisableLlms,
          tooltipPrefix: 'llm',
          removeConfirmation: 'modals.confirmationPhrase.deleteSelectedLlms',
        };
      }

      if (props.map) {
        return {
          remove: bulkRemoveMaps,
          tooltipPrefix: 'map',
        };
      }

      if (props.metaAlarmRule) {
        return {
          remove: bulkRemoveMetaAlarmRules,
          enable: bulkEnableMetaAlarmRules,
          disable: bulkDisableMetaAlarmRules,
          tooltipPrefix: 'metaAlarmRule',
        };
      }

      if (props.pbehavior) {
        return {
          remove: bulkRemovePbehaviors,
          enable: bulkEnablePbehaviors,
          disable: bulkDisablePbehaviors,
          tooltipPrefix: 'pbehavior',
        };
      }

      if (props.playlist) {
        return {
          remove: bulkRemovePlaylist,
          enable: bulkEnablePlaylist,
          disable: bulkDisablePlaylist,
          tooltipPrefix: 'playlist',
        };
      }

      if (props.pbehaviorType) {
        return {
          remove: bulkRemovePbehaviorTypes,
          hide: bulkHidePbehaviorTypes,
          unhide: bulkUnhidePbehaviorTypes,
          tooltipPrefix: 'pbehavior',
        };
      }

      if (props.pbehaviorReason) {
        return {
          remove: bulkRemovePbehaviorReasons,
          hide: bulkHidePbehaviorReasons,
          unhide: bulkUnhidePbehaviorReasons,
          tooltipPrefix: 'pbehavior',
        };
      }

      if (props.pbehaviorException) {
        return {
          remove: bulkRemovePbehaviorExceptions,
          hide: bulkHidePbehaviorExceptions,
          unhide: bulkUnhidePbehaviorExceptions,
          tooltipPrefix: 'pbehavior',
        };
      }

      if (props.user) {
        return {
          remove: bulkRemoveUsers,
          enable: bulkEnableUsers,
          disable: bulkDisableUsers,
          tooltipPrefix: 'user',
        };
      }

      if (props.role) {
        return {
          remove: bulkRemoveRoles,
          tooltipPrefix: 'role',
        };
      }

      if (props.instruction) {
        return {
          remove: bulkRemoveRemediationInstructions,
          enable: bulkEnableRemediationInstructions,
          disable: bulkDisableRemediationInstructions,
          tooltipPrefix: 'remediation.instruction',
        };
      }

      if (props.remediationConfiguration) {
        return {
          remove: bulkRemoveRemediationConfigurations,
          tooltipPrefix: 'remediation.configuration',
        };
      }

      if (props.remediationJob) {
        return {
          remove: bulkRemoveRemediationJobs,
          tooltipPrefix: 'remediation.job',
        };
      }

      if (props.anomalyMonitoredConnector) {
        return {
          remove: bulkRemoveAnomalyMonitoredConnectors,
          enable: bulkEnableAnomalyMonitoredConnectors,
          disable: bulkDisableAnomalyMonitoredConnectors,
          tooltipPrefix: 'anomalyMonitoredConnector',
        };
      }

      if (props.snmpRule) {
        return {
          remove: bulkRemoveSnmpRules,
          enable: bulkEnableSnmpRules,
          disable: bulkDisableSnmpRules,
          tooltipPrefix: 'snmpRule',
        };
      }

      return {};
    });

    const removeTooltip = computed(() => {
      const massRemoveTooltipKey = `${activeConfig.value.tooltipPrefix}.massRemove`;

      return te(massRemoveTooltipKey) ? t(massRemoveTooltipKey) : '';
    });

    const afterSubmit = async () => {
      emit('clear:items');
      emit('refresh');
    };

    const showErrorPopups = (response = [], messageKey) => {
      response.forEach(({ status, message, error }, index) => {
        if (status !== RESPONSE_STATUSES.success) {
          const ruleName = props.items[index]?.name;
          const text = te(messageKey) ? t(messageKey, { name: ruleName }) : message || error;

          popups.error({ text });
        }
      });
    };

    const showRemoveModal = () => {
      if (activeConfig.value.removeConfirmation) {
        const confirmationKey = activeConfig.value.removeConfirmation;
        const count = props.items.length;

        modals.show({
          name: MODALS.confirmationPhrase,
          config: {
            title: t(`${confirmationKey}.title`, count, { count }),
            text: t(`${confirmationKey}.text`, count, { count }),
            phraseText: t(`${confirmationKey}.phraseText`),
            phrase: t(`${confirmationKey}.phrase`),
            action: async () => {
              const response = await activeConfig.value.remove({
                data: pickIds(props.items, props.itemId),
              });

              showErrorPopups(response, `${activeConfig.value.tooltipPrefix}.removeForbidden`);

              return afterSubmit();
            },
          },
        });

        return;
      }

      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const response = await activeConfig.value.remove({
              data: pickIds(props.items, props.itemId),
            });

            showErrorPopups(response, `${activeConfig.value.tooltipPrefix}.removeForbidden`);

            return afterSubmit();
          },
        },
      });
    };

    const showEnableModal = () =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await activeConfig.value.enable({ data: pickIds(props.items, props.itemId) });

            return afterSubmit();
          },
        },
      });

    const showDisableModal = () =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const response = await activeConfig.value.disable({
              data: pickIds(props.items, props.itemId),
            });

            showErrorPopups(response, `${activeConfig.value.tooltipPrefix}.disableForbidden`);

            return afterSubmit();
          },
        },
      });

    const showUnhideModal = () =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await activeConfig.value.unhide({ data: pickIds(props.items, props.itemId) });

            return afterSubmit();
          },
        },
      });

    const showHideModal = () =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const response = await activeConfig.value.hide({
              data: pickIds(props.items, props.itemId),
            });

            showErrorPopups(response, `${activeConfig.value.tooltipPrefix}.hideForbidden`);

            return afterSubmit();
          },
        },
      });

    return {
      t,
      activeConfig,
      removeTooltip,
      someOneEnable,
      someOneDisable,
      someOneRemovable,
      someOneVisible,
      someOneHidden,
      showRemoveModal,
      showEnableModal,
      showDisableModal,
      showUnhideModal,
      showHideModal,
    };
  },
};
</script>
