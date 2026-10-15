<template>
  <div class="more-infos">
    <c-compiled-template
      v-if="template"
      :template-id="templateId"
      :template="template"
      :context="templateContext"
      @select:tag="$emit('select:tag', $event)"
      @remove:tag="$emit('remove:tag', $event)"
    />
    <v-layout
      v-else
      justify-center
    >
      <v-icon color="info">
        infos
      </v-icon>
      <p class="ma-0">
        {{ $t('alarm.moreInfos.defineATemplate') }}
      </p>
    </v-layout>
  </div>
</template>

<script>
import { USER_PERMISSIONS } from '@/constants';

import { handlebarsLinksHelperCreator } from '@/mixins/handlebars/links-helper-creator';

export default {
  mixins: [
    handlebarsLinksHelperCreator(
      'alarm.links',
      USER_PERMISSIONS.business.alarmsList.actions.links,
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
    templateId: {
      type: String,
      required: false,
    },
    selectedTags: {
      type: Array,
      default: () => [],
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
