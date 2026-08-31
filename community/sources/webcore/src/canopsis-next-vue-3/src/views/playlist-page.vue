<template>
  <div ref="playlistEl">
    <v-fade-transition mode="out-in">
      <c-progress-overlay v-if="pending" pending />

      <div v-else-if="playlist" class="playlist">
        <c-page-header v-if="!isFullscreenMode">
          {{ playlist.name }}
        </c-page-header>

        <Teleport v-if="mounted && additionalTopBarItemsTarget" :to="additionalTopBarItemsTarget">
          <v-fade-transition>
            <c-layout v-if="!pending" class="playlist__actions mr-2" align="center">
              <span class="playlist__timer text-white mr-2">
                {{ formattedTime }}
              </span>

              <v-btn
                :disabled="!activeTab"
                class="text-white"
                icon="skip_previous"
                variant="text"
                @click="prevTab"
              />

              <v-btn
                v-if="pausing || !played"
                :disabled="!activeTab"
                class="text-white"
                icon="play_arrow"
                variant="text"
                @click="play"
              />
              <v-btn
                v-else
                :disabled="!activeTab"
                class="text-white"
                icon="pause"
                variant="text"
                @click="pause"
              />

              <v-btn
                :disabled="!activeTab"
                class="text-white"
                icon="skip_next"
                variant="text"
                @click="nextTab"
              />

              <c-action-btn
                :disabled="!activeTab"
                :tooltip="t('common.fullscreen')"
                icon="fullscreen"
                color="white"
                @click="toggleFullScreenMode"
              />
            </c-layout>
          </v-fade-transition>
        </Teleport>

        <div class="position-relative playlist__tabs-wrapper">
          <div v-if="!played" class="playlist__play-button-wrapper">
            <v-btn
              color="primary"
              size="large"
              icon
              @click="play"
            >
              <v-icon icon="play_arrow" size="large" />
            </v-btn>
          </div>

          <v-fade-transition mode="out-in">
            <view-tab-widgets
              v-if="activeTab"
              :key="activeTab._id"
              :tab="activeTab"
              visible
            />
          </v-fade-transition>
        </div>
      </div>
    </v-fade-transition>
  </div>
</template>

<script>
import { computed, onBeforeUnmount, onMounted, provide, ref, watch } from 'vue';

import { TELEPORTS_NAMES, VIEW_SCREEN_MODES } from '@/constants';

import Observer from '@/services/observer';

import { convertDurationToString, toSeconds } from '@/helpers/date/duration';

import { useFullscreen } from '@/hooks/fullscreen';
import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useActiveView } from '@/hooks/store/modules/active-view';
import { usePlaylist } from '@/hooks/store/modules/playlist';
import { useView } from '@/hooks/store/modules/view';

import ViewTabWidgets from '@/components/other/view/view-tab-widgets.vue';

