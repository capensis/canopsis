<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.planningAdministration') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(fixed-tabs, slider-color="primary")
            template(v-if="hasReadAnyPbehaviorTypeAccess")
              v-tab {{ $t('planning.tabs.type') }}
              v-tab-item
                v-card-text
                  planning-types-list
            v-tab {{ $t('planning.tabs.reason') }}
            v-tab-item
              v-card-text
                | Reason
            v-tab {{ $t('planning.tabs.datesOfExceptions') }}
            v-tab-item
              v-card-text
                | Dates of exceptions
    planning-fab-buttons
</template>

<script>
import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';

import PlanningFabButtons from '@/components/other/planning/admin/planning-fab-buttons.vue';
import PlanningTypesList from '@/components/other/pbehavior/types/pbehavior-types-list.vue';

export default {
  components: { PlanningTypesList, PlanningFabButtons },
  mixins: [rightsTechnicalPbehaviorTypesMixin],
};
</script>
