<template lang="pug">
  v-card
    v-card-text
      div(v-if="!template")
        v-layout(justify-center)
          v-icon(color="info") infos
          p(class="ma-0") {{ $t('moreInfosModal.defineATemplate') }}
      div(v-else, v-html="output")
</template>

<script>
import HandleBars from 'handlebars';
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('modal');

export default {
  name: 'more-infos',
  props: {
    template: {
      type: String,
    },
  },
  computed: {

    ...mapGetters(['config']),

    output() {
      const output = HandleBars.compile(this.template);
      const context = { alarm: this.config.alarm.props, entity: this.config.alarm.props.entity };
      return output(context);
    },
  },
};
</script>
