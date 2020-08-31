<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text.white--text
      v-layout
        v-flex(xs6)
          h3.headline.text-xs-center.my-1 {{ $t('context.impacts') }}
          v-container
            v-card
              v-card-text
                entities-iterator(:entities="impact")
        v-flex(xs6)
          h3.headline.text-xs-center.my-1 {{ $t('context.dependencies') }}
          v-container
            v-card
              v-card-text
                entities-iterator(:entities="depends")
</template>

<script>
import EntitiesIterator from '@/components/other/context/partials/entities-iterator.vue';

export default {
  components: { EntitiesIterator },
  props: {
    impact: {
      type: Array,
      default: () => [],
    },
    depends: {
      type: Array,
      default: () => [],
    },
  },
};
</script>

