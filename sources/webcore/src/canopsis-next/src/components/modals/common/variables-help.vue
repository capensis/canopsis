<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.variablesHelp.variables') }}
    template(slot="text")
      v-treeview(
        :items="config.variables",
        item-key="name"
      )
        template(slot="prepend", slot-scope="props", v-if="props.item.isArray")
          div.caption.font-italic (Array)
        template(slot="label", slot-scope="props")
          div {{ props.item.name }}
            span.pl-1(v-if="props.leaf") :
              c-ellipsis.pl-1.d-inline-block.grey--text.body-1.pre-wrap(:text="String(props.item.value)")
            span.pl-1(v-else-if="!props.leaf && !(props.item.children && props.item.children.length)") :
              .pl-1.d-inline-block.grey--text.text--darken-1.body-1.font-italic {{ $t('common.emptyObject') }}
        template(slot="append", slot-scope="props", v-if="props.leaf")
          v-tooltip(left)
            v-btn(
              v-clipboard:copy="props.item.path",
              v-clipboard:success="() => $popups.success({ text: $t('success.pathCopied') })",
              v-clipboard:error="() => $popups.error({ text: $t('errors.default') })",
              slot="activator",
              small,
              icon
            )
              v-icon content_copy
            span {{ $t('modals.variablesHelp.copyToClipboard') }}
</template>

<script>
import { MODALS } from '@/constants';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.variablesHelp,
  components: { ModalWrapper },
};
</script>
