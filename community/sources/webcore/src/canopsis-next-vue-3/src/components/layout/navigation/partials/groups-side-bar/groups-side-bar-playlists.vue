<template>
  <v-expansion-panels v-if="availablePlaylists.length" class="groups-side-bar-playlists" multiple>
    <v-expansion-panel class="secondary lighten-1 groups-side-bar-playlists__panel" elevation="0">
      <v-expansion-panel-title>
        <div class="panel-header">
          <span>{{ $t(`pageHeaders.${USER_PERMISSIONS.technical.playlist}.title`) }}</span>
        </div>
      </v-expansion-panel-title>
      <v-expansion-panel-text>
        <router-link
          v-for="playlist in availablePlaylists"
          :key="playlist._id"
          :title="playlist.name"
          :to="{ name: ROUTES_NAMES.playlist, params: { id: playlist._id } }"
        >
          <v-card class="secondary lighten-2 groups-side-bar-playlists__button" elevation="0">
            <div class="pa-4">
              <c-layout align="center" justify="space-between">
                <c-layout align="center" class="min-width-0">
                  <span class="text-truncate pl-3">
                    {{ playlist.name }}
                  </span>
                </c-layout>
              </c-layout>
            </div>
          </v-card>
        </router-link>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script>
import { ROUTES_NAMES, USER_PERMISSIONS } from '@/constants';
import layoutNavigationGroupsBarPlaylistsMixin from '@/mixins/layout/navigation/groups-bar-playlists';

export default {
  mixins: [layoutNavigationGroupsBarPlaylistsMixin],

  data() {
    return {
      ROUTES_NAMES,
      USER_PERMISSIONS,
    };
  },
};
</script>

<style lang="scss" scoped>
.groups-side-bar-playlists {
  padding: 10px;
  box-shadow: none;

  &__button {
    border-radius: 0;
  }

  :deep(.v-expansion-panel-title) {
    height: 48px;
    min-height: 48px;
  }

  :deep(.v-expansion-panel-text__wrapper) {
    padding: 0;
  }
}

a {
  color: inherit;
  text-decoration: none;

  &.router-link-active :deep(.v-card) {
    background-color: #73879a !important;
    border-color: #73879a !important;
  }
}

.min-width-0 {
  min-width: 0;
}
</style>
