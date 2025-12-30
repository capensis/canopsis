<template>
  <v-layout>
    <c-action-btn
      v-if="removable && config.remove"
      :tooltip="config.removeTooltip"
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
    <c-action-btn
      v-if="unhideable && someOneHidden && config.unhide"
      :tooltip="$t(`${config.tooltipPrefix}.massUnhide`)"
      icon="check_circle"
      color="primary"
      @click="showUnhideModal"
    />
    <c-action-btn
      v-if="hideable && someOneVisible && config.hide"
      :tooltip="$t(`${config.tooltipPrefix}.massHide`)"
      icon="cancel"
      color="error"
      @click="showHideModal"
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

import { mapIds, pickIds } from '@/helpers/array';

import { useI18n } from '@/hooks/i18n';
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
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason';
import { usePbehaviorException } from '@/hooks/store/modules/pbehavior-exception';
import { usePlaylist } from '@/hooks/store/modules/playlist';
import { useMaps } from '@/hooks/store/modules/maps';

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
    pbehavior: {
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
    map: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { te, t } = useI18n();
    const modals = useModals();

    const {
      bulkEnableDynamicInfos,
      bulkDisableDynamicInfos,
      bulkRemoveDynamicInfos,
    } = useDynamicInfo();

    const {
      bulkEnableEventFilters,
      bulkDisableEventFilters,
      bulkRemoveEventFilters,
    } = useEventFilter();

    const {
      bulkEnableFlappingRules,
      bulkDisableFlappingRules,
      bulkRemoveFlappingRules,
    } = useFlappingRules();

    const {
      bulkEnableResolveRules,
      bulkDisableResolveRules,
      bulkRemoveResolveRules,
    } = useResolveRules();

    const {
      bulkEnableIdleRules,
      bulkDisableIdleRules,
      bulkRemoveIdleRules,
    } = useIdleRules();

    const {
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
      bulkEnableSnmpRules,
      bulkDisableSnmpRules,
      bulkRemoveSnmpRules,
    } = useSnmpRule();

    const {
      bulkEnableScenarios,
      bulkDisableScenarios,
      bulkRemoveScenarios,
    } = useScenario();

    const {
      bulkEnableDeclareTicketRules,
      bulkDisableDeclareTicketRules,
      bulkRemoveDeclareTicketRules,
    } = useDeclareTicketRule();

    const {
      bulkUpdatePbehaviors,
      bulkRemovePbehaviors,
    } = usePbehavior();

    const {
      bulkEnablePlaylists,
      bulkDisablePlaylists,
      bulkRemovePlaylists,
    } = usePlaylist();

    const {
      bulkHidePbehaviorTypes,
      bulkUnhidePbehaviorTypes,
      bulkRemovePbehaviorTypes,
    } = usePbehaviorType();

    const {
      bulkHidePbehaviorReasons,
      bulkUnhidePbehaviorReasons,
      bulkRemovePbehaviorReasons,
    } = usePbehaviorReason();

    const {
      bulkHidePbehaviorExceptions,
      bulkUnhidePbehaviorExceptions,
      bulkRemovePbehaviorExceptions,
    } = usePbehaviorException();

    const {
      bulkRemoveMaps,
    } = useMaps();

    const itemsIds = computed(() => mapIds(props.items));
    const enablableItems = computed(() => (
      props.pbehavior ? props.items.filter(({ editable }) => editable) : props.items
    ));

    const enablableItemsIds = computed(() => pickIds(enablableItems.value));

    const someOneEnable = computed(() => enablableItems.value.some(({ enabled }) => enabled));
    const someOneDisable = computed(() => enablableItems.value.some(({ enabled }) => !enabled));

    const hideableItems = computed(() => props.items);
    const hideableItemsIds = computed(() => pickIds(hideableItems.value));
    const someOneVisible = computed(() => hideableItems.value.some(({ hidden }) => !hidden));
    const someOneHidden = computed(() => hideableItems.value.some(({ hidden }) => hidden));

    const config = computed(() => {
      const activeConfig = {
        [props.pbehavior]: {
          remove: bulkRemovePbehaviors,
          enable: bulkUpdatePbehaviors,
          disable: bulkUpdatePbehaviors,
          tooltipPrefix: 'pbehavior',
          exportProps: { pbehavior: true },
        },
        [props.dynamicInfo]: {
          remove: bulkRemoveDynamicInfos,
          enable: bulkEnableDynamicInfos,
          disable: bulkDisableDynamicInfos,
          tooltipPrefix: 'dynamicInfo',
          exportProps: { dynamicInfo: true },
        },
        [props.eventFilter]: {
          remove: bulkRemoveEventFilters,
          enable: bulkEnableEventFilters,
          disable: bulkDisableEventFilters,
          tooltipPrefix: 'eventFilter',
          exportProps: { eventFilter: true },
        },
        [props.flappingRule]: {
          remove: bulkRemoveFlappingRules,
          enable: bulkEnableFlappingRules,
          disable: bulkDisableFlappingRules,
          tooltipPrefix: 'flappingRule',
          exportProps: { flappingRule: true },
        },
        [props.resolveRule]: {
          remove: bulkRemoveResolveRules,
          enable: bulkEnableResolveRules,
          disable: bulkDisableResolveRules,
          tooltipPrefix: 'resolveRule',
          exportProps: { resolveRule: true },
        },
        [props.idleRule]: {
          remove: bulkRemoveIdleRules,
          enable: bulkEnableIdleRules,
          disable: bulkDisableIdleRules,
          tooltipPrefix: 'idleRule',
          exportProps: { idleRule: true },
        },
        [props.linkRule]: {
          remove: bulkRemoveLinkRules,
          enable: bulkEnableLinkRules,
          disable: bulkDisableLinkRules,
          tooltipPrefix: 'linkRule',
          exportProps: { linkRule: true },
        },
        [props.metaAlarmRule]: {
          remove: bulkRemoveMetaAlarmRules,
          enable: bulkEnableMetaAlarmRules,
          disable: bulkDisableMetaAlarmRules,
          tooltipPrefix: 'metaAlarmRule',
          exportProps: { metaAlarmRule: true },
        },
        [props.snmpRule]: {
          remove: bulkRemoveSnmpRules,
          enable: bulkEnableSnmpRules,
          disable: bulkDisableSnmpRules,
          tooltipPrefix: 'snmpRule',
        },
        [props.scenario]: {
          remove: bulkRemoveScenarios,
          enable: bulkEnableScenarios,
          disable: bulkDisableScenarios,
          tooltipPrefix: 'scenario',
          exportProps: { scenario: true },
        },
        [props.declareTicket]: {
          remove: bulkRemoveDeclareTicketRules,
          enable: bulkEnableDeclareTicketRules,
          disable: bulkDisableDeclareTicketRules,
          tooltipPrefix: 'declareTicket',
          exportProps: { declareTicket: true },
        },
        [props.playlist]: {
          remove: bulkRemovePlaylists,
          enable: bulkEnablePlaylists,
          disable: bulkDisablePlaylists,
          tooltipPrefix: 'playlist',
        },
        [props.pbehaviorType]: {
          remove: bulkRemovePbehaviorTypes,
          hide: bulkHidePbehaviorTypes,
          unhide: bulkUnhidePbehaviorTypes,
          tooltipPrefix: 'pbehavior',
        },
        [props.pbehaviorReason]: {
          remove: bulkRemovePbehaviorReasons,
          hide: bulkHidePbehaviorReasons,
          unhide: bulkUnhidePbehaviorReasons,
          tooltipPrefix: 'pbehavior',
        },
        [props.pbehaviorException]: {
          remove: bulkRemovePbehaviorExceptions,
          hide: bulkHidePbehaviorExceptions,
          unhide: bulkUnhidePbehaviorExceptions,
          tooltipPrefix: 'pbehavior',
        },
        [props.map]: {
          remove: bulkRemoveMaps,
          tooltipPrefix: 'map',
        },
      }.true ?? {};

      const massRemoveTooltipKey = `${activeConfig.tooltipPrefix}.massRemove`;

      if (te(massRemoveTooltipKey)) {
        activeConfig.removeTooltip = t(massRemoveTooltipKey);
      }

      return activeConfig;
    });

    const afterSubmit = async () => {
      emit('clear:items');
      emit('refresh');
    };

    /**
     * Shows a confirmation modal for bulk remove operation.
     * On confirmation, removes selected items and refreshes the list.
     */
    const showRemoveModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await config.value.remove({ data: pickIds(props.items) });

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

    /**
     * Shows a confirmation modal for bulk hide operation.
     * On confirmation, hides selected items and refreshes the list.
     */
    const showHideModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await config.value.hide({ data: hideableItemsIds.value });

          return afterSubmit();
        },
      },
    });

    /**
     * Shows a confirmation modal for bulk unhide operation.
     * On confirmation, unhides selected items and refreshes the list.
     */
    const showUnhideModal = () => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await config.value.unhide({ data: hideableItemsIds.value });

          return afterSubmit();
        },
      },
    });

    return {
      config,
      itemsIds,
      someOneEnable,
      someOneDisable,
      someOneVisible,
      someOneHidden,

      showRemoveModal,
      showEnableModal,
      showDisableModal,
      showHideModal,
      showUnhideModal,
    };
  },
};
</script>
