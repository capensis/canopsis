<template lang="pug">
  v-tabs(color="secondary lighten-1", dark, centered, slider-color="primary")
    v-tab {{ $tc('common.information') }}
    v-tab-item
      v-layout.pa-3
        v-flex(xs12)
          v-card.pa-3
            v-layout(column)
              draggable-playlist-tabs(:tabs="availableTabs", disabled)
</template>

<script>
import { permissionsEntitiesPlaylistTabMixin } from '@/mixins/permissions/entities/playlist-tab';

import DraggablePlaylistTabs from '../form/partials/draggable-playlist-tabs.vue';

export default {
  components: { DraggablePlaylistTabs },
  mixins: [permissionsEntitiesPlaylistTabMixin],
  props: {
    playlist: {
      type: Object,
      required: true,
    },
  },
  computed: {
    availableTabs() {
      const tabsIds = (this.playlist && this.playlist.tabs_list) || [];

      return this.getAvailableTabsByIds(tabsIds);
    },
  },
};
</script>
