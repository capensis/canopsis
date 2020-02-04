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
        v-show="isEditingMode",
        depressed,
        small,
        icon,
        color="grey darken-2",
        @click.prevent="showDuplicateViewModal"
      )
        v-icon(small) file_copy
</template>

<script>
import layoutNavigationViewItemMixin from '@/mixins/layout/navigation/view-item';

export default {
  mixins: [layoutNavigationViewItemMixin],
};
</script>
