<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text.white--text
      v-layout
        v-flex(xs6)
          h3.headline.text-xs-center.my-1 {{ $t('context.impacts') }}
          v-container
            v-card
              v-card-text
                v-data-iterator(:items="impact")
                  v-flex(slot="item", slot-scope="props")
                    v-card
                      v-card-title {{ props.item }}
                  v-flex(slot="no-data")
                    v-card
                      v-card-title {{ $t('tables.noData') }}
        v-flex(xs6)
          h3.headline.text-xs-center.my-1 {{ $t('context.dependencies') }}
          v-container
            v-card
              v-card-text
                v-data-iterator(:items="depends")
                  v-flex(slot="item", slot-scope="props")
                    v-card
                      v-card-title {{ props.item }}
                  v-flex(slot="no-data")
                    v-card
                      v-card-title {{ $t('tables.noData') }}
</template>

<script>
export default {
  props: {
    impact: {
      type: Array,
      required: true,
    },
    depends: {
      type: Array,
      required: true,
    },
  },
};
</script>

