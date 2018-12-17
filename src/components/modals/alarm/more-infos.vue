<template lang="pug">
  v-card
    v-card-text
      div(v-if="!config.template")
        v-layout(justify-center)
          v-icon(color="info") infos
          p(class="ma-0") {{ $t('modals.moreInfos.defineATemplate') }}
      div(v-else, v-html="output")
</template>

<script>
import HandleBars from 'handlebars';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import { MODALS } from '@/constants';

/**
 * Modal showing more infos on an alarm
 *
 * @prop {String} [template] - template to be shown on the modal
 */
export default {
  name: MODALS.moreInfos,
  mixins: [modalInnerItemsMixin],
  computed: {
    output() {
      const output = HandleBars.compile(this.config.template);
      const context = { alarm: this.firstItem, entity: this.firstItem.entity };

      return output(context);
    },
  },
};
</script>
