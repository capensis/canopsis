<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.planningAdministration') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(v-model="activeTab", fixed-tabs, slider-color="primary")
            template(v-if="hasReadAnyPbehaviorTypeAccess")
              v-tab {{ $t('planning.tabs.type') }}
              v-tab-item
                v-card-text
                  planning-types(:params.sync="typesParams")
            v-tab {{ $t('planning.tabs.reason') }}
            v-tab-item
              v-card-text
                | Reason
            v-tab {{ $t('planning.tabs.datesOfExceptions') }}
            v-tab-item
              v-card-text
                | Dates of exceptions
    fab-buttons(@create="create", @refresh="refresh")
      span {{ tooltipText }}
</template>

<script>
import { MODALS, PLANNING_TABS_INDEXES } from '@/constants';

import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';

import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue';
import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';

export default {
  components: { FabButtons, PlanningTypes },
  mixins: [rightsTechnicalPbehaviorTypesMixin, entitiesPbehaviorTypesMixin],
  data() {
    return {
      activeTab: PLANNING_TABS_INDEXES.types,
      typesParams: {},
    };
  },
  computed: {
    tooltipText() {
      return {
        [PLANNING_TABS_INDEXES.types]: this.$t('modals.createPbehaviorType.title'),
        [PLANNING_TABS_INDEXES.resons]: this.$t('modals.createPbehaviorReason.title'),
        [PLANNING_TABS_INDEXES.datesOfExceptions]: this.$t('modals.createPbehaviorException.title'),
      }[this.activeTab];
    },
  },
  mounted() {
    this.fetchTypesList();
  },
  methods: {
    refresh() {
      switch (this.activeTab) {
        case PLANNING_TABS_INDEXES.types:
          this.fetchTypesList();
          break;
        case PLANNING_TABS_INDEXES.resons:
        case PLANNING_TABS_INDEXES.datesOfExceptions:
      }
    },

    create() {
      switch (this.activeTab) {
        case PLANNING_TABS_INDEXES.types:
          this.showCreateTypeModal();
          break;
        case PLANNING_TABS_INDEXES.resons:
        case PLANNING_TABS_INDEXES.datesOfExceptions:
      }
    },

    fetchTypesList() {
      this.fetchPbehaviorTypesList({ params: this.typesParams });
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
  },
};
</script>
