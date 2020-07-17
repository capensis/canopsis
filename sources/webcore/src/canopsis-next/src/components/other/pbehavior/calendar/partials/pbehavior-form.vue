<template lang="pug">
  div
    v-layout(wrap)
      v-flex(xs12)
        pbehavior-general-form(v-field="form")
      v-flex(xs12)
        pbehavior-comments-form(v-field="form.comments")
      v-flex(xs12)
        v-btn.ml-0.btn-filter(
          :color="errors.has('filter') ? 'error' : 'primary'",
          @click="showCreateFilterModal"
        ) {{ hasFilter ? 'Edit filter' : 'Add filter' }}
        v-tooltip(v-show="hasFilter", fixed, top)
          v-btn(slot="activator", icon)
            v-icon(color="grey darken-1") info
          span.pre {{ form.filter.filter | json }}
        v-alert(:value="errors.has('filter')", type="error") {{ errors.first('filter') }}
      v-flex(xs12)
        v-btn.ml-0(
          color="primary",
          @click="showCreateRRuleModal"
        )  {{ hasRRule ? 'Edit RRule' : 'Add RRule' }}
        v-tooltip(v-show="hasRRule", fixed, top)
          v-btn(slot="activator", icon)
            v-icon(color="grey darken-1") info
          span {{ form.rrule }}
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import formMixin from '@/mixins/form/object';

import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorCommentsForm from './pbehavior-comments-form.vue';

export default {
  components: {
    PbehaviorGeneralForm,
    PbehaviorCommentsForm,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      required: true,
    },
    noFilter: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasFilter() {
      return this.form.filter && !isEmpty(this.form.filter.filter);
    },
    hasRRule() {
      return !isEmpty(this.form.rrule);
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
          filter: this.form.filter,
          hiddenFields: ['title'],
          action: (filter) => {
            this.updateField('filter', filter);
            this.$nextTick(() => this.$validator.validate('filter'));
          },
        },
      });
    },
    showCreateRRuleModal() {
      this.$modals.show({
        name: MODALS.createRRule,
        dialogProps: {
          zIndex: 300,
        },
        config: {
          rrule: this.form.rrule,
          action: rrule => this.updateField('rrule', rrule),
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .pre {
    white-space: pre;
  }

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
