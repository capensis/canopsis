<template>
  <c-compiled-template
    :template="template"
    :context="templateContext"
    :template-props="templateProps"
  />
</template>

<script>
import { USERS_PERMISSIONS } from '@/constants';

import { handlebarsLinksHelperCreator } from '@/mixins/handlebars/links-helper-creator';

export default {
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
  computed: {
    templateContext() {
      return { entity: this.entity };
    },

    templateProps() {
      return { entity: this.entity };
    },
  },
};
</script>
