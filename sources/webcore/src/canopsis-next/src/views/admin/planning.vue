<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.planningAdministration') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(v-model="activeTab", fixed-tabs, slider-color="primary")
            template(v-if="hasReadAnyPbehaviorTypeAccess")
              v-tab(href="#types") {{ $t('planning.tabs.type') }}
              v-tab-item(value="types")
                v-card-text
                  planning-types(:params.sync="typesParams")
            template(v-if="hasReadAnyPbehaviorReasonAccess")
              v-tab(href="#reasons") {{ $t('planning.tabs.reason') }}
              v-tab-item(value="reasons")
                v-card-text
                  planning-reasons(:params.sync="reasonsParams")
            template(v-if="hasReadAnyPbehaviorDateExceptionAccess")
              v-tab(href="#datesExceptions") {{ $t('planning.tabs.datesExceptions') }}
              v-tab-item(value="datesExceptions")
                v-card-text
                  planning-dates-exceptions(:params.sync="datesExceptionsParams")
    fab-buttons(@create="create", @refresh="refresh", :has-access="hasAccess")
      span {{ tooltipText }}
</template>

<script>
import { MODALS } from '@/constants';

import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import rightsTechnicalPbehaviorReasonsMixin from '@/mixins/rights/technical/pbehavior-reasons';
import rightsTechnicalPbehaviorDatesExceptionsMixin from '@/mixins/rights/technical/pbehavior-dates-exceptions';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';
import entitiesPbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';
import entitiesPbehaviorDatesExceptionsMixin from '@/mixins/entities/pbehavior/dates-exceptions';

import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue';
import PlanningReasons from '@/components/other/pbehavior/reasons/planning-reasons.vue';
import PlanningDatesExceptions from '@/components/other/pbehavior/dates-exceptions/planning-dates-exceptions.vue';
import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';

export const PLANNING_TABS = {
  types: 'types',
  reasons: 'reasons',
  datesExceptions: 'datesExceptions',
};

export default {
  components: {
    PlanningDatesExceptions,
    FabButtons,
    PlanningTypes,
    PlanningReasons,
  },
  mixins: [
    rightsTechnicalPbehaviorTypesMixin,
    rightsTechnicalPbehaviorReasonsMixin,
    rightsTechnicalPbehaviorDatesExceptionsMixin,
    entitiesPbehaviorTypesMixin,
    entitiesPbehaviorReasonsMixin,
    entitiesPbehaviorDatesExceptionsMixin,
  ],
  data() {
    return {
      activeTab: PLANNING_TABS.types,
      typesParams: {},
      reasonsParams: {},
      datesExceptionsParams: {},
    };
  },
  computed: {
    tooltipText() {
      return {
        [PLANNING_TABS.types]: this.$t('modals.createPbehaviorType.title'),
        [PLANNING_TABS.reasons]: this.$t('modals.createPbehaviorReason.title'),
        [PLANNING_TABS.datesExceptions]: this.$t('modals.createPbehaviorDateException.title'),
      }[this.activeTab];
    },

    hasAccess() {
      return {
        [PLANNING_TABS.types]: this.hasCreateAnyPbehaviorTypeAccess,
        [PLANNING_TABS.reasons]: this.hasCreateAnyPbehaviorReasonAccess,
        [PLANNING_TABS.datesExceptions]: this.hasCreateAnyPbehaviorDateExceptionAccess,
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
        case PLANNING_TABS.datesExceptions:
          this.fetchDatesExceptionsList();
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
        case PLANNING_TABS.datesExceptions:
          this.showCreateDateExceptionModal();
          break;
      }
    },

    fetchTypesList() {
      this.fetchPbehaviorTypesList({ params: this.typesParams });
    },

    fetchReasonsList() {
      this.fetchPbehaviorReasonsList({ params: this.reasonsParams });
    },

    fetchDatesExceptionsList() {
      this.fetchPbehaviorDatesExceptionsList({ params: this.datesExceptionsParams });
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

    showCreateDateExceptionModal() {
      this.$modals.show({
        name: MODALS.createPbehaviorDateException,
        config: {
          action: async (data) => {
            await this.createPbehaviorDateException({ data });
            await this.fetchDatesExceptionsList();
          },
        },
      });
    },
  },
};
</script>
