<template>
  <c-layout gap="1">
    <c-action-btn
      v-if="removable && activeConfig.remove"
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
  </c-layout>
</template>

<script>
import { computed } from 'vue'

import { MODALS, RESPONSE_STATUSES } from '@/constants'

import { pickIds } from '@/helpers/array'

import { useI18n } from '@/hooks/i18n'
import { useModals } from '@/hooks/modals'
import { usePopups } from '@/hooks/popups'
import { useDynamicInfo } from '@/hooks/store/modules/dynamic-info'
import { useDynamicInfoTemplates } from '@/hooks/store/modules/dynamic-info-templates'
import { useEventFilter } from '@/hooks/store/modules/event-filter'
import { useFlappingRules } from '@/hooks/store/modules/flapping-rules'
import { useIdleRules } from '@/hooks/store/modules/idle-rules'
import { usePbehaviorException } from '@/hooks/store/modules/pbehavior-exception'
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason'
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type'
import { usePlaylist } from '@/hooks/store/modules/playlist'
import { useResolveRules } from '@/hooks/store/modules/resolve-rules'
import { useScenario } from '@/hooks/store/modules/scenario'

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
  },

  emits: ['clear:items', 'refresh'],

  setup(props, { emit }) {
    const { te, t } = useI18n()
    const modals = useModals()
    const popups = usePopups()

    const { bulkEnableDynamicInfos, bulkDisableDynamicInfos, bulkRemoveDynamicInfos } =
      useDynamicInfo()

    const { bulkRemoveDynamicInfoTemplates } = useDynamicInfoTemplates()

    const { bulkEnableScenarios, bulkDisableScenarios, bulkRemoveScenarios } = useScenario()
    const { bulkEnableEventFilters, bulkDisableEventFilters, bulkRemoveEventFilters } =
      useEventFilter()
    const { bulkEnableFlappingRules, bulkDisableFlappingRules, bulkRemoveFlappingRules } =
      useFlappingRules()
    const { bulkEnableResolveRules, bulkDisableResolveRules, bulkRemoveResolveRules } =
      useResolveRules()
    const { bulkEnableIdleRules, bulkDisableIdleRules, bulkRemoveIdleRules } = useIdleRules()
    const { bulkEnablePlaylist, bulkDisablePlaylist, bulkRemovePlaylist } = usePlaylist()
    const { bulkHidePbehaviorTypes, bulkUnhidePbehaviorTypes, bulkRemovePbehaviorTypes } =
      usePbehaviorType()
    const { bulkHidePbehaviorReasons, bulkUnhidePbehaviorReasons, bulkRemovePbehaviorReasons } =
      usePbehaviorReason()
    const {
      bulkHidePbehaviorExceptions,
      bulkUnhidePbehaviorExceptions,
      bulkRemovePbehaviorExceptions,
    } = usePbehaviorException()

    const someOneEnable = computed(() =>
      props.items.some(({ enabled, hidden }) => enabled || hidden === false),
    )
    const someOneDisable = computed(() =>
      props.items.some(({ enabled, hidden }) => enabled === false || hidden === true),
    )

    const activeConfig = computed(() => {
      if (props.dynamicInfo) {
        return {
          remove: bulkRemoveDynamicInfos,
          enable: bulkEnableDynamicInfos,
          disable: bulkDisableDynamicInfos,
          tooltipPrefix: 'dynamicInfo',
        }
      }

      if (props.dynamicInfoTemplate) {
        return {
          remove: bulkRemoveDynamicInfoTemplates,
          tooltipPrefix: 'dynamicInfo.templatesList',
        }
      }

      if (props.scenario) {
        return {
          remove: bulkRemoveScenarios,
          enable: bulkEnableScenarios,
          disable: bulkDisableScenarios,
          tooltipPrefix: 'scenario',
        }
      }

      if (props.eventFilter) {
        return {
          remove: bulkRemoveEventFilters,
          enable: bulkEnableEventFilters,
          disable: bulkDisableEventFilters,
          tooltipPrefix: 'eventFilter',
        }
      }

      if (props.flappingRule) {
        return {
          remove: bulkRemoveFlappingRules,
          enable: bulkEnableFlappingRules,
          disable: bulkDisableFlappingRules,
          tooltipPrefix: 'flappingRule',
        }
      }

      if (props.resolveRule) {
        return {
          remove: bulkRemoveResolveRules,
          enable: bulkEnableResolveRules,
          disable: bulkDisableResolveRules,
          tooltipPrefix: 'resolveRule',
        }
      }

      if (props.idleRule) {
        return {
          remove: bulkRemoveIdleRules,
          enable: bulkEnableIdleRules,
          disable: bulkDisableIdleRules,
          tooltipPrefix: 'idleRules',
        }
      }

      if (props.playlist) {
        return {
          remove: bulkRemovePlaylist,
          enable: bulkEnablePlaylist,
          disable: bulkDisablePlaylist,
          tooltipPrefix: 'playlist',
        }
      }

      if (props.pbehaviorType) {
        return {
          remove: bulkRemovePbehaviorTypes,
          enable: bulkUnhidePbehaviorTypes,
          disable: bulkHidePbehaviorTypes,
          tooltipPrefix: 'pbehavior',
        }
      }

      if (props.pbehaviorReason) {
        return {
          remove: bulkRemovePbehaviorReasons,
          enable: bulkUnhidePbehaviorReasons,
          disable: bulkHidePbehaviorReasons,
          tooltipPrefix: 'pbehavior',
        }
      }

      if (props.pbehaviorException) {
        return {
          remove: bulkRemovePbehaviorExceptions,
          enable: bulkUnhidePbehaviorExceptions,
          disable: bulkHidePbehaviorExceptions,
          tooltipPrefix: 'pbehavior',
        }
      }

      return {}
    })

    const removeTooltip = computed(() => {
      const massRemoveTooltipKey = `${activeConfig.value.tooltipPrefix}.massRemove`

      return te(massRemoveTooltipKey) ? t(massRemoveTooltipKey) : ''
    })

    const afterSubmit = async () => {
      emit('clear:items')
      emit('refresh')
    }

    const showErrorPopups = (response = [], messageKey) => {
      response.forEach(({ status, message, error }, index) => {
        if (status !== RESPONSE_STATUSES.success) {
          const ruleName = props.items[index]?.name
          const text = te(messageKey) ? t(messageKey, { name: ruleName }) : message || error

          popups.error({ text })
        }
      })
    }

    const showRemoveModal = () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const response = await activeConfig.value.remove({
              data: pickIds(props.items, props.itemId),
            })

            showErrorPopups(response, `${activeConfig.value.tooltipPrefix}.removeForbidden`)

            return afterSubmit()
          },
        },
      })
    }

    const showEnableModal = () =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await activeConfig.value.enable({ data: pickIds(props.items, props.itemId) })

            return afterSubmit()
          },
        },
      })

    const showDisableModal = () =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const response = await activeConfig.value.disable({
              data: pickIds(props.items, props.itemId),
            })

            showErrorPopups(response, `${activeConfig.value.tooltipPrefix}.disableForbidden`)

            return afterSubmit()
          },
        },
      })

    return {
      t,
      activeConfig,
      removeTooltip,
      someOneEnable,
      someOneDisable,
      showRemoveModal,
      showEnableModal,
      showDisableModal,
    }
  },
}
</script>
