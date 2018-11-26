<template lang="pug">
v-card.white--text(:class="getItemClasses", tile)
  v-layout(justify-start, align-center)
    v-flex(xs2)
      component.ma-2(:is="format.icon")
    v-flex(xs10)
      div.watcherName.pt-2(v-html="compiledTemplate")
    v-btn.pauseIcon.white(v-if="watcher.active_pb_some && !watcher.active_pb_all", fab, icon, small)
      v-icon pause
</template>

<script>
import modalMixin from '@/mixins/modal/modal';
import compile from '@/helpers/handlebars';

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
    template: {
      type: String,
    },
    widget: {
      type: Object,
    },
  },
  computed: {
    format() {
      const hasActivePb = this.watcher.active_pb_all || this.watcher.active_pb_watcher;
      const iconsMap = {
        [this.$constants.ENTITIES_STATES.ok]: SunIcon,
        [this.$constants.ENTITIES_STATES.minor]: CloudySunIcon,
        [this.$constants.ENTITIES_STATES.major]: CloudIcon,
        [this.$constants.ENTITIES_STATES.critical]: RainingCloudIcon,
      };

      if (hasActivePb) {
        return { icon: PauseIcon, color: this.$constants.WATCHER_PBEHAVIOR_COLOR };
      }

      return {
        icon: iconsMap[this.watcher.state.val],
        color: this.$constants.WATCHER_STATES_COLORS[this.watcher.state.val],
      };
    },
    compiledTemplate() {
      return compile(this.template, { watcher: this.watcher });
    },
    getItemClasses() {
      return [
        this.format.color,
        `mt-${this.widget.parameters.margin.top}`,
        `mr-${this.widget.parameters.margin.right}`,
        `mb-${this.widget.parameters.margin.bottom}`,
        `ml-${this.widget.parameters.margin.left}`,
      ];
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
  .pauseIcon {
    position: absolute;
    right: 0;
    bottom: 0;
    cursor: inherit;
  }

  .iconContainer {
    font-size: 48px;
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
