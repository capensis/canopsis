<template>
  <v-list-item :to="viewLink" class="groups-top-bar-view">
    <v-list-item-title>
      <span>{{ view.title }}</span>
      <c-icon v-if="view.is_private" size="16">lock</c-icon>
      <v-btn
        v-show="(view.is_private || hasViewEditButtonAccess) && isNavigationEditingMode"
        class="ml-2"
        variant="text"
        size="small"
        icon
        @click.prevent="showEditViewModal"
      >
        <c-icon size="16">edit</c-icon>
      </v-btn>
      <v-btn
        v-show="(view.is_private || hasCreateAnyViewAccess) && isNavigationEditingMode"
        variant="text"
        size="small"
        icon
        @click.prevent="showDuplicateViewModal"
      >
        <c-icon size="16">file_copy</c-icon>
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

<style lang="scss" scoped>
.groups-top-bar-view :deep(.v-list-item-title) {
  display: flex;
  align-items: center;
  min-width: 180px;
  color: white !important;
}

.groups-top-bar-view :deep(.v-list-item__content) {
  color: white;
}
</style>
