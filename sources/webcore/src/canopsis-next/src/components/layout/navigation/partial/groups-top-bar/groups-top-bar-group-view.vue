<template lang="pug">
  v-list-tile(:to="viewLink")
    v-list-tile-title
      span {{ view.title }}
      v-btn.edit-view-button(
        :data-test="`editViewButton-view-${view._id}`",
        v-show="hasViewEditButtonAccess",
        color="grey darken-2",
        depressed,
        small,
        icon,
        @click.prevent="showEditViewModal"
      )
        v-icon(small) edit
      v-btn.duplicate-view-button(
        :data-test="`copyViewButton-view-${view._id}`",
        v-show="isNavigationEditingMode",
        depressed,
        small,
        icon,
        color="grey darken-2",
        @click.prevent="showDuplicateViewModal"
      )
        v-icon(small) file_copy
</template>

<script>
import layoutNavigationGroupsBarGroupViewMixin from '@/mixins/layout/navigation/groups-bar-group-view';

export default {
  mixins: [layoutNavigationGroupsBarGroupViewMixin],
};
</script>
