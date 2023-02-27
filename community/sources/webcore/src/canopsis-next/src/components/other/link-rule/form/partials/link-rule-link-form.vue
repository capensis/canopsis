<template lang="pug">
  v-card
    v-card-text
      v-layout(column)
        v-layout(row)
          v-text-field.mr-2(
            v-field="form.label",
            v-validate="'required'",
            :label="$t('common.label')",
            :error-messages="errors.collect(labelFieldName)",
            :name="labelFieldName"
          )
          v-btn.mr-0(icon, @click="remove")
            v-icon(color="error") delete
        v-layout(row)
          v-flex.pr-2(xs8)
            v-text-field(
              v-field="form.category",
              :label="$t('common.category')"
            )
          v-flex(xs4)
            c-icon-field(
              v-field="form.icon_name",
              :label="$t('common.icon')",
              :name="iconFieldName",
              required
            )
        v-text-field(
          v-field="form.url",
          v-validate="'required'",
          :label="$t('common.url')",
          :error-messages="errors.collect(urlFieldName)",
          :name="urlFieldName"
        )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'link',
    },
  },
  computed: {
    labelFieldName() {
      return `${this.name}.label`;
    },

    iconFieldName() {
      return `${this.name}.icon`;
    },

    urlFieldName() {
      return `${this.name}.url`;
    },
  },
  methods: {
    remove() {
      this.$emit('remove');
    },
  },
};
</script>
