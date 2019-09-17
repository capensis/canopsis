<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      v-layout(row, wrap)
        v-flex(xs3)
          v-text-field(v-model="infosSearchingText", :label="$t('context.moreInfos.infosSearchLabel')", dark)
      v-data-table(
        :items="items",
        item-key="item.name",
        :headers="headers",
        :search="infosSearchingText"
      )
        template(slot="items", slot-scope="props")
          td {{ props.item.name }}
          td {{ props.item.description }}
          td {{ props.item.value }}
</template>

<script>
export default {
  props: {
    infos: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      infosSearchingText: '',
      headers: [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.description'),
          value: 'description',
        },
        {
          text: this.$t('common.value'),
          value: 'value',
        },
      ],
    };
  },
  computed: {
    items() {
      return Object.keys(this.infos).map(info => ({
        name: info,
        description: this.infos[info].description,
        value: this.infos[info].value,
      }));
    },
  },
};
</script>
