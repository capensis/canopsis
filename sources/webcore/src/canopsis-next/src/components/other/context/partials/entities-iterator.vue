<template lang="pug">
  v-data-iterator(:items="entities")
    v-layout(slot="header")
      v-layout.font-weight-bold.subtitle-1(xs6, justify-center) {{ $t('common.id') }}
      v-layout.font-weight-bold.subtitle-1(xs6, justify-center) {{ $t('common.name') }}
    v-flex(slot="item", slot-scope="props")
      v-card
        v-card-title
          v-flex(xs6) {{ props.item.id }}
          v-flex(xs6) {{ props.item.name }}
    v-flex(slot="no-data")
      v-card
        v-card-title {{ $t('tables.noData') }}
</template>

<script>
export default {
  props: {
    entities: {
      type: Array,
      default: () => [],
    },
  },
};
</script>

