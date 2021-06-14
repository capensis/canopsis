<template lang="pug">
  v-flex
    v-btn.ml-0.btn-filter(
      :color="errors.has('filter') ? 'error' : 'primary'",
      @click="showCreateFilterModal"
    ) {{ hasFilter ? $t('pbehavior.buttons.editFilter') : $t('pbehavior.buttons.addFilter') }}
    v-tooltip(v-show="hasFilter", fixed, top)
      v-btn(slot="activator", icon)
        v-icon(color="grey darken-1") info
      span.pre {{ form.filter | json }}
    v-alert(:value="errors.has('filter')", type="error") {{ errors.first('filter') }}
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import formMixin from '@/mixins/form/object';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    hasFilter() {
      return this.form.filter && !isEmpty(this.form.filter);
    },
  },
  created() {
    this.$validator.attach({
      name: 'filter',
      rules: 'required:true',
      getter: () => !isEmpty(this.form.filter),
      context: () => this,
      vm: this,
    });
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        dialogProps: {
          zIndex: 300,
        },
        config: {
          filter: { filter: this.form.filter },
          hiddenFields: ['title'],
          action: ({ filter }) => {
            this.updateField('filter', filter);
            this.$nextTick(() => this.$validator.validate('filter'));
          },
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .btn-filter.error {
    animation: shake .6s cubic-bezier(.25,.8,.5,1);
  }

  @keyframes shake {
    10%, 90% {
      transform: translate3d(-1px, 0, 0);
    }

    20%, 80% {
      transform: translate3d(2px, 0, 0);
    }

    30%, 50%, 70% {
      transform: translate3d(-4px, 0, 0);
    }

    40%, 60% {
      transform: translate3d(4px, 0, 0);
    }
  }
</style>
