<template lang="pug">
  v-card(dark)
    v-card-title.primary.pa-2.white--text
      v-layout(justify-space-between, align-center)
        h4 {{ $t('alarmList.infoPopup') }}
        v-btn.ma-0.ml-3(
          color="white",
          icon,
          small,
          @click="$emit('close')"
        )
          v-icon(small, color="error") close
    v-fade-transition
      v-card-text.pa-2
        v-runtime-template(:template="popupTextContent")
</template>

<script>
import VRuntimeTemplate from 'v-runtime-template';

import { compile } from '@/helpers/handlebars';

export default {
  components: { VRuntimeTemplate },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
      default: '',
    },
  },
  asyncComputed: {
    popupTextContent: {
      lazy: true,

      async get() {
        const context = { alarm: this.alarm, entity: this.alarm.entity || {} };
        const compiledTemplate = await compile(this.template, context);

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
};
</script>
