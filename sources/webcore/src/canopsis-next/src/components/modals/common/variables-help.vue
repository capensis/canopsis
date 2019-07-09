<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.variablesHelp.variables') }}
    v-card-text
      v-treeview(
      :items="config.variables",
      item-key="name"
      )
        template(slot="prepend", slot-scope="props", v-if="props.item.isArray")
          div.caption.font-italic (Array)
        template(slot="label", slot-scope="props")
          div {{ props.item.name }}
            span.pl-1(v-if="props.leaf") :
              ellipsis.pl-1.d-inline-block.grey--text.body-1(:text="String(props.item.value)")
            span.pl-1(v-else-if="!props.leaf && !(props.item.children && props.item.children.length)") :
              .pl-1.d-inline-block.grey--text.text--darken-1.body-1.font-italic {{ $t('common.emptyObject') }}
        template(slot="append", slot-scope="props", v-if="props.leaf")
          v-tooltip(left)
            v-btn(
              v-clipboard:copy="props.item.path",
              v-clipboard:success="() => addSuccessPopup({ text: $t('success.pathCopied') })",
              v-clipboard:error="() => addErrorPopup({ text: $t('errors.default') })",
              slot="activator",
              small,
              icon
              )
              v-icon file_copy
            span {{ $t('modals.variablesHelp.copyToClipboard') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import popupMixin from '@/mixins/popup';
import Ellipsis from '@/components/tables/ellipsis.vue';

export default {
  name: MODALS.variablesHelp,
  components: { Ellipsis },
  mixins: [modalInnerMixin, popupMixin],
};
</script>
