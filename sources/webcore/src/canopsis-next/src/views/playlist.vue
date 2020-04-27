<template lang="pug">
  div
    h1 Playlist title
    portal(to="additional-top-bar-items")
      v-fade-transition
        v-toolbar-items.mr-2(v-if="!pending")
          span.playlist-timer.white--text.mr-2 {{ time | duration }}
          v-btn(dark, icon, @click="prevTab")
            v-icon skip_previous
          v-btn(v-if="playing", dark, icon, @click="pause")
            v-icon pause
          v-btn(v-else, dark, icon, @click="play")
            v-icon play_arrow
          v-btn(dark, icon, @click="nextTab")
            v-icon skip_next
          v-btn(dark, icon, @click="toggleFullScreenMode")
            v-icon fullscreen
    div.position-relative(ref="playlistWrapper")
      div.play-button-wrapper(v-if="!playing")
        v-btn.play-button(color="primary", large, @click="play")
          v-icon(large) play_arrow
      v-fade-transition(v-if="activeTab", mode="out-in")
        view-tab-rows(:tab="activeTab", :key="activeTab._id")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { ENTITIES_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';

import ViewTabRows from '@/components/other/view/view-tab-rows.vue';

const { mapActions } = createNamespacedHelpers('playlist');

const {
  mapActions: mapGroupsActions,
  mapGetters: mapGroupsGetters,
} = createNamespacedHelpers('view/group');

const { mapGetters: mapEntitiesGetters } = createNamespacedHelpers('entities');

export default {
  components: { ViewTabRows },
  mixins: [authMixin],
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
    ...mapEntitiesGetters(['getList']),

    ...mapGroupsGetters({
      groupsPending: 'pending',
    }),

    availableTabs() {
      const tabsIds = (this.playlist && this.playlist.tabs) || [];
      const tabs = this.getList(ENTITIES_TYPES.viewTab, tabsIds, true);

      return tabs.filter(tab => tab[SCHEMA_EMBEDDED_KEY].parents.some(parent => this.checkReadAccess(parent.id)));
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
    this.pending = false;
    this.time = this.playlist.interval.value;
  },
  beforeDestroy() {
    this.stopTimer();
  },
  methods: {
    ...mapActions({
      fetchPlaylistItemWithoutStore: 'fetchItemWithoutStore',
    }),

    ...mapGroupsActions({
      fetchGroupsList: 'fetchList',
    }),

    play() {
      this.playing = true;
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
      this.time = this.playlist.interval.value;

      this.stopTimer();
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
