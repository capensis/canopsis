<template lang="pug">
  v-card
    v-card-text
      div(v-if="!config.widget.more_infos_popup")
        v-layout(justify-center)
          v-icon(color="info") infos
          p(class="ma-0") {{ $t('modals.moreInfos.defineATemplate') }}
      div(v-else, v-html="output")
</template>

<script>
import HandleBars from 'handlebars';

import modalInnerMixin from '@/mixins/modal/modal-inner';
import modalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import { MODALS } from '@/constants';

/**
 * Modal showing more infos on an alarm
 *
 * @prop {String} [template] - template to be shown on the modal
 */
export default {
  name: MODALS.moreInfos,
  mixins: [modalInnerMixin, modalInnerItemsMixin],
  computed: {
    output() {
      const output = HandleBars.compile(this.config.widget.more_infos_popup);
      const context = { alarm: this.items[0], entity: this.items[0].entity };
      return output(context);
    },
  },
};
</script>
