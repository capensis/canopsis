<template lang="pug">
  div.position-relative
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import VRuntimeTemplate from 'v-runtime-template';

import { compile } from '@/helpers/handlebars';

export default {
  components: {
    VRuntimeTemplate,
  },
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.widget.parameters?.template ?? '');

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
};
</script>
