<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsNumbers.title') }}
    v-container
      v-list-group(v-for="(stat, key) in stats", :key="key")
        v-list-tile(slot="activator") {{ key }}
        v-switch(v-model="yesNoMode", :label="$t('settings.statsNumbers.yesNoMode')")
        v-layout.d-flex(v-if="yesNoMode", justify-space-around)
          v-flex.my-2(xs5)
            div
              chrome(v-model="colors.yesNoMode.yes")
          v-flex.my-2(xs5)
            div
              chrome(v-model="colors.yesNoMode.no")
        v-layout.d-flex(v-else, justify-space-around, wrap)
          v-flex.my-2(xs5)
            div
              chrome(v-model="colors.normalMode.ok")
          v-flex.my-2(xs5)
            div
              chrome(v-model="colors.normalMode.major")
          v-flex.my-2(xs5)
            div
              chrome(v-model="colors.normalMode.critical")
</template>

<script>
import { Chrome } from 'vue-color';

export default {
  components: {
    Chrome,
  },
  props: {
    stats: {
      type: Object,
    },
  },
  data() {
    return {
      yesNoMode: false,
      colors: {
        yesNoMode: {
          yes: '',
          no: '',
        },
        normalMode: {
          ok: '',
          major: '',
          critical: '',
        },
      },
    };
  },
};
</script>

<style lang="scss" scoped>
  .vc-chrome {
    max-width: 100%;
  }
</style>

