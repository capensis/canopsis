<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.columnNames')}}
    v-container
      v-card.my-2(v-for="(column, index) in value", :key="`alarm-settings-column-${index}`")
        v-layout.pt-2(justify-space-between)
          v-flex(xs3)
            v-layout.text-xs-center.pl-2(justify-space-between)
              v-flex(xs1)
                v-btn(icon, @click.prevent="up(index)")
                  v-icon arrow_upward
              v-flex(xs5)
                v-btn(icon, @click.prevent="down(index)")
                  v-icon arrow_downward
          v-flex.d-flex(xs3)
            div.text-xs-right.pr-2
              v-btn(icon, @click.prevent="remove(index)")
                v-icon(color="red") close
        v-layout(justify-center wrap)
          v-flex(xs11)
            v-text-field(
            :placeholder="$t('common.label')",
            :error-messages="errors.collect(`value[${index}]`)",
            v-model="column.label",
            v-validate="'required'",
            :data-vv-name="`value[${index}]`"
            )
          v-flex(xs11)
            v-text-field(
            :placeholder="$t('common.value')",
            :error-messages="errors.collect(`value[${index}]`)",
            v-model="column.value",
            v-validate="'required'",
            :data-vv-name="`value[${index}]`"
            )
      v-btn(color="success", @click.prevent="add") Add
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [Array, Object],
      default: () => [],
    },
  },
  methods: {
    add() {
      this.$emit('input', [...this.value, { label: '', value: '' }]);
    },
    remove(index) {
      this.$emit('input', this.value.filter((v, i) => i !== index));
    },
    up(index) {
      if (index > 0) {
        const value = [...this.value];
        const temp = value[index];

        value[index] = value[index - 1];
        value[index - 1] = temp;

        this.$emit('input', value);
      }
    },
    down(index) {
      if (index < this.value.length - 1) {
        const value = [...this.value];
        const temp = value[index];

        value[index] = value[index + 1];
        value[index + 1] = temp;

        this.$emit('input', value);
      }
    },
  },
};
</script>
