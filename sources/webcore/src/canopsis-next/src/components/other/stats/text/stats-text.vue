<template lang="pug">

</template>

<script>
import Handlebars from 'handlebars';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

export default {
  computed: {
    compiledTemplate() {
      return `<div>${compile(this.template)}</div>`;
    },
  },
  beforeCreate() {
    registerHelper('entities', ({ hash }) => {
      const entityNameField = hash.name || 'entity.name';

      return new Handlebars.SafeString(`
        <div class="mt-2" v-for="watcherEntity in watcherEntities">
          <watcher-entity
          :entity="watcherEntity"
          :template="entityTemplate"
          entityNameField="${entityNameField}"
          @addEvent="addEventToQueue"></watcher-entity>
        </div>
      `);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
};
</script>
