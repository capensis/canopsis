<template>
  <v-container>
    <c-page-header :name="GROUPED_USER_PERMISSIONS_KEYS.planning" />
    <v-layout wrap>
      <v-flex xs12>
        <v-card class="ma-2">
          <v-tabs
            v-model="activeTab"
            slider-color="primary"
            fixed-tabs
          >
            <template v-if="hasReadAnyPbehaviorTypeAccess">
              <v-tab :href="`#${PLANNING_TABS.types}`">
                {{ $t('pbehavior.tabs.type') }}
              </v-tab>
              <v-tab-item :value="PLANNING_TABS.types">
                <v-card-text>
                  <planning-types />
                </v-card-text>
              </v-tab-item>
            </template>
            <template v-if="hasReadAnyPbehaviorReasonAccess">
              <v-tab :href="`#${PLANNING_TABS.reasons}`">
                {{ $t('pbehavior.tabs.reason') }}
              </v-tab>
              <v-tab-item :value="PLANNING_TABS.reasons">
                <v-card-text>
                  <planning-reasons />
                </v-card-text>
              </v-tab-item>
            </template>
            <template v-if="hasReadAnyPbehaviorExceptionAccess">
              <v-tab :href="`#${PLANNING_TABS.exceptions}`">
                {{ $t('pbehavior.tabs.exceptions') }}
              </v-tab>
              <v-tab-item :value="PLANNING_TABS.exceptions">
                <v-card-text>
                  <planning-exceptions />
                </v-card-text>
              </v-tab-item>
            </template>
          </v-tabs>
        </v-card>
      </v-flex>
    </v-layout>
    <c-fab-expand-btn
      v-if="isExceptionTab"
      :has-access="hasCreateAccess"
      @refresh="refresh"
    >
      <c-action-fab-btn
        :tooltip="$t('modals.importPbehaviorException.title')"
        color="indigo"
        icon="upload_file"
        top
        @click="showImportExceptionsModal"
      />
      <c-action-fab-btn
        :tooltip="$t('modals.createPbehaviorException.title')"
        color="deep-purple"
        icon="event"
        top
        @click="showCreateExceptionModal"
      />
    </c-fab-expand-btn>
    <c-fab-btn
      v-else
      :has-access="hasCreateAccess"
      @create="create"
      @refresh="refresh"
    >
      <span>{{ tooltipText }}</span>
    </c-fab-btn>
  </v-container>
</template>

<script>
import { ref, computed } from 'vue';

import { GROUPED_USER_PERMISSIONS_KEYS, MODALS, PLANNING_TABS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason';
import { usePbehaviorException } from '@/hooks/store/modules/pbehavior-exception';

import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue';
import PlanningReasons from '@/components/other/pbehavior/reasons/planning-reasons.vue';
import PlanningExceptions from '@/components/other/pbehavior/exceptions/planning-exceptions.vue';

export default {
  components: {
    PlanningExceptions,
    PlanningTypes,
    PlanningReasons,
  },
  setup() {
    const modals = useModals();
    const { t } = useI18n();

    const {
      hasCreateAccess: hasCreateAnyPbehaviorTypeAccess,
      hasReadAccess: hasReadAnyPbehaviorTypeAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.planningType);

    const {
      hasCreateAccess: hasCreateAnyPbehaviorReasonAccess,
      hasReadAccess: hasReadAnyPbehaviorReasonAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.planningReason);

    const {
      hasCreateAccess: hasCreateAnyPbehaviorExceptionAccess,
      hasReadAccess: hasReadAnyPbehaviorExceptionAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.planningExceptions);

    const {
      fetchPbehaviorTypesListWithPreviousParams,
      createPbehaviorType,
    } = usePbehaviorType();

    const {
      fetchPbehaviorReasonsListWithPreviousParams,
      createPbehaviorReason,
    } = usePbehaviorReason();

    const {
      fetchPbehaviorExceptionsListWithPreviousParams,
      createPbehaviorException,
      importPbehaviorException,
    } = usePbehaviorException();

    const activeTab = ref(PLANNING_TABS.types);

    const tooltipText = computed(() => ({
      [PLANNING_TABS.types]: t('modals.createPbehaviorType.title'),
      [PLANNING_TABS.reasons]: t('modals.createPbehaviorReason.title'),
      [PLANNING_TABS.exceptions]: t('modals.createPbehaviorException.title'),
    }[activeTab.value]));

    const hasCreateAccess = computed(() => ({
      [PLANNING_TABS.types]: hasCreateAnyPbehaviorTypeAccess.value,
      [PLANNING_TABS.reasons]: hasCreateAnyPbehaviorReasonAccess.value,
      [PLANNING_TABS.exceptions]: hasCreateAnyPbehaviorExceptionAccess.value,
    }[activeTab.value]));

    const isExceptionTab = computed(() => activeTab.value === PLANNING_TABS.exceptions);

    /**
     * Shows the modal for creating a new pbehavior type.
     * After successful creation, refreshes the types list.
     */
    const showCreateTypeModal = () => {
      modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          action: async (data) => {
            await createPbehaviorType({ data });
            await fetchPbehaviorTypesListWithPreviousParams();
          },
        },
      });
    };

    /**
     * Shows the modal for creating a new pbehavior reason.
     * After successful creation, refreshes the reasons list.
     */
    const showCreateReasonModal = () => {
      modals.show({
        name: MODALS.createPbehaviorReason,
        config: {
          action: async (data) => {
            await createPbehaviorReason({ data });
            await fetchPbehaviorReasonsListWithPreviousParams();
          },
        },
      });
    };

    /**
     * Refreshes the list based on the active tab.
     */
    const refresh = () => {
      switch (activeTab.value) {
        case PLANNING_TABS.types:
          fetchPbehaviorTypesListWithPreviousParams();
          break;
        case PLANNING_TABS.reasons:
          fetchPbehaviorReasonsListWithPreviousParams();
          break;
        case PLANNING_TABS.exceptions:
          fetchPbehaviorExceptionsListWithPreviousParams();
          break;
      }
    };

    /**
     * Shows the create modal based on the active tab.
     */
    const create = () => {
      switch (activeTab.value) {
        case PLANNING_TABS.types:
          showCreateTypeModal();
          break;
        case PLANNING_TABS.reasons:
          showCreateReasonModal();
          break;
      }
    };

    /**
     * Shows the modal for creating a new pbehavior exception.
     * After successful creation, refreshes the exceptions list.
     */
    const showCreateExceptionModal = () => {
      modals.show({
        name: MODALS.createPbehaviorException,
        config: {
          action: async (data) => {
            await createPbehaviorException({ data });
            await fetchPbehaviorExceptionsListWithPreviousParams();
          },
        },
      });
    };

    /**
     * Shows the modal for importing pbehavior exceptions.
     * After successful import, refreshes the exceptions list.
     */
    const showImportExceptionsModal = () => {
      modals.show({
        name: MODALS.importPbehaviorException,
        config: {
          action: async (data) => {
            await importPbehaviorException({ data });
            await fetchPbehaviorExceptionsListWithPreviousParams();
          },
        },
      });
    };

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
    };
  },
};
</script>
