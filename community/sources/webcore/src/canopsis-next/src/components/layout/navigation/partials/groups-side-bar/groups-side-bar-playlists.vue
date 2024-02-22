<template>
  <v-expansion-panels
    v-if="availablePlaylists.length"
    class="groups-side-bar-playlists"
    multiple
    dark
  >
    <v-expansion-panel class="secondary lighten-1">
      <v-expansion-panel-header>
        <div class="panel-header">
          <span>{{ $t(`pageHeaders.${$constants.USERS_PERMISSIONS.technical.playlist}.title`) }}</span>
        </div>
      </v-expansion-panel-header>
      <v-expansion-panel-content>
        <router-link
          v-for="playlist in availablePlaylists"
          :key="playlist._id"
          :title="playlist.name"
          :to="{ name: $constants.ROUTES_NAMES.playlist, params: { id: playlist._id } }"
        >
          <v-card
            class="secondary lighten-2 groups-side-bar-playlists__button"
            flat
          >
            <div class="pa-4">
              <v-layout
                align-center
                justify-space-between
              >
                <v-flex>
                  <v-layout align-center>
                    <span class="text-truncate pl-3">{{ playlist.name }}</span>
                  </v-layout>
                </v-flex>
              </v-layout>
            </div>
          </v-card>
        </router-link>
      </v-expansion-panel-content>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script>
import layoutNavigationGroupsBarPlaylistsMixin from '@/mixins/layout/navigation/groups-bar-playlists';

export default {
  mixins: [
    layoutNavigationGroupsBarPlaylistsMixin,
  ],
};
</script>

<style lang="scss" scoped>
.groups-side-bar-playlists {
  padding: 10px;
  box-shadow: none;

  &__button {
    border-radius: 0;
  }

  & ::v-deep .v-expansion-panel-header {
    height: 48px;
    min-height: 48px;
  }
}

a {
  color: inherit;
  text-decoration: none;

  &.router-link-active ::v-deep .v-card {
    background-color: #73879a !important;
    border-color: #73879a !important;
  }
}
</style>
