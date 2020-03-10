<template lang="pug">
  .more-infos(v-if="!template", data-test="moreInfosTemplateContent")
    v-layout(justify-center)
      v-icon(color="info") infos
      p(class="ma-0") {{ $t('alarmList.moreInfos.defineATemplate') }}
  .more-infos(v-else, data-test="moreInfosContent")
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import VRuntimeTemplate from 'v-runtime-template';

import { compile } from '@/helpers/handlebars';

export default {
  components: { VRuntimeTemplate },
  props: {
    template: {
      type: String,
      required: false,
    },
    alarm: {
      type: Object,
      required: false,
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.template, { alarm: this.alarm, entity: this.alarm.entity });

        return `<div>${compiledTemplate}</div>`;
      },
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
