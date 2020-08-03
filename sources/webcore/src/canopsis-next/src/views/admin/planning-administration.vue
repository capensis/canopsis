<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.planningAdministration') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(fixed-tabs, slider-color="primary")
            v-tab {{ $t('planningAdministration.tabs.type') }}
            v-tab-item
              v-card-text
                | Type
            v-tab {{ $t('planningAdministration.tabs.reason') }}
            v-tab-item
              v-card-text
                | Reason
            v-tab {{ $t('planningAdministration.tabs.datesOfExceptions') }}
            v-tab-item
              v-card-text
                | Dates of exceptions
    planning-administration-fab-buttons
</template>

<script>
import PlanningAdministrationFabButtons from '@/components/other/planning-administration/admin/planning-administration-fab-buttons.vue';

export default {
  components: { PlanningAdministrationFabButtons },
};
</script>
