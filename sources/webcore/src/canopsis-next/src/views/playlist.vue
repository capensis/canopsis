<template lang="pug">
  div
    v-fade-transition(mode="out-in")
      progress-overlay(v-if="pending", :pending="true")
      div(v-else-if="playlist")
        h2.text-xs-center.my-3.display-1.font-weight-medium {{ playlist.name }}
        portal(to="additional-top-bar-items")
          v-fade-transition
            v-toolbar-items.mr-2(v-if="!pending")
              span.playlist-timer.white--text.mr-2 {{ time | duration }}
              v-btn(:disabled="!activeTab", dark, icon, @click="prevTab")
                v-icon skip_previous
              v-btn(v-if="playing", :disabled="!activeTab", dark, icon, @click="pause")
                v-icon pause
              v-btn(v-else, :disabled="!activeTab", dark, icon, @click="play")
                v-icon play_arrow
              v-btn(:disabled="!activeTab", dark, icon, @click="nextTab")
                v-icon skip_next
              v-btn(:disabled="!activeTab", dark, icon, @click="toggleFullScreenMode")
                v-icon fullscreen
        div.white.position-relative(ref="playlistWrapper", v-if="activeTab")
          div.play-button-wrapper(v-if="!playing")
            v-btn.play-button(color="primary", large, @click="play")
              v-icon(large) play_arrow
          v-fade-transition(mode="out-in")
            view-tab-rows(:tab="activeTab", :key="activeTab._id")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { getSecondsByUnit } from '@/helpers/time';

import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import rightsEntitiesPlaylistTabMixin from '@/mixins/rights/entities/playlist-tab';

import ViewTabRows from '@/components/other/view/view-tab-rows.vue';
import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

const { mapActions } = createNamespacedHelpers('playlist');

export default {
  components: { ViewTabRows, ProgressOverlay },
  mixins: [entitiesViewGroupMixin, rightsEntitiesPlaylistTabMixin],
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      time: 0,
      pending: false,
      playing: false,
      playlist: null,
      activeTabIndex: 0,
    };
  },
  computed: {
    availableTabs() {
      const tabsIds = (this.playlist && this.playlist.tabs_list) || [];

      return this.getAvailableTabsByIds(tabsIds);
    },

    activeTab() {
      return this.availableTabs[this.activeTabIndex];
    },
  },
  async mounted() {
    this.pending = true;

    if (!this.groupsPending) {
      await this.fetchGroupsList();
    }

    this.playlist = await this.fetchPlaylistItemWithoutStore({ id: this.id });
    this.initTime();

    this.pending = false;
  },
  beforeDestroy() {
    this.stopTimer();
  },
  methods: {
    ...mapActions({
      fetchPlaylistItemWithoutStore: 'fetchItemWithoutStore',
    }),

    initTime() {
      const { interval, unit } = this.playlist.interval;

      this.time = getSecondsByUnit(interval, unit);
    },

    play() {
      this.playing = true;

      if (this.playlist.fullscreen) {
        this.toggleFullScreenMode();
      }

      this.startTimer();
    },

    pause() {
      this.playing = false;
      this.stopTimer();
    },

    prevTab() {
      if (this.availableTabs.length) {
        const lastIndex = this.availableTabs.length - 1;

        this.activeTabIndex = this.activeTabIndex <= 0 ? lastIndex : this.activeTabIndex - 1;
        this.time = this.playlist.interval.value;

        this.restartTimer();
      }
    },

    nextTab() {
      if (this.availableTabs.length) {
        const lastIndex = this.availableTabs.length - 1;
        this.activeTabIndex = this.activeTabIndex >= lastIndex ? 0 : this.activeTabIndex + 1;

        this.restartTimer();
      }
    },

    timerTick() {
      this.time -= 1;


      if (this.time <= 0) {
        return this.nextTab();
      }

      return this.startTimer();
    },

    startTimer() {
      this.timer = setTimeout(this.timerTick, 1000);
    },

    stopTimer() {
      clearTimeout(this.timer);
    },

    restartTimer() {
      this.stopTimer();
      this.initTime();
      this.startTimer();
    },

    toggleFullScreenMode() {
      this.$fullscreen.toggle(this.$refs.playlistWrapper, {
        fullscreenClass: 'full-screen',
        background: 'white',
      });
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

  .playlist-timer {
    line-height: 48px;
  }
</style>
