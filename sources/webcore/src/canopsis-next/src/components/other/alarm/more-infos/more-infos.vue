<template lang="pug">
  .more-infos(data-test="moreInfosTemplateContent", v-if="!template")
    v-layout(justify-center)
      v-icon(color="info") infos
      p(class="ma-0") {{ $t('alarmList.moreInfos.defineATemplate') }}
  .more-infos(data-test="moreInfosContent", v-else, v-html="output")
</template>

<script>
import { compile } from '@/helpers/handlebars';

export default {
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
  computed: {
    output() {
      return compile(this.template, { alarm: this.alarm, entity: this.alarm.entity });
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
