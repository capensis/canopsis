<template lang="pug">
  v-container.pa-3(data-test="statSelector", fluid)
    v-layout(align-center, justify-space-between)
      .subheading(:class="validationHeaderClass") {{ $t('settings.statsSelect.title') }}
        .font-italic.caption.ml-1 ({{ $t('settings.statsNumbers.defaultStat') }})
      v-btn.primary(
        data-test="selectButton",
        small,
        @click="openAddStatModal"
      ) {{ $t('common.select') }}
</template>

<script>
import { STATS_TYPES, MODALS } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
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
      vm: this,
    });
  },
  methods: {
    openAddStatModal() {
      this.$modals.show({
        name: MODALS.addStat,
        config: {
          title: this.$t('modals.addStat.title.add'),
          stat: this.stat,
          statTitle: this.stat.title,
          action: stat => this.$emit('input', stat),
        },
      });
    },
  },
};
</script>
