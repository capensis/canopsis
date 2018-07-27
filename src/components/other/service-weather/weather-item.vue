<template lang="pug">
v-card.ma-2(:class="values[watcher.state].color")
  v-layout(justify-start, align-center)
    v-flex(xs2)
      component.ma-2(:is="values[watcher.state].icon")
    v-flex(xs10)
      p.watcherName {{ watcher._id }}
  v-layout
    v-flex(xs12)
      div.moreInfos.py-1
        v-layout(justify-center)
          div More infos
          v-icon.pl-1(color="white", small) arrow_forward
</template>

<script>
import sun from './icons/sun.vue';
import cloudySun from './icons/cloudy-sun.vue';
import cloud from './icons/cloud.vue';
import rainingCloud from './icons/raining-cloud.vue';

export default {
  components: {
    sun,
    cloudySun,
    cloud,
    rainingCloud,
  },
  props: {
    watcher: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      values: {
        0: {
          color: 'white--text green darken-1',
          icon: sun,
        },
        1: {
          color: 'yellow darken-1',
          icon: cloudySun,
        },
        2: {
          color: 'orange darken-1',
          icon: cloud,
        },
        3: {
          color: 'red darken-1',
          icon: rainingCloud,
        },
        pbehavior: {
          color: 'grey lighten-1',
          icon: 'pause',
        },
      },
    };
  },
  methods: {
    test() {
      console.log(this.watcher);
    },
  },
};
</script>

<style lang="scss" scoped>
  .watcherName {
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .moreInfos {
    background-color: rgba(0,0,0,0.2);
  }
</style>
