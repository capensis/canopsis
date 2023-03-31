<template lang="pug">
  c-runtime-template(:template="compiledTemplate")
</template>

<script>
import VRuntimeTemplate from 'v-runtime-template';

import { USERS_PERMISSIONS } from '@/constants';

import { compile } from '@/helpers/handlebars';

import { handlebarsLinksHelperCreator } from '@/mixins/handlebars/links-helper-creator';

export default {
  components: {
    VRuntimeTemplate,
  },
  mixins: [
    handlebarsLinksHelperCreator(
      'entity.links',
      USERS_PERMISSIONS.business.serviceWeather.actions.entityLinks,
    ),
  ],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
      default: '',
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.template, { entity: this.entity });

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
};
</script>
