<template lang="pug">
  div.more-infos
    c-compiled-template(
      v-if="template",
      :template="template",
      :context="templateContext",
      @select:tag="$emit('select:tag', $event)",
      @clear:tag="$emit('clear:tag')"
    )
    v-layout(v-else, justify-center)
      v-icon(color="info") infos
      p.ma-0 {{ $t('alarm.moreInfos.defineATemplate') }}
</template>

<script>
import { USERS_PERMISSIONS } from '@/constants';

import { handlebarsLinksHelperCreator } from '@/mixins/handlebars/links-helper-creator';

export default {
  mixins: [
    handlebarsLinksHelperCreator(
      'alarm.links',
      USERS_PERMISSIONS.business.alarmsList.actions.links,
    ),
  ],
  props: {
    alarm: {
      type: Object,
      required: false,
    },
    template: {
      type: String,
      required: false,
    },
    selectedTag: {
      type: String,
      default: '',
    },
  },
  computed: {
    templateContext() {
      return {
        alarm: this.alarm,
        entity: this.alarm.entity,
      };
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