export default {
  components: {
    ViewTabWidgets,
  },

  props: {
    id: {
      type: String,
      required: true,
    },
    autoplay: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const mounted = ref(false);
    const pending = ref(false);
    const pausing = ref(false);
    const played = ref(false);
    const playlist = ref(null);
    const time = ref(0);
    const timer = ref(null);
    const activeTabIndex = ref(0);
    const availableTabs = ref([]);
    const playlistEl = ref(null);
    const isFullscreenMode = ref(false);
    const periodicRefresh = new Observer();

    provide('$periodicRefresh', periodicRefresh);

    const { t } = useI18n();
    const popups = usePopups();
    const fullscreen = useFullscreen();
    const { setActiveViewScreenMode, clearActiveView } = useActiveView();
    const { fetchPlaylistWithoutStore } = usePlaylist();
    const { fetchViewTab } = useView();

    const additionalTopBarItemsTarget = computed(() => document.querySelector(`[data-teleport-target="${TELEPORTS_NAMES.additionalTopBarItems}"]`)
        ? `[data-teleport-target="${TELEPORTS_NAMES.additionalTopBarItems}"]`
        : '',
    );
    const tabsIds = computed(() => playlist.value?.tabs_list ?? []);
    const activeTab = computed(() => availableTabs.value[activeTabIndex.value]);
    const formattedTime = computed(() => convertDurationToString(time.value));

    /**
     * Clears the active countdown timeout.
     */
    const stopTimer = () => {
      clearTimeout(timer.value);
      timer.value = null;
    };

    /**
     * Resets the remaining time from the playlist interval.
     */
    const initTime = () => {
      const { value, unit } = playlist.value?.interval ?? {};

      time.value = toSeconds(value ?? 0, unit);
    };

    /**
     * Schedules the next countdown tick after one second.
     */
    const startTimer = () => {
      stopTimer();
      timer.value = setTimeout(timerTick, 1000);
    };

    /**
     * Resets the countdown and resumes it when playback is active.
     */
    const restartTimer = () => {
      stopTimer();
      initTime();

      if (played.value && !pausing.value) {
        startTimer();
      }
    };

    /**
     * Switches to the next available tab and restarts the countdown.
     */
    const nextTab = () => {
      if (!availableTabs.value.length) {
        return;
      }

      const lastIndex = availableTabs.value.length - 1;

      activeTabIndex.value = activeTabIndex.value >= lastIndex ? 0 : activeTabIndex.value + 1;

      restartTimer();
    };

    /**
     * Decrements the remaining time and advances to the next tab when it reaches zero.
     */
    function timerTick() {
      time.value -= 1;

      if (time.value <= 0) {
        nextTab();

        return;
      }

      startTimer();
    }

    /**
     * Switches to the previous available tab and restarts the countdown.
     */
    const prevTab = () => {
      if (!availableTabs.value.length) {
        return;
      }

      const lastIndex = availableTabs.value.length - 1;

      activeTabIndex.value = activeTabIndex.value <= 0 ? lastIndex : activeTabIndex.value - 1;

      restartTimer();
    };

    /**
     * Toggles fullscreen mode for the playlist container.
     */
    const toggleFullScreenMode = () => {
      fullscreen.toggle(playlistEl.value, {
        fullscreenClass: 'full-screen',
        callback: value => (isFullscreenMode.value = value),
      });
    };

    /**
     * Starts playback, optionally entering fullscreen when autoplay is enabled.
     */
    const play = () => {
      played.value = true;
      pausing.value = false;

      if (playlist.value?.fullscreen && !isFullscreenMode.value && props.autoplay) {
        toggleFullScreenMode();
      }

      startTimer();
    };

    /**
     * Pauses playback and stops the countdown.
     */
    const pause = () => {
      pausing.value = true;
      stopTimer();
    };

    /**
     * Loads view tabs referenced by the playlist and resets the active tab index.
     *
     * @returns {Promise<void>}
     */
    const fetchAvailableTabs = async () => {
      availableTabs.value = (await Promise.all(tabsIds.value.map(id => fetchViewTab({ id })))).filter(
        Boolean,
      );
      activeTabIndex.value = 0;
    };

    /**
     * Loads the playlist, its tabs, and starts autoplay when enabled.
     *
     * @returns {Promise<void>}
     */
    const fetchPlaylist = async () => {
      try {
        pending.value = true;
        stopTimer();
        played.value = false;
        pausing.value = false;
        playlist.value = await fetchPlaylistWithoutStore({ id: props.id });

        await fetchAvailableTabs();
        initTime();

        if (props.autoplay) {
          play();
        }
      } catch (err) {
        console.error(err);

        popups.error({ text: err.description || t('errors.default') });
      } finally {
        pending.value = false;
      }
    };

    onMounted(() => {
      mounted.value = true;
      clearActiveView();
      setActiveViewScreenMode(VIEW_SCREEN_MODES.default);
      fetchPlaylist();
    });

    onBeforeUnmount(() => {
      stopTimer();
      periodicRefresh.unregisterAll();
      fullscreen.exit();
      clearActiveView();
    });
    watch(() => props.id, fetchPlaylist);

    return {
      activeTab,
      additionalTopBarItemsTarget,
      formattedTime,
      isFullscreenMode,
      mounted,
      pausing,
      pending,
      play,
      played,
      playlist,
      playlistEl,
      prevTab,
      nextTab,
      pause,
      t,
      toggleFullScreenMode,
    };
  },
};
</script>

<style lang="scss">
.playlist {
  &__tabs-wrapper {
    min-height: 60px;
  }

  &__play-button-wrapper {
    position: absolute;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;
    z-index: 3;
    background: rgba(255, 255, 255, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  &__actions {
    width: 310px;
    align-items: center;
  }

  &__timer {
    line-height: 48px;
  }
}
</style>
