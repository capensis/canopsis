<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.infoPopup.title')}}
    v-container
      v-card.my-2(v-for="(popup, index) in popups", :key="`settings-info-popup-${index}`")
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
import TextEditor from '@/components/other/text-editor/text-editor.vue';
import formMixin from '@/mixins/form';

export default {
  components: {
    TextEditor,
  },
  mixins: [
    formMixin,
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
  },
  methods: {
    add() {
      this.addItemIntoArray({ column: '', template: '' });
    },
  },
};
</script>
