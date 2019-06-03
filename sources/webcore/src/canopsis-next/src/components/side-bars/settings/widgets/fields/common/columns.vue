<template lang="pug">
  v-list-group
    v-list-tile(slot="activator")
      .validation-header(:class="{ 'error--text': hasAnyError }") {{ $t('settings.columnNames') }}
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
            v-validate="'required'",
            :placeholder="$t('common.label')",
            :error-messages="errors.collect(`label[${index}]`)",
            :name="`label[${index}]`",
            :value="column.label",
            @input="updateFieldInArrayItem(index, 'label', $event)"
            )
          v-flex(xs11)
            v-text-field(
            v-validate="'required'",
            :placeholder="$t('common.value')",
            :error-messages="errors.collect(`value[${index}]`)",
            :value="column.value"
            :name="`value[${index}]`",
            @input="updateFieldInArrayItem(index, 'value', $event)"
            )
          v-flex(v-if="withHtml", xs11)
            v-switch(
            :label="$t('settings.columns.isHtml')",
            :input-value="column.isHtml",
            @change="updateFieldInArrayItem(index, 'isHtml', $event)"
            )
      v-btn(color="primary", @click.prevent="add") {{ $t('common.add') }}
</template>

<script>
import formArrayMixin from '@/mixins/form/array';
import formValidatorErrorMixin from '@/mixins/form/validator-error';

export default {
  inject: ['$validator'],
  mixins: [
    formArrayMixin,
    formValidatorErrorMixin,
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
    withHtml: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    add() {
      const column = { label: '', value: '' };

      if (this.withHtml) {
        column.isHtml = false;
      }

      this.addItemIntoArray(column);
    },
    up(index) {
      if (index > 0) {
        const columns = [...this.columns];
        const temp = columns[index];

        columns[index] = columns[index - 1];
        columns[index - 1] = temp;

        this.$emit('input', columns);
      }
    },
    down(index) {
      if (index < this.columns.length - 1) {
        const columns = [...this.columns];
        const temp = columns[index];

        columns[index] = columns[index + 1];
        columns[index + 1] = temp;

        this.$emit('input', columns);
      }
    },
  },
};
</script>
