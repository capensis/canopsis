<template lang="pug">
  v-list-tile(:to="viewLink")
    v-list-tile-title
      span {{ view.title }}
      v-icon.ml-1(v-if="view.is_private", color="white", small) lock
      v-btn.edit-view-button(
        v-show="hasViewEditButtonAccess",
        color="grey darken-2",
        depressed,
        small,
        icon,
        @click.prevent="showEditViewModal"
      )
        v-icon(small) edit
      v-btn.duplicate-view-button(
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
