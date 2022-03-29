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
    v-alert(
      :value="errors.has('filter')",
      type="error",
      transition="fade-transition"
    ) {{ errors.first('filter') }}
    v-alert(
      v-model="countAlertShown",
      type="warning",
      transition="fade-transition",
      dismissible
    )
      span {{ countAlertMessage }}
</template>

<script>
import { isEmpty } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

const { mapActions } = createNamespacedHelpers('pbehavior');

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
  data() {
    return {
      countAlertShown: false,
      countAlertMessage: '',
      countPending: false,
    };
  },
  computed: {
    hasFilter() {
      return this.form.filter && !isEmpty(this.form.filter);
    },
  },
  created() {
    this.attachFilterRule();
  },
  methods: {
    ...mapActions({
      fetchPbehaviorEntitiesCountWithoutStore: 'fetchEntitiesCountWithoutStore',
    }),

    attachFilterRule() {
      this.$validator.attach({
        name: 'filter',
        rules: 'required:true',
        getter: () => this.hasFilter,
        context: () => this,
        vm: this,
      });
    },

    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        dialogProps: {
          zIndex: 300,
        },
        config: {
          filter: this.form.filter,
          withEntity: true,
          action: (filter) => {
            this.updateField('filter', filter);
            this.fetchCountForFilter(filter);
            this.$nextTick(() => this.$validator.validate('filter'));
          },
        },
      });
    },

    async fetchCountForFilter(filter) {
      try {
        this.countPending = true;

        const {
          over_limit: overLimit,
          total_count: totalCount,
        } = await this.fetchPbehaviorEntitiesCountWithoutStore({ data: { filter } });

        if (overLimit) {
          this.countAlertMessage = this.$t('entitiesCountAlerts.filter.countOverLimit', { count: totalCount });
          this.countAlertShown = true;

          return;
        }

        this.countAlertShown = false;
      } catch (err) {
        this.countAlertMessage = this.$t('entitiesCountAlerts.filter.countRequestError');
        this.countAlertShown = true;
      } finally {
        this.countPending = false;
      }
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
