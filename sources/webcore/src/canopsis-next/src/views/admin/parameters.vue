<template lang="pug">
  v-container
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-tabs(fixed-tabs, slider-color="primary")
            v-tab {{ $t('common.parameters') }}
            v-tab-item
              v-card-text
                user-interface-form(:disabled="!hasUpdateParametersAccess")
            v-tab {{ $t('common.import') }}/{{ $t('common.export') }}
            v-tab-item
              view-import-export
</template>

<script>
import rightsTechnicalParametersMixin from '@/mixins/rights/technical/parameters';

import UserInterfaceForm from '@/components/forms/user-interface.vue';
import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';
import ViewImportExport from '@/components/other/view/views-import-export.vue';

import layoutNavigationGroupsBarMixin from '@/mixins/layout/navigation/groups-bar';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';

export default {
  components: { ViewImportExport, UserInterfaceForm, GroupsSideBarGroup },
  mixins: [
    rightsTechnicalParametersMixin,
    layoutNavigationGroupsBarMixin,
    rightsEntitiesGroupMixin,
  ],
};
</script>
