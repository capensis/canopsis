<template lang="pug">
  div
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import EntityLinks from './entity-links.vue';

export default {
  components: {
    VRuntimeTemplate,
    EntityLinks,
  },
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
  beforeCreate() {
    registerHelper('links', ({ hash }) => {
      const category = hash.category ? `'${hash.category}'` : null;

      return new Handlebars.SafeString(`
        <div>
          <entity-links
          :links="entity.linklist"
          :category="${category}"
          ></entity-links>
        </div>
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('links');
  },
};
</script>
