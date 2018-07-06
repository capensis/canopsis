<template lang="pug">
  v-card
    v-card-text
      div(v-if="!template")
        v-layout(justify-center)
          v-icon(color="info") infos
          p(class="ma-0") {{ $t('modals.moreInfos.defineATemplate') }}
      div(v-else, v-html="output")
</template>

<script>
import HandleBars from 'handlebars';
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

const { mapGetters } = createNamespacedHelpers('modal');

/**
 * Modal showing more infos on an alarm
 *
 * @prop {String} [template] - template to be shown on the modal
 */
export default {
  name: MODALS.moreInfos,
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
