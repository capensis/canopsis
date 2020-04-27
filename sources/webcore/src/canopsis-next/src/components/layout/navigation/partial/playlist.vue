<template lang="pug">
  v-toolbar-items.mr-2(v-if="playlist")
    span.playlist-timer.white--text.mr-2(v-if="!isEmptyTabs") {{ time | duration }}
    v-btn(:disabled="isEmptyTabs", dark, icon, @click="prevTab")
      v-icon skip_previous
    v-btn(v-if="playing", :disabled="isEmptyTabs", dark, icon, @click="stop")
      v-icon stop
    v-btn(v-else, :disabled="isEmptyTabs", dark, icon, @click="play")
      v-icon play_arrow
    v-btn(dark, :disabled="isEmptyTabs", icon, @click="nextTab")
      v-icon skip_next
    v-btn(:disabled="isEmptyTabs", dark, icon)
      v-icon fullscreen
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('playlistPlayer');

const interval = 130;

export default {
  data() {
    return {
      time: interval,
    };
  },
  computed: {
    ...mapGetters(['playing', 'activeTabIndex', 'playlist']),

    isEmptyTabs() {
      const { tabs = [] } = this.playlist;

      return tabs.length === 0;
    },
  },
  watch: {
    playing(value, oldValue) {
      if (value !== oldValue) {
        if (value) {
          this.time = interval;
          this.startTimer();
        } else {
          this.time = interval;
          this.stopTimer();
        }
      }
    },
    activeTabIndex(value, oldValue) {
      if (value !== oldValue) {
        this.time = interval;
        this.stopTimer();
        this.startTimer();
      }
    },
  },
  methods: {
    ...mapActions(['play', 'stop', 'nextTab', 'prevTab']),

    timerTick() {
      this.time -= 1;

      if (this.time <= 0) {
        this.nextTab();

        return this.stopTimer();
      }

      return this.startTimer();
    },

    startTimer() {
      this.timer = setTimeout(this.timerTick, 1000);
    },

    stopTimer() {
      clearTimeout(this.timer);
    },
  },
};
</script>

<style lang="scss">
  .playlist-timer {
    line-height: 48px;
  }
</style>
