<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.infoPopup.title')}}
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
            v-text-field(v-model="popup.column", :placeholder="$t('settings.infoPopup.fields.column')")
          v-flex(xs11)
            v-textarea(v-model="popup.template", :placeholder="$t('settings.infoPopup.fields.template')")

      v-btn(color="success", @click="add") {{ $t('common.add') }}
</template>

<script>
export default {
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
  },
};
</script>
