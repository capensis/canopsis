<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.columnNames')}}
    v-container
      v-card.my-2(v-for="(column, index) in columns", :key="`settings-column-${index}`")
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
              v-btn(icon, @click.prevent="removeItemFromArray(index)")
                v-icon(color="red") close
        v-layout(justify-center wrap)
          v-flex(xs11)
            v-text-field(
            :placeholder="$t('common.label')",
            :error-messages="errors.collect(`label[${index}]`)",
            @input="updateFieldInArrayItem(index, 'label', $event)"
            v-validate="'required'",
            :data-vv-name="`label[${index}]`",
            :value="column.label"
            )
          v-flex(xs11)
            v-text-field(
            :placeholder="$t('common.value')",
            :error-messages="errors.collect(`value[${index}]`)",
            @input="updateFieldInArrayItem(index, 'value', $event)"
            v-validate="'required'",
            :data-vv-name="`value[${index}]`",
            :value="column.value"
            )
      v-btn(color="success", @click.prevent="add") Add
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [
    formMixin,
  ],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    columns: {
      type: [Array, Object],
      default: () => [],
    },
  },
  methods: {
    add() {
      this.addItemIntoArray({ label: '', value: '' });
    },
    up(index) {
      if (index > 0) {
        const value = [...this.columns];
        const temp = value[index];

        value[index] = value[index - 1];
        value[index - 1] = temp;

        this.$emit('input', value);
      }
    },
    down(index) {
      if (index < this.columns.length - 1) {
        const value = [...this.columns];
        const temp = value[index];

        value[index] = value[index + 1];
        value[index + 1] = temp;

        this.$emit('input', value);
      }
    },
  },
};
</script>
