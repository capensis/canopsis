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
            template(v-if="hasReadAnyPbehaviorExceptionAccess")
              v-tab(href="#exceptions") {{ $t('planning.tabs.exceptions') }}
              v-tab-item(value="exceptions")
                v-card-text
                  planning-exceptions(:params.sync="exceptionsParams")
    fab-buttons(@create="create", @refresh="refresh", :has-access="hasAccess")
      span {{ tooltipText }}
</template>

<script>
import { MODALS } from '@/constants';

import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import rightsTechnicalPbehaviorReasonsMixin from '@/mixins/rights/technical/pbehavior-reasons';
import rightsTechnicalPbehaviorExceptionsMixin from '@/mixins/rights/technical/pbehavior-exceptions';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';
import entitiesPbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';
import entitiesPbehaviorExceptionsMixin from '@/mixins/entities/pbehavior/exceptions';

import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue';
import PlanningReasons from '@/components/other/pbehavior/reasons/planning-reasons.vue';
import PlanningExceptions from '@/components/other/pbehavior/exceptions/planning-exceptions.vue';
import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';

export const PLANNING_TABS = {
  types: 'types',
  reasons: 'reasons',
  exceptions: 'exceptions',
};

export default {
  components: {
    PlanningExceptions,
    FabButtons,
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
  ],
  data() {
    return {
      activeTab: PLANNING_TABS.types,
      typesParams: {},
      reasonsParams: {},
      exceptionsParams: {},
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

    hasAccess() {
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
      this.fetchPbehaviorTypesList({ params: this.typesParams });
    },

    fetchReasonsList() {
      this.fetchPbehaviorReasonsList({ params: this.reasonsParams });
    },

    fetchExceptionsList() {
      this.fetchPbehaviorExceptionsList({ params: this.exceptionsParams });
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
