<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.infoPopup.title')}}
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
            v-text-field(
            :placeholder="$t('settings.infoPopup.fields.column')",
            :value="popup.column",
            @input="updateValue(index, 'column', $event)"
            )
          v-flex(xs11)
            v-textarea(
            :placeholder="$t('settings.infoPopup.fields.template')",
            @input="updateValue(index, 'template', $event)",
            :value="popup.template"
            )
      v-btn(color="success", @click="add") {{ $t('common.add') }}
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
