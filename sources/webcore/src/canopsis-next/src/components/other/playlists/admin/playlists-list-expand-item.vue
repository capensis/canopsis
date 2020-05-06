<template lang="pug">
  div
</template>

<script>
import rightsEntitiesPlaylistTabMixin from '@/mixins/rights/entities/playlist-tab';

export default {
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
