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
import { MODALS } from '@/constants';
import { compile } from '@/helpers/handlebars';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';

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
      return compile(this.config.template, { alarm: this.firstItem, entity: this.firstItem.entity });
    },
  },
};
</script>
