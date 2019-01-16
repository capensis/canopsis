<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ $t('settings.infoPopup.title') }}
      v-layout(justify-end)
        v-btn.primary(
        small,
        @click="edit"
        ) {{ $t('common.create') }}/{{ $t('common.edit') }}
      //v-card.my-2(v-for="(popup, index) in popups", :key="`settings-info-popup-${index}`")
        v-layout(justify-space-between)
          v-flex(xs3)
          v-flex.d-flex(xs3)
            div.text-xs-right.pr-2
              v-btn(icon, @click.prevent="removeItemFromArray(index)")
                v-icon(color="red") close
        v-layout(justify-center wrap)
          v-flex(xs11)
            v-text-field(
            :placeholder="$t('settings.infoPopup.fields.column')",
            :value="popup.column",
            @input="updateFieldInArrayItem(index, 'column', $event)"
            )
          v-flex(xs11)
            text-editor(:value="popup.template", @input="updateFieldInArrayItem(index, 'template', $event)")
        v-btn(color="success", @click="add") {{ $t('common.add') }}
</template>

<script>
import { MODALS } from '@/constants';

import TextEditor from '@/components/other/text-editor/text-editor.vue';

import formMixin from '@/mixins/form';
import modalMixin from '@/mixins/modal';

export default {
  components: {
    TextEditor,
  },
  mixins: [
    formMixin,
    modalMixin,
  ],
  model: {
    prop: 'popups',
    event: 'input',
  },
  props: {
    popups: {
      type: [Array, Object],
      default: () => [],
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    add() {
      this.addItemIntoArray({ column: '', template: '' });
    },
    edit() {
      this.showModal({
        name: MODALS.infoPopupSetting,
        config: {
          infoPopups: this.popups,
          columns: this.columns,
          action: popups => this.$emit('input', popups),
        },
      });
    },
  },
};
</script>
