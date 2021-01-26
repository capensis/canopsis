<template lang="pug">
  v-layout(row, wrap, align-center)
    v-flex(xs11)
      v-chip(
        v-for="entity in entities",
        :key="entity[itemKey]",
        close,
        @input="$emit('remove', entity)"
      )
        slot(:entity="entity") {{ contentKey ? entity[contentKey] : entity }}
    v-flex(v-if="clearable && entities.length", xs1)
      v-tooltip(right)
        v-btn(
          slot="activator",
          small,
          icon,
          @click="$emit('clear')"
        )
          v-icon(color="error") delete
        span {{ $t('common.deleteAll') }}
</template>

<script>
export default {
  props: {
    entities: {
      type: Array,
      required: true,
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    contentKey: {
      type: String,
      required: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
