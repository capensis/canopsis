<template>
  <v-layout>
    <c-action-btn
      v-if="removable && config.remove"
      :tooltip="$t(`${config.tooltipPrefix}.massRemove`)"
      type="delete"
      @click="showRemoveModal"
    />
    <c-action-btn
      v-if="enablable && someOneDisable && config.enable"
      :tooltip="$t(`${config.tooltipPrefix}.massEnable`)"
      icon="check_circle"
      color="primary"
      @click="showEnableModal"
    />
    <c-action-btn
      v-if="disablable && someOneEnable && config.disable"
      :tooltip="$t(`${config.tooltipPrefix}.massDisable`)"
      icon="cancel"
      color="error"
      @click="showDisableModal"
    />
    <c-db-export-btn
      v-if="config.exportProps"
      :ids="itemsIds"
      v-bind="config.exportProps"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { mapIds } from '@/helpers/array';

import { useModals } from '@/hooks/modals';
import { useDynamicInfo } from '@/hooks/store/modules/dynamic-info';
import { useEventFilter } from '@/hooks/store/modules/event-filter';
import { useFlappingRules } from '@/hooks/store/modules/flapping-rules';
import { useResolveRules } from '@/hooks/store/modules/resolve-rules';
import { useIdleRules } from '@/hooks/store/modules/idle-rules';
import { useLinkRule } from '@/hooks/store/modules/link-rule';
import { useMetaAlarmRule } from '@/hooks/store/modules/meta-alarm-rule';
import { useSnmpRule } from '@/hooks/store/modules/snmp-rule';
import { useScenario } from '@/hooks/store/modules/scenario';
import { useDeclareTicketRule } from '@/hooks/store/modules/declare-ticket-rule';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';
import { usePlaylist } from '@/hooks/store/modules/playlist';

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
    pbehavior: {
      type: Boolean,
      default: false,
    },
    dynamicInfo: {
      type: Boolean,
      default: false,
    },
    eventFilter: {
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
    metaAlarmRule: {
      type: Boolean,
      default: false,
    },
    snmpRule: {
      type: Boolean,
      default: false,
    },
    scenario: {
      type: Boolean,
      default: false,
    },
    declareTicket: {
      type: Boolean,
      default: false,
    },
    playlist: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const modals = useModals();

    const {
      fetchDynamicInfosListWithPreviousParams,
      bulkEnableDynamicInfos,
      bulkDisableDynamicInfos,
      bulkRemoveDynamicInfos,
    } = useDynamicInfo();

    const {
      fetchEventFiltersListWithPreviousParams,
      bulkEnableEventFilters,
      bulkDisableEventFilters,
      bulkRemoveEventFilters,
    } = useEventFilter();

    const {
      fetchFlappingRulesListWithPreviousParams,
      bulkEnableFlappingRules,
      bulkDisableFlappingRules,
      bulkRemoveFlappingRules,
    } = useFlappingRules();

    const {
      fetchResolveRulesListWithPreviousParams,
      bulkEnableResolveRules,
      bulkDisableResolveRules,
      bulkRemoveResolveRules,
    } = useResolveRules();

    const {
      fetchIdleRulesListWithPreviousParams,
      bulkEnableIdleRules,
      bulkDisableIdleRules,
      bulkRemoveIdleRules,
    } = useIdleRules();

    const {
      fetchLinkRulesListWithPreviousParams,
      bulkEnableLinkRules,
      bulkDisableLinkRules,
      bulkRemoveLinkRules,
    } = useLinkRule();

    const {
      bulkEnableMetaAlarmRules,
      bulkDisableMetaAlarmRules,
      bulkRemoveMetaAlarmRules,
    } = useMetaAlarmRule();

    const {
      fetchSnmpRulesListWithPreviousParams,
      bulkEnableSnmpRules,
      bulkDisableSnmpRules,
      bulkRemoveSnmpRules,
    } = useSnmpRule();

    const {
      fetchScenariosListWithPreviousParams,
      bulkEnableScenarios,
      bulkDisableScenarios,
      bulkRemoveScenarios,
    } = useScenario();

    const {
      fetchDeclareTicketRulesListWithPreviousParams,
      bulkEnableDeclareTicketRules,
      bulkDisableDeclareTicketRules,
      bulkRemoveDeclareTicketRules,
    } = useDeclareTicketRule();

    const {
      fetchPbehaviorsListWithPreviousParams,
      bulkUpdatePbehaviors,
      bulkRemovePbehaviors,
    } = usePbehavior();

    const {
      fetchPlaylistsListWithPreviousParams,
      bulkEnablePlaylists,
      bulkDisablePlaylists,
      bulkRemovePlaylists,
    } = usePlaylist();

    const itemsIds = computed(() => mapIds(props.items));
    const enablableItems = computed(() => (
      props.pbehavior ? props.items.filter(({ editable }) => editable) : props.items
    ));

    const enablableItemsIds = computed(() => mapIds(enablableItems.value));

    const someOneEnable = computed(() => enablableItems.value.some(({ enabled }) => enabled));
    const someOneDisable = computed(() => enablableItems.value.some(({ enabled }) => !enabled));

    const config = computed(() => ({
      [props.pbehavior]: {
        afterSubmit: fetchPbehaviorsListWithPreviousParams,
        remove: bulkRemovePbehaviors,
        enable: bulkUpdatePbehaviors,
        disable: bulkUpdatePbehaviors,
        tooltipPrefix: 'pbehavior',
        exportProps: { pbehavior: true },
      },
      [props.dynamicInfo]: {
        afterSubmit: fetchDynamicInfosListWithPreviousParams,
        remove: bulkRemoveDynamicInfos,
        enable: bulkEnableDynamicInfos,
        disable: bulkDisableDynamicInfos,
        tooltipPrefix: 'dynamicInfo',
        exportProps: { dynamicInfo: true },
      },
      [props.eventFilter]: {
        afterSubmit: fetchEventFiltersListWithPreviousParams,
        remove: bulkRemoveEventFilters,
        enable: bulkEnableEventFilters,
        disable: bulkDisableEventFilters,
        tooltipPrefix: 'eventFilter',
        exportProps: { eventFilter: true },
      },
      [props.flappingRule]: {
        afterSubmit: fetchFlappingRulesListWithPreviousParams,
        remove: bulkRemoveFlappingRules,
        enable: bulkEnableFlappingRules,
        disable: bulkDisableFlappingRules,
        tooltipPrefix: 'flappingRule',
        exportProps: { flappingRule: true },
      },
      [props.resolveRule]: {
        afterSubmit: fetchResolveRulesListWithPreviousParams,
        remove: bulkRemoveResolveRules,
        enable: bulkEnableResolveRules,
        disable: bulkDisableResolveRules,
        tooltipPrefix: 'resolveRule',
        exportProps: { resolveRule: true },
      },
      [props.idleRule]: {
        afterSubmit: fetchIdleRulesListWithPreviousParams,
        remove: bulkRemoveIdleRules,
        enable: bulkEnableIdleRules,
        disable: bulkDisableIdleRules,
        tooltipPrefix: 'idleRule',
        exportProps: { idleRule: true },
      },
      [props.linkRule]: {
        afterSubmit: fetchLinkRulesListWithPreviousParams,
        remove: bulkRemoveLinkRules,
        enable: bulkEnableLinkRules,
        disable: bulkDisableLinkRules,
        tooltipPrefix: 'linkRule',
        exportProps: { linkRule: true },
      },
      [props.metaAlarmRule]: {
        afterSubmit: () => emit('refresh'), // TODO: rewrite fetching
        remove: bulkRemoveMetaAlarmRules,
        enable: bulkEnableMetaAlarmRules,
        disable: bulkDisableMetaAlarmRules,
        tooltipPrefix: 'metaAlarmRule',
        exportProps: { metaAlarmRule: true },
      },
      [props.snmpRule]: {
        afterSubmit: fetchSnmpRulesListWithPreviousParams,
        remove: bulkRemoveSnmpRules,
        enable: bulkEnableSnmpRules,
        disable: bulkDisableSnmpRules,
        tooltipPrefix: 'snmpRule',
      },
      [props.scenario]: {
        afterSubmit: fetchScenariosListWithPreviousParams,
        remove: bulkRemoveScenarios,
        enable: bulkEnableScenarios,
        disable: bulkDisableScenarios,
        tooltipPrefix: 'scenario',
        exportProps: { scenario: true },
      },
      [props.declareTicket]: {
        afterSubmit: fetchDeclareTicketRulesListWithPreviousParams,
        remove: bulkRemoveDeclareTicketRules,
        enable: bulkEnableDeclareTicketRules,
        disable: bulkDisableDeclareTicketRules,
        tooltipPrefix: 'declareTicket',
        exportProps: { declareTicket: true },
      },
      [props.playlist]: {
        afterSubmit: fetchPlaylistsListWithPreviousParams,
        remove: bulkRemovePlaylists,
        enable: bulkEnablePlaylists,
        disable: bulkDisablePlaylists,
        tooltipPrefix: 'playlist',
      },
    }.true ?? {}));

    const afterSubmit = async () => {
      emit('clear:items');

      return config.value.afterSubmit?.();
    };

    /**
     * Shows a confirmation modal for bulk remove operation.
     * On confirmation, removes selected items and refreshes the list.
     */
    const showRemoveModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await config.value.remove({ data: itemsIds.value });

          return afterSubmit();
        },
      },
    });

    /**
     * Shows a confirmation modal for bulk enable operation.
     * On confirmation, enables selected items and refreshes the list.
     */
    const showEnableModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await config.value.enable({ data: enablableItemsIds.value });

          return afterSubmit();
        },
      },
    });

    /**
     * Shows a confirmation modal for bulk disable operation.
     * On confirmation, disables selected items and refreshes the list.
     */
    const showDisableModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await config.value.disable({ data: enablableItemsIds.value });

          return afterSubmit();
        },
      },
    });

    return {
      config,
      itemsIds,
      someOneEnable,
      someOneDisable,

      showRemoveModal,
      showEnableModal,
      showDisableModal,
    };
  },
};
</script>
