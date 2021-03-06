<template lang="pug">
  v-container
    c-the-page-header {{ $t('common.planningAdministration') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(v-model="activeTab", fixed-tabs, slider-color="primary")
            template(v-if="hasReadAnyPbehaviorTypeAccess")
              v-tab(:href="`#${$constants.PLANNING_TABS.types}`") {{ $t('planning.tabs.type') }}
              v-tab-item(:value="$constants.PLANNING_TABS.types")
                v-card-text
                  planning-types
            template(v-if="hasReadAnyPbehaviorReasonAccess")
              v-tab(:href="`#${$constants.PLANNING_TABS.reasons}`") {{ $t('planning.tabs.reason') }}
              v-tab-item(:value="$constants.PLANNING_TABS.reasons")
                v-card-text
                  planning-reasons
            template(v-if="hasReadAnyPbehaviorExceptionAccess")
              v-tab(:href="`#${$constants.PLANNING_TABS.exceptions}`") {{ $t('planning.tabs.exceptions') }}
              v-tab-item(:value="$constants.PLANNING_TABS.exceptions")
                v-card-text
                  planning-exceptions
    c-fab-btn(@create="create", @refresh="refresh", :has-access="hasCreateAccess")
      span {{ tooltipText }}
</template>

<script>
import { MODALS, PLANNING_TABS } from '@/constants';

import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import rightsTechnicalPbehaviorReasonsMixin from '@/mixins/rights/technical/pbehavior-reasons';
import rightsTechnicalPbehaviorExceptionsMixin from '@/mixins/rights/technical/pbehavior-exceptions';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';
import entitiesPbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';
import entitiesPbehaviorExceptionsMixin from '@/mixins/entities/pbehavior/exceptions';
import queryMixin from '@/mixins/query';

import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue';
import PlanningReasons from '@/components/other/pbehavior/reasons/planning-reasons.vue';
import PlanningExceptions from '@/components/other/pbehavior/exceptions/planning-exceptions.vue';

export default {
  components: {
    PlanningExceptions,
    PlanningTypes,
    PlanningReasons,
  },
  mixins: [
    rightsTechnicalPbehaviorTypesMixin,
    rightsTechnicalPbehaviorReasonsMixin,
    rightsTechnicalPbehaviorExceptionsMixin,
    entitiesPbehaviorTypesMixin,
    entitiesPbehaviorReasonsMixin,
    entitiesPbehaviorExceptionsMixin,
    queryMixin,
  ],
  data() {
    return {
      activeTab: PLANNING_TABS.types,
    };
  },
  computed: {
    tooltipText() {
      return {
        [PLANNING_TABS.types]: this.$t('modals.createPbehaviorType.title'),
        [PLANNING_TABS.reasons]: this.$t('modals.createPbehaviorReason.title'),
        [PLANNING_TABS.exceptions]: this.$t('modals.createPbehaviorException.title'),
      }[this.activeTab];
    },

    hasCreateAccess() {
      return {
        [PLANNING_TABS.types]: this.hasCreateAnyPbehaviorTypeAccess,
        [PLANNING_TABS.reasons]: this.hasCreateAnyPbehaviorReasonAccess,
        [PLANNING_TABS.exceptions]: this.hasCreateAnyPbehaviorExceptionAccess,
      }[this.activeTab];
    },
  },
  methods: {
    refresh() {
      switch (this.activeTab) {
        case PLANNING_TABS.types:
          this.fetchTypesList();
          break;
        case PLANNING_TABS.reasons:
          this.fetchReasonsList();
          break;
        case PLANNING_TABS.exceptions:
          this.fetchExceptionsList();
          break;
      }
    },

    create() {
      switch (this.activeTab) {
        case PLANNING_TABS.types:
          this.showCreateTypeModal();
          break;
        case PLANNING_TABS.reasons:
          this.showCreateReasonModal();
          break;
        case PLANNING_TABS.exceptions:
          this.showCreateExceptionModal();
          break;
      }
    },

    fetchTypesList() {
      this.fetchPbehaviorTypesListWithPreviousParams();
    },

    fetchReasonsList() {
      this.fetchPbehaviorReasonsListWithPreviousParams();
    },

    fetchExceptionsList() {
      this.fetchPbehaviorExceptionsListWithPreviousParams();
    },

    showCreateTypeModal() {
      this.$modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          action: async (data) => {
            await this.createPbehaviorType({ data });
            await this.fetchTypesList();
          },
        },
      });
    },

    showCreateReasonModal() {
      this.$modals.show({
        name: MODALS.createPbehaviorReason,
        config: {
          action: async (data) => {
            await this.createPbehaviorReason({ data });
            await this.fetchReasonsList();
          },
        },
      });
    },

    showCreateExceptionModal() {
      this.$modals.show({
        name: MODALS.createPbehaviorException,
        config: {
          action: async (data) => {
            await this.createPbehaviorException({ data });
            await this.fetchExceptionsList();
          },
        },
      });
    },
  },
};
</script>
