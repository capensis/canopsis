<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
      .font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      v-text-field(:value="value", type="number", @input="$emit('input', parseInt($event, 10))")
</template>

<script>
/**
 * Component to select a Number type setting
 *
 * @prop {Number} [value=undefined] - Value
 * @prop {String} [title=''] - Title
 */
export default {
  props: {
    value: {
      type: Number,
      default: undefined,
    },
    title: {
      type: String,
      default: '',
    },
  },
};
</script>
