<template lang="pug">
v-card.ma-2.white--text(:class="format.color", tile, raised)
  div.pauseContainer(v-if="watcher.active_pb_some && !watcher.active_pb_all")
    v-icon.pauseIcon pause
  v-layout(justify-start, align-center)
    v-flex(xs2)
      component.ma-2(:is="format.icon")
    v-flex(xs10)
      p.watcherName {{ watcher.display_name }}
  v-layout
    v-flex(xs12)
      div.moreInfos.py-1(@click="showWeatherWatcherModal")
        v-layout(justify-center)
          div More infos
          v-icon.pl-1(color="white", small) arrow_forward
</template>

<script>
import { ENTITIES_STATES, MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';

import SunIcon from './icons/sun.vue';
import CloudySunIcon from './icons/cloudy-sun.vue';
import CloudIcon from './icons/cloud.vue';
import RainingCloudIcon from './icons/raining-cloud.vue';
import PauseIcon from './icons/pause.vue';

export default {
  components: {
    SunIcon,
    CloudySunIcon,
    CloudIcon,
    RainingCloudIcon,
    PauseIcon,
  },
  mixins: [modalMixin],
  props: {
    watcher: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      values: {
        [ENTITIES_STATES.ok]: {
          color: 'green darken-1',
          icon: SunIcon,
        },
        [ENTITIES_STATES.minor]: {
          color: 'yellow darken-1',
          icon: CloudySunIcon,
        },
        [ENTITIES_STATES.major]: {
          color: 'orange darken-1',
          icon: CloudIcon,
        },
        [ENTITIES_STATES.critical]: {
          color: 'red darken-1',
          icon: RainingCloudIcon,
        },
      },
    };
  },
  computed: {
    format() {
      const hasActivePb = this.watcher.active_pb_all || this.watcher.active_pb_watcher;

      if (hasActivePb) {
        return { icon: PauseIcon, color: 'grey lighten-1' };
      }

      return this.values[this.watcher.state.val];
    },
  },
  methods: {
    showWeatherWatcherModal() {
      this.showModal({
        name: MODALS.weatherWatcher,
        config: {
          watcherId: this.watcher.entity_id,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .iconContainer {
    font-size: 48px;
  }

  .pauseContainer {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    width: 25%;
    clip-path: polygon(100% 0 , 0 100%, 100% 100%);
    background-color: white;
    z-index: 1;
    position: absolute;
    right: 0;
  }

  .pauseIcon {
    z-index: 4;
    position: relative;
    top: 1em;
    left: 20%;
    color: black;
  }

  .watcherName {
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .moreInfos {
    z-index: 2;
    background-color: rgba(0,0,0,0.2);
    cursor: pointer;
  }
</style>
