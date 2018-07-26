<template lang="pug">
  v-container(fluid)
    v-layout(wrap)
      v-flex(xs3)
        v-card(
          class=['white--text', 'green darken-1']
        )
          v-layout(justify-start, align-center)
            v-flex(xs2)
              raining-cloud.ma-2
            div Watcher Name
          v-layout
            v-flex(xs12)
              div.moreInfos.py-1
                v-layout(justify-center)
                  div More infos
                  v-icon.pl-1(color="white", small) arrow_forward

</template>

<script>
import RainingCloud from './icons/cloud.vue';

export default {
  components: { RainingCloud },
  data() {
    return {
      values: {
        ok: {
          color: 'green darken-1',
          icon: 'sun',
        },
        minor: {
          color: 'yellow darken-1',
          icon: 'cloudy sun',
        },
        major: {
          color: 'orange darken-1',
          icon: 'cloud',
        },
        critical: {
          color: 'red darken-1',
          icon: 'raining cloud',
        },
        pbehavior: {
          color: 'grey lighten-1',
          icon: 'pause',
        },
      },
    };
  },
};
</script>

<style lang="scss" scoped>
  .moreInfos {
    background-color: rgba(0,0,0,0.2);
  }
</style>

