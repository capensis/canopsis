<template>
  <div>
    <c-page-header :name="GROUPED_USER_PERMISSIONS_KEYS.planning" />

    <v-card class="ma-4 mt-0">
      <v-tabs v-model="activeTab" color="primary" align-tabs="center">
        <v-tab v-if="hasReadAnyPbehaviorTypeAccess" :value="PLANNING_TABS.types">
          {{ t('pbehavior.tabs.type') }}
        </v-tab>
        <v-tab v-if="hasReadAnyPbehaviorReasonAccess" :value="PLANNING_TABS.reasons">
          {{ t('pbehavior.tabs.reason') }}
        </v-tab>
        <v-tab v-if="hasReadAnyPbehaviorExceptionAccess" :value="PLANNING_TABS.exceptions">
          {{ t('pbehavior.tabs.exceptions') }}
        </v-tab>
      </v-tabs>

      <v-card-text>
        <v-window v-model="activeTab">
          <v-window-item :value="PLANNING_TABS.types">
            <planning-types />
          </v-window-item>
          <v-window-item :value="PLANNING_TABS.reasons">
            <planning-reasons />
          </v-window-item>
          <v-window-item :value="PLANNING_TABS.exceptions">
            <planning-exceptions />
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <c-fab-expand-btn v-if="isExceptionTab" :has-access="hasCreateAccess" @refresh="refresh">
      <c-action-fab-btn
        :tooltip="t('modals.importPbehaviorException.title')"
        color="indigo"
        icon="upload_file"
        @click="showImportExceptionsModal"
      />
      <c-action-fab-btn
        :tooltip="t('modals.createPbehaviorException.create.title')"
        color="deep-purple"
        icon="event"
        @click="showCreateExceptionModal"
      />
    </c-fab-expand-btn>
    <c-fab-btn v-else :has-access="hasCreateAccess" @create="create" @refresh="refresh">
      <span>{{ tooltipText }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { computed, ref, watch } from 'vue'

import { GROUPED_USER_PERMISSIONS_KEYS, MODALS, PLANNING_TABS, USER_PERMISSIONS } from '@/constants'

import { useCRUDPermissions } from '@/hooks/auth'
import { useI18n } from '@/hooks/i18n'
import { useModals } from '@/hooks/modals'
import { usePbehaviorException } from '@/hooks/store/modules/pbehavior-exception'
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason'
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type'

import PlanningExceptions from '@/components/other/pbehavior/exceptions/planning-exceptions.vue'
import PlanningReasons from '@/components/other/pbehavior/reasons/planning-reasons.vue'
import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue'

export default {
  components: {
    PlanningExceptions,
    PlanningReasons,
    PlanningTypes,
  },

  setup() {
    const { t } = useI18n()
    const modals = useModals()

    const {
      hasCreateAccess: hasCreateAnyPbehaviorTypeAccess,
      hasReadAccess: hasReadAnyPbehaviorTypeAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.planningType)
    const {
      hasCreateAccess: hasCreateAnyPbehaviorReasonAccess,
      hasReadAccess: hasReadAnyPbehaviorReasonAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.planningReason)
    const {
      hasCreateAccess: hasCreateAnyPbehaviorExceptionAccess,
      hasReadAccess: hasReadAnyPbehaviorExceptionAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.planningExceptions)

    const { fetchPbehaviorTypesListWithPreviousParams, createPbehaviorType } = usePbehaviorType()
    const { fetchPbehaviorReasonsListWithPreviousParams, createPbehaviorReason } =
      usePbehaviorReason()
    const {
      fetchPbehaviorExceptionsListWithPreviousParams,
      createPbehaviorException,
      importPbehaviorException,
    } = usePbehaviorException()

    const availableTabs = computed(() =>
      [
        hasReadAnyPbehaviorTypeAccess.value && PLANNING_TABS.types,
        hasReadAnyPbehaviorReasonAccess.value && PLANNING_TABS.reasons,
        hasReadAnyPbehaviorExceptionAccess.value && PLANNING_TABS.exceptions,
      ].filter(Boolean),
    )

    const activeTab = ref(availableTabs.value[0] || PLANNING_TABS.types)
    const tooltipText = computed(
      () =>
        ({
          [PLANNING_TABS.types]: t('modals.createPbehaviorType.create.title'),
          [PLANNING_TABS.reasons]: t('modals.createPbehaviorReason.create.title'),
          [PLANNING_TABS.exceptions]: t('modals.createPbehaviorException.create.title'),
        })[activeTab.value],
    )
    const hasCreateAccess = computed(
      () =>
        ({
          [PLANNING_TABS.types]: hasCreateAnyPbehaviorTypeAccess.value,
          [PLANNING_TABS.reasons]: hasCreateAnyPbehaviorReasonAccess.value,
          [PLANNING_TABS.exceptions]: hasCreateAnyPbehaviorExceptionAccess.value,
        })[activeTab.value],
    )
    const isExceptionTab = computed(() => activeTab.value === PLANNING_TABS.exceptions)

    watch(availableTabs, (tabs) => {
      if (!tabs.includes(activeTab.value)) {
        activeTab.value = tabs[0] || PLANNING_TABS.types
      }
    })

    const refresh = () => {
      switch (activeTab.value) {
        case PLANNING_TABS.types:
          fetchPbehaviorTypesListWithPreviousParams()
          break
        case PLANNING_TABS.reasons:
          fetchPbehaviorReasonsListWithPreviousParams()
          break
        case PLANNING_TABS.exceptions:
          fetchPbehaviorExceptionsListWithPreviousParams()
          break
      }
    }

    const showCreateTypeModal = () =>
      modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          action: async (data) => {
            await createPbehaviorType({ data })
            await fetchPbehaviorTypesListWithPreviousParams()
          },
        },
      })

    const showCreateReasonModal = () =>
      modals.show({
        name: MODALS.createPbehaviorReason,
        config: {
          action: async (data) => {
            await createPbehaviorReason({ data })
            await fetchPbehaviorReasonsListWithPreviousParams()
          },
        },
      })

    const showCreateExceptionModal = () =>
      modals.show({
        name: MODALS.createPbehaviorException,
        config: {
          action: async (data) => {
            await createPbehaviorException({ data })
            await fetchPbehaviorExceptionsListWithPreviousParams()
          },
        },
      })

    const showImportExceptionsModal = () =>
      modals.show({
        name: MODALS.importPbehaviorException,
        config: {
          action: async (data) => {
            await importPbehaviorException({ data })
            await fetchPbehaviorExceptionsListWithPreviousParams()
          },
        },
      })

    const create = () => {
      switch (activeTab.value) {
        case PLANNING_TABS.types:
          showCreateTypeModal()
          break
        case PLANNING_TABS.reasons:
          showCreateReasonModal()
          break
      }
    }

    return {
      GROUPED_USER_PERMISSIONS_KEYS,
      PLANNING_TABS,
      activeTab,
      tooltipText,
      hasCreateAccess,
      isExceptionTab,
      hasReadAnyPbehaviorTypeAccess,
      hasReadAnyPbehaviorReasonAccess,
      hasReadAnyPbehaviorExceptionAccess,
      refresh,
      create,
      showCreateExceptionModal,
      showImportExceptionsModal,
      t,
    }
  },
}
</script>
