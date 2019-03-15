<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      div.subheading {{ $t('settings.statsSelect.title') }}
        .font-italic.caption.ml-1 ({{ $t('settings.statsNumbers.defaultStat') }})
      v-btn.primary(
      small,
      @click="openAddStatModal"
      ) {{ $t('common.select') }}
</template>


<script>
import modalMixin from '@/mixins/modal';
import { STATS_TYPES, MODALS } from '@/constants';

export default {
  inject: ['$validator'],
  mixins: [modalMixin],
  model: {
    prop: 'stat',
    event: 'input',
  },
  props: {
    stat: {
      type: Object,
      required: true,
    },
  },
  computed: {
    /**
     * Get stats different types from constant, and return an object with stat's value and stat's translated title
     */
    statsTypes() {
      return Object.values(STATS_TYPES)
        .map(item => ({ value: item.value, text: this.$t(`stats.types.${item.value}`), options: item.options }));
    },
  },
  created() {
    this.$validator.attach({
      name: 'stat',
      rules: 'required:true',
      getter: () => this.stat,
      context: () => this,
    });
  },
  methods: {
    openAddStatModal() {
      this.showModal({
        name: MODALS.addStat,
        config: {
          title: 'modals.addStat.title.add',
          stat: this.stat,
          statTitle: this.stat.title,
          action: stat => this.$emit('input', stat),
        },
      });
    },
  },
};
</script>
