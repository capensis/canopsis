<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.infoPopup')}}
    v-container
      v-card.my-2(v-for="(popup, index) in columns", :key="`info-popup-${index}`")
        v-layout(justify-space-between)
          v-flex(xs3)
          v-flex.d-flex(xs3)
            div.text-xs-right.pr-2
              v-btn(icon, @click.prevent="remove(index)")
                v-icon(color="red") close
        v-layout(justify-center wrap)
          v-flex(xs11)
            v-text-field(placeholder="Column", @input="updateValue(index, 'column', $event)", v-model="popup.column")
          v-flex(xs11)
            v-text-field(
            placeholder="Template",
            :multi-line="true",
            @input="updateValue(index, 'template', $event)",
            v-model="popup.template"
            )

      v-btn(color="success", @click="add") Add
</template>

<script>
import settingsColumnMixin from '@/mixins/settings-column';

export default {
  mixins: [
    settingsColumnMixin,
  ],
  methods: {
    add() {
      this.columns = [...this.columns, { column: '', template: '' }];

      this.$emit('input', this.columns);
    },
    remove(index) {
      this.columns = this.columns.filter((v, i) => i !== index);

      this.$emit('input', this.columns);
    },
  },
};
</script>
