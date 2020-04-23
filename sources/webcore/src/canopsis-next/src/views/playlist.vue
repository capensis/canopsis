<template lang="pug">
  div
    h1 Playlist title
    div.position-relative(ref="playlistWrapper")
      div.play-button-wrapper(v-if="!playing")
        v-btn.play-button(color="primary", large, @click="play")
          v-icon(large) play_arrow
      v-fade-transition(v-if="activeTab", mode="out-in")
        view-tab-rows(:tab="activeTab", :key="activeTab._id")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ENTITIES_TYPES } from '@/constants';

import ViewTabRows from '@/components/other/view/view-tab-rows.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('playlist');
const {
  mapActions: mapPlaylistPlayerActions,
  mapGetters: mapPlaylistPlayerGetters,
} = createNamespacedHelpers('playlistPlayer');

const { mapActions: mapGroupsActions } = createNamespacedHelpers('view/group');

const { mapGetters: mapEntitiesGetters } = createNamespacedHelpers('entities');

export default {
  components: { ViewTabRows },
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      pl: {
        _id: 'asd',
        name: 'Playlist #1',
        fullscreen: true,
        interval: {
          value: 10,
          unit: 'm',
        },
        tabs: [
          '875df4c2-027b-4549-8add-e20ed7ff7d4f', // Alarm default
          'view-tab_5a339b3a-0611-4d4c-b307-dc1b92aeb27d', // Meteo technic
          'view-tab_c02ae48e-7f0a-4ba4-9215-ba5662e1550c', // Meteo correct
        ],
      },
    };
  },
  computed: {
    ...mapEntitiesGetters(['getList']),

    ...mapGetters({
      getPlaylistItem: 'getItem',
    }),

    ...mapPlaylistPlayerGetters(['playing', 'activeTabIndex']),

    availableTabs() {
      return this.getList(this.pl.tabs, ENTITIES_TYPES.viewTab, true);
    },

    playlist() {
      return this.getPlaylistItem(this.id);
    },

    tabs() {
      return (this.playlist && this.playlist.tabs) || [];
    },

    activeTab() {
      return this.availableTabs[this.activeTabIndex];
    },
  },
  async mounted() {
    await this.fetchGroupsList();

    await this.fetchPlaylistItem({ id: this.id });

    this.setPlaylist({ playlist: this.playlist });
  },
  beforeDestroy() {
    this.setPlaylist({ playlist: {} });
  },
  methods: {
    ...mapActions({
      fetchPlaylistItem: 'fetchItem',
    }),

    ...mapGroupsActions({
      fetchGroupsList: 'fetchList',
    }),

    ...mapPlaylistPlayerActions(['setPlaylist', 'play']),
  },
};
</script>

<style lang="scss">
  .play-button-wrapper {
    position: absolute;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;
    z-index: 2;
    background: rgba(255, 255, 255, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
