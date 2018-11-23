<template lang="pug">
v-card.ma-2.white--text(:class="format.color", tile, raised)
  div.pauseContainer(v-if="watcher.active_pb_some && !watcher.active_pb_all")
    v-icon.pauseIcon pause
  v-layout(justify-start, align-center)
    v-flex(xs2)
      v-icon.px-3.py-2.white--text(size="5em") {{ format.icon }}
    v-flex(xs10)
      div.watcherName.pt-2(v-html="compiledTemplate")
  v-layout
    v-flex(xs12)
      div.moreInfos.py-1(@click="showWatcherModal")
        v-layout(justify-center)
          div {{ $t('weather.moreInfos') }}
          v-icon.pl-1(color="white", small) arrow_forward
</template>

<script>
import find from 'lodash/find';
import modalMixin from '@/mixins/modal/modal';
import compile from '@/helpers/handlebars';
import { WATCHER_PBEHAVIOR_COLOR, PBEHAVIOR_TYPES } from '@/constants';

export default {
  mixins: [modalMixin],
  props: {
    watcher: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
    },
    widget: {
      type: Object,
    },
  },
  computed: {
    isPaused() {
      return this.watcher.active_pb_all;
    },
    isPbehavior() {
      return this.watcher.active_pb_watcher;
    },
    format() {
      const iconsMap = {
        [this.$constants.ENTITIES_STATES.ok]: 'wb_sunny',
        [this.$constants.ENTITIES_STATES.minor]: 'person',
        [this.$constants.ENTITIES_STATES.major]: 'person',
        [this.$constants.ENTITIES_STATES.critical]: 'wb_cloudy',
      };

      if (!this.isPaused) {
        if (this.isPbehavior) {
          /* eslint-disable */
          const maintenancePbehavior =
            find(this.watcher.watcher_pbehavior, pbehavior => pbehavior.type_ === PBEHAVIOR_TYPES.maintenance);
          const outOfSurveillancePbehavior =
            find(this.watcher.watcher_pbehavior, pbehavior => pbehavior.type_ === PBEHAVIOR_TYPES.maintenance);
          /* eslint-enable */
          if (maintenancePbehavior) {
            return {
              icon: 'build',
              color: WATCHER_PBEHAVIOR_COLOR,
            };
          } else if (outOfSurveillancePbehavior) {
            return {
              icon: 'brightness_3',
              color: WATCHER_PBEHAVIOR_COLOR,
            };
          }

          return {
            icon: 'pause',
            color: WATCHER_PBEHAVIOR_COLOR,
          };
        }

        return {
          icon: iconsMap[this.watcher.state.val],
          color: this.$constants.WATCHER_STATES_COLORS[this.watcher.state.val],
        };
      }
      return {
        icon: 'pause',
        color: WATCHER_PBEHAVIOR_COLOR,
      };
    },
    compiledTemplate() {
      return compile(this.template, { watcher: this.watcher });
    },
  },
  methods: {
    showWatcherModal() {
      this.showModal({
        name: this.$constants.MODALS.watcher,
        config: {
          watcherId: this.watcher.entity_id,
          entityTemplate: this.widget.parameters.entityTemplate,
          modalTemplate: this.widget.parameters.modalTemplate,
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
