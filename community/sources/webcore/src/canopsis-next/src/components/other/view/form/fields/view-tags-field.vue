<template>
  <v-combobox
    v-field="value"
    :label="$t('view.groupTags')"
    append-icon=""
    tags
    clearable
    multiple
    chips
    deletable-chips
  />
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
