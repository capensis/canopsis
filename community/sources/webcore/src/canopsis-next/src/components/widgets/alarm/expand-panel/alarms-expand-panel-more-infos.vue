<template lang="pug">
  div.more-infos()
    v-runtime-template(v-if="template", :template="compiledTemplate")
    v-layout(v-else, justify-center)
      v-icon(color="info") infos
      p.ma-0 {{ $t('alarm.moreInfos.defineATemplate') }}
</template>

<script>
import VRuntimeTemplate from 'v-runtime-template';

import { compile } from '@/helpers/handlebars';

export default {
  components: { VRuntimeTemplate },
  props: {
    alarm: {
      type: Object,
      required: false,
    },
    template: {
      type: String,
      required: false,
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.template, { alarm: this.alarm, entity: this.alarm.entity });

        return `<div>${compiledTemplate}</div>`;
      },
      lazy: true,
      default: '',
    },
  },
};
</script>

<style lang="scss" scoped>
.more-infos {
  width: 90%;
  margin: 0 auto;
}
</style>
