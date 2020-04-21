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

const { mapGetters } = createNamespacedHelpers('entities');

export default {
  components: { ViewTabRows },
  data() {
    return {
      playing: false,
      activeTabIndex: 0,
      playlist: {
        _id: 'id123',
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
    ...mapGetters(['getList']),

    tabs() {
      return this.getList(ENTITIES_TYPES.viewTab, this.playlist.tabs);
    },

    activeTab() {
      return this.tabs[this.activeTabIndex];
    },
  },
  mounted() {
    this.play();
  },
  beforeDestroy() {
    this.stopTabsChanging();
  },
  methods: {
    play() {
      if (this.playing) {
        return;
      }

      if (!this.playlist.fullscreen) {
        this.startTabsChanging();
        this.playing = true;

        return;
      }

      if (this.$refs.playlistWrapper) {
        this.$fullscreen.toggle(this.$refs.playlistWrapper, {
          fullscreenClass: 'full-screen',
          background: 'white',
          callback: (value) => {
            this.playing = value;

            if (value) {
              this.startTabsChanging();
            } else {
              this.stopTabsChanging();
            }
          },
        });
      }
    },
    startTabsChanging() {
      this.timer = setTimeout(this.changeTabTick, 10000);
    },

    changeTabTick() {
      this.activeTabIndex = this.activeTabIndex >= this.tabs.length - 1 ? 0 : this.activeTabIndex + 1;

      this.timer = setTimeout(this.changeTabTick, 10000);
    },

    stopTabsChanging() {
      clearTimeout(this.timer);
    },
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
