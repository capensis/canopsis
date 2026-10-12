<template>
  <v-tabs
    v-model="activeTab"
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.information') }}</v-tab>

    <v-tabs-items
      v-model="activeTab"
      mandatory
    >
      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card class="pa-3">
              <v-layout column>
                <draggable-playlist-tabs
                  :tabs="availableTabs"
                  disabled
                />
              </v-layout>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import { permissionsEntitiesPlaylistTabMixin } from '@/mixins/permissions/entities/playlist-tab';

import DraggablePlaylistTabs from '../form/fields/draggable-playlist-tabs.vue';

export default {
  components: { DraggablePlaylistTabs },
  mixins: [permissionsEntitiesPlaylistTabMixin],
  props: {
    playlist: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      activeTab: 0,
    };
  },
  computed: {
    availableTabs() {
      const tabsIds = (this.playlist && this.playlist.tabs_list) || [];

      return this.getAvailableTabsByIds(tabsIds);
    },
  },
};
</script>
