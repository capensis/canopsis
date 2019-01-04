<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Variables
    v-card-text
      v-treeview(
      :items="config.variables",
      item-key="name"
      )
        template(slot="prepend", slot-scope="props", v-if="props.item.isArray")
          div.caption.font-italic (Array)
        template(slot="append", slot-scope="props", v-if="!props.item.children")
          v-tooltip(left)
            v-btn(@click="copyPathToClipBoard(props.item.path)", slot="activator", small, icon)
              v-icon file_copy
            span Copy to clipboard
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import popupMixin from '@/mixins/popup';

export default {
  name: MODALS.alarmsListVariablesHelp,
  mixins: [modalInnerItemsMixin, popupMixin],
  methods: {
    async copyPathToClipBoard(itemPath) {
      try {
        await this.$copyText(itemPath);
        this.addSuccessPopup({ text: this.$t('success.pathCopied') });
      } catch (err) {
        this.addErrorPopup({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
