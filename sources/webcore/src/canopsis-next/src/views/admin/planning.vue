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
              v-tab(href="#resons") {{ $t('planning.tabs.reason') }}
              v-tab-item(value="resons")
                v-card-text
                  planning-reasons(:params.sync="reasonsParams")
            v-tab(href="#datesOfExceptions") {{ $t('planning.tabs.datesOfExceptions') }}
            v-tab-item(value="datesOfExceptions")
              v-card-text
                | Dates of exceptions
    fab-buttons(@create="create", @refresh="refresh")
      span {{ tooltipText }}
</template>

<script>
import { MODALS } from '@/constants';

import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import rightsTechnicalPbehaviorReasonsMixin from '@/mixins/rights/technical/pbehavior-reasons';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';
import entitiesPbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';

import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue';
import PlanningReasons from '@/components/other/pbehavior/reasons/planning-reasons.vue';
import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';

export const PLANNING_TABS = {
  types: 'types',
  resons: 'resons',
  datesOfExceptions: 'datesOfExceptions',
};

export default {
  components: {
    FabButtons,
    PlanningTypes,
    PlanningReasons,
  },
  mixins: [
    rightsTechnicalPbehaviorTypesMixin,
    rightsTechnicalPbehaviorReasonsMixin,
    entitiesPbehaviorTypesMixin,
    entitiesPbehaviorReasonsMixin,
  ],
  data() {
    return {
      activeTab: PLANNING_TABS.types,
      typesParams: {},
      reasonsParams: {},
    };
  },
  computed: {
    tooltipText() {
      return {
        [PLANNING_TABS.types]: this.$t('modals.createPbehaviorType.title'),
        [PLANNING_TABS.resons]: this.$t('modals.createPbehaviorReason.title'),
        [PLANNING_TABS.datesOfExceptions]: this.$t('modals.createPbehaviorException.title'),
      }[this.activeTab];
    },
  },
  methods: {
    refresh() {
      switch (this.activeTab) {
        case PLANNING_TABS.types:
          this.fetchTypesList();
          break;
        case PLANNING_TABS.resons:
          this.fetchReasonsList();
          break;
        case PLANNING_TABS.datesOfExceptions:
      }
    },

    create() {
      switch (this.activeTab) {
        case PLANNING_TABS.types:
          this.showCreateTypeModal();
          break;
        case PLANNING_TABS.resons:
          this.showCreateReasonModal();
          break;
        case PLANNING_TABS.datesOfExceptions:
      }
    },

    fetchTypesList() {
      this.fetchPbehaviorTypesList({ params: this.typesParams });
    },

    fetchReasonsList() {
      this.fetchPbehaviorReasonsList({ params: this.reasonsParams });
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
  },
};
</script>
