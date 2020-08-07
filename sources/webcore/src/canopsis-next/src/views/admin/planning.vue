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
                  planning-types
            v-tab {{ $t('planning.tabs.reason') }}
            v-tab-item
              v-card-text
                | Reason
            v-tab {{ $t('planning.tabs.datesOfExceptions') }}
            v-tab-item
              v-card-text
                | Dates of exceptions
</template>

<script>
import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';

import PlanningTypes from '@/components/other/pbehavior/types/planning-types.vue';

export default {
  components: { PlanningTypes },
  mixins: [rightsTechnicalPbehaviorTypesMixin],
};
</script>
