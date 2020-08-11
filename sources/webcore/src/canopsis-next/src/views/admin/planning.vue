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
            v-tab(href="#resons") {{ $t('planning.tabs.reason') }}
            v-tab-item(value="resons")
              v-card-text
                | Reason
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
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';

import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue';
import FabButtons from '@/components/other/fab-buttons/fab-buttons.vue';

export const PLANNING_TABS = {
  types: 'types',
  resons: 'resons',
  datesOfExceptions: 'datesOfExceptions',
};

export default {
  components: { FabButtons, PlanningTypes },
  mixins: [rightsTechnicalPbehaviorTypesMixin, entitiesPbehaviorTypesMixin],
  data() {
    return {
      activeTab: PLANNING_TABS.types,
      typesParams: {},
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
        case PLANNING_TABS.datesOfExceptions:
      }
    },

    create() {
      switch (this.activeTab) {
        case PLANNING_TABS.types:
          this.showCreateTypeModal();
          break;
        case PLANNING_TABS.resons:
        case PLANNING_TABS.datesOfExceptions:
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
