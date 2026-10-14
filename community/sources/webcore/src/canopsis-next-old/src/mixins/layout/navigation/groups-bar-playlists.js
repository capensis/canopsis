import { entitiesPlaylistMixin } from '@/mixins/entities/playlist';

export default {
  mixins: [entitiesPlaylistMixin],
  computed: {
    availablePlaylists() {
      return this.playlists.filter(({ enabled }) => enabled);
    },
  },
  mounted() {
    this.fetchPlaylistsList();
  },
};
