<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.infoPopup')}}
    v-container
      v-card.my-2(v-for="(popup, index) in value", :key="`info-popup-${index}`")
        v-layout(justify-space-between)
          v-flex(xs3)
          v-flex.d-flex(xs3)
            div.text-xs-right.pr-2
              v-btn(icon, @click.prevent="remove(index)")
                v-icon(color="red") close
        v-layout(justify-center wrap)
          v-flex(xs11)
            v-text-field(placeholder="Column", v-model="popup.column")
          v-flex(xs11)
            text-editor(v-model="popup.template", @update:value="updateTemplate($event, index)")
      v-btn(color="success", @click="add") Add
</template>

<script>
import TextEditor from '@/components/other/text-editor/text-editor.vue';

export default {
  components: {
    TextEditor,
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    add() {
      this.$emit('input', [...this.value, { column: '', template: '' }]);
    },
    remove(index) {
      this.$emit('input', this.value.filter((v, i) => i !== index));
    },
    updateTemplate(event, index) {
      this.value[index].template = event;
    },
  },
};
</script>
