<template lang="pug">
  v-layout(row)
    v-flex
      draggable-playlist-tabs(:tabs="availableTabs", disabled)
</template>

<script>
import rightsEntitiesPlaylistTabMixin from '@/mixins/rights/entities/playlist-tab';

import DraggablePlaylistTabs from '../form/partials/draggable-playlist-tabs.vue';

export default {
  components: { DraggablePlaylistTabs },
  mixins: [rightsEntitiesPlaylistTabMixin],
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
