<template>
  <v-list-item :to="viewLink">
    <v-list-item-title>
      <span>{{ view.title }}</span>
      <v-btn
        class="edit-view-button"
        :data-test="`editViewButton-view-${view._id}`"
        v-show="hasViewEditButtonAccess"
        color="grey darken-2"
        depressed
        small
        icon
        @click.prevent="showEditViewModal"
      >
        <v-icon small>
          edit
        </v-icon>
      </v-btn>
      <v-btn
        class="duplicate-view-button"
        :data-test="`copyViewButton-view-${view._id}`"
        v-show="isNavigationEditingMode"
        depressed
        small
        icon
        color="grey darken-2"
        @click.prevent="showDuplicateViewModal"
      >
        <v-icon small>
          file_copy
        </v-icon>
      </v-btn>
    </v-list-item-title>
  </v-list-item>
</template>

<script>
import layoutNavigationGroupsBarGroupViewMixin from '@/mixins/layout/navigation/groups-bar-group-view';

export default {
  mixins: [layoutNavigationGroupsBarGroupViewMixin],
};
</script>
