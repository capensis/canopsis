<template>
  <v-list-item :to="viewLink">
    <v-list-item-title>
      <span>{{ view.title }}</span>
      <v-icon
        v-if="view.is_private"
        color="white"
        small
      >
        lock
      </v-icon>
      <v-btn
        class="edit-view-button"
        v-show="view.is_private || hasViewEditButtonAccess"
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
        v-show="view.is_private || isNavigationEditingMode"
        depressed
        small
        icon
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
