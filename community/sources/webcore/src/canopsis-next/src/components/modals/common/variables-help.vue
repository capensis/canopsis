<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.variablesHelp.variables') }}
    template(#text="")
      v-treeview(:items="config.variables", item-key="name")
        template(#prepend="{ item }")
          div.caption.font-italic(v-if="item.isArray") {{`(${$t('common.variableTypes.array')})`}}

        template(#label="{ item, leaf }")
          div {{ item.name }}
            span.pl-1(v-if="leaf") :
              c-ellipsis.pl-1.d-inline-block.grey--text.body-1.pre-wrap(:text="String(item.value)")
            span.pl-1(v-else-if="!leaf && !(item.children && item.children.length)") :
              div.pl-1.d-inline-block.grey--text.text--darken-1.body-1.font-italic {{ $t('common.emptyObject') }}

        template(#append="{ leaf, item }")
          c-copy-btn(
            v-if="leaf",
            :value="item.path",
            :tooltip="$t('modals.variablesHelp.copyToClipboard')",
            left,
            @success="onSuccessCopied",
            @error="onErrorCopied"
          )
          c-action-btn(
            v-else,
            :tooltip="$t('common.export')",
            icon="file_download",
            left,
            @click="exportAsJson(item)"
          )
</template>

<script>
import { DATETIME_FORMATS, MODALS } from '@/constants';

import { saveJsonFile } from '@/helpers/file/files';
import { convertTreeviewToObject } from '@/helpers/treeview';
import { convertDateToString, getNowTimestamp } from '@/helpers/date/date';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.variablesHelp,
  components: { ModalWrapper },
  mixins: [modalInnerMixin],
  methods: {
    exportAsJson(item) {
      const object = convertTreeviewToObject(item);

      const dateString = convertDateToString(
        getNowTimestamp(),
        DATETIME_FORMATS.long,
      );

      saveJsonFile(object, `${item.name}-${dateString}`);
    },

    onSuccessCopied() {
      this.$popups.success({ text: this.$t('success.pathCopied') });
    },

    onErrorCopied() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>
